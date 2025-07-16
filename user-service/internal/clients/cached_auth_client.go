package clients

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"shared/proto/auth_service"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

type CachedAuthClient struct {
	authClient   *AuthServiceClient
	l1Cache      *cache.Cache  // In-memory cache (L1)
	l2Cache      *redis.Client // Redis cache (L2)
	l1TTL        time.Duration
	l2TTL        time.Duration
	cacheEnabled bool
	redisEnabled bool
	metrics      *CacheMetrics
}

type CacheMetrics struct {
	L1Hits        int64
	L1Misses      int64
	L2Hits        int64
	L2Misses      int64
	L2Errors      int64
	GrpcCalls     int64
	TotalRequests int64
}

type CachedToken struct {
	UserID    uint32 `json:"user_id"`
	Email     string `json:"email"`
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	ExpiresAt int64  `json:"expires_at"`
	IssuedAt  int64  `json:"issued_at"`
	CachedAt  int64  `json:"cached_at"`
}

func NewCachedAuthClient(authServiceAddr, redisAddr string, l1TTL, l2TTL time.Duration) (*CachedAuthClient, error) {
	// Initialize auth service client
	authClient, err := NewAuthServiceClient(authServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %v", err)
	}

	// Initialize L1 cache (in-memory)
	l1Cache := cache.New(l1TTL, l1TTL*2) // cleanup interval = 2x TTL

	// Initialize L2 cache (Redis)
	var l2Cache *redis.Client
	var redisEnabled bool

	if redisAddr != "" {
		l2Cache = redis.NewClient(&redis.Options{
			Addr:         redisAddr,
			PoolSize:     10,
			MinIdleConns: 2,
			MaxRetries:   3,
			ReadTimeout:  100 * time.Millisecond,
			WriteTimeout: 100 * time.Millisecond,
		})

		// Test Redis connection
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := l2Cache.Ping(ctx).Err(); err != nil {
			log.Printf("Redis connection failed, running without L2 cache: %v", err)
			redisEnabled = false
		} else {
			log.Println("Redis L2 cache connected successfully")
			redisEnabled = true
		}
	}

	return &CachedAuthClient{
		authClient:   authClient,
		l1Cache:      l1Cache,
		l2Cache:      l2Cache,
		l1TTL:        l1TTL,
		l2TTL:        l2TTL,
		cacheEnabled: true,
		redisEnabled: redisEnabled,
		metrics:      &CacheMetrics{},
	}, nil
}

func (h *CachedAuthClient) ValidateToken(ctx context.Context, token string) (*auth_service.ValidateTokenResponse, error) {
	h.metrics.TotalRequests++

	if !h.cacheEnabled {
		h.metrics.GrpcCalls++
		return h.authClient.ValidateToken(ctx, token)
	}

	cacheKey := h.generateCacheKey(token)

	// Step 1: Check L1 cache (in-memory)
	if cachedToken, found := h.checkL1Cache(cacheKey); found {
		h.metrics.L1Hits++
		log.Printf("L1 CACHE HIT for token: %s", cacheKey[:12]+"...")
		return h.buildResponse(cachedToken), nil
	}
	h.metrics.L1Misses++

	// Step 2: Check L2 cache (Redis)
	if h.redisEnabled {
		if cachedToken, found := h.checkL2Cache(ctx, cacheKey); found {
			h.metrics.L2Hits++
			log.Printf("L2 CACHE HIT for token: %s", cacheKey[:12]+"...")

			// Promote to L1 cache
			h.storeInL1(cacheKey, cachedToken)
			return h.buildResponse(cachedToken), nil
		}
		h.metrics.L2Misses++
	}

	// Step 3: Cache miss - call auth service
	log.Printf("CACHE MISS - calling auth service for token: %s", cacheKey[:12]+"...")
	h.metrics.GrpcCalls++

	response, err := h.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Step 4: Store in both caches if validation successful
	if response.IsValid && response.Claims != nil {
		cachedToken := h.convertToCachedToken(response.Claims)

		// Store in L1 (always)
		h.storeInL1(cacheKey, cachedToken)

		// Store in L2 (if Redis available)
		if h.redisEnabled {
			h.storeInL2(ctx, cacheKey, cachedToken)
		}
	}

	return response, nil
}

func (h *CachedAuthClient) checkL1Cache(cacheKey string) (*CachedToken, bool) {
	if item, found := h.l1Cache.Get(cacheKey); found {
		if cachedToken, ok := item.(*CachedToken); ok {
			// Check if token is still valid
			if time.Now().Unix() < cachedToken.ExpiresAt {
				return cachedToken, true
			}
			// Token expired, remove from cache
			h.l1Cache.Delete(cacheKey)
		}
	}
	return nil, false
}

func (h *CachedAuthClient) checkL2Cache(ctx context.Context, cacheKey string) (*CachedToken, bool) {
	val, err := h.l2Cache.Get(ctx, cacheKey).Result()
	if err != nil {
		if err != redis.Nil {
			h.metrics.L2Errors++
			log.Printf("Redis L2 cache error: %v", err)
		}
		return nil, false
	}

	var cachedToken CachedToken
	if err := json.Unmarshal([]byte(val), &cachedToken); err != nil {
		log.Printf("Failed to unmarshal cached token: %v", err)
		// Remove corrupted cache entry
		h.l2Cache.Del(ctx, cacheKey)
		return nil, false
	}

	// Check if token is still valid
	if time.Now().Unix() < cachedToken.ExpiresAt {
		return &cachedToken, true
	}

	// Token expired, remove from cache
	h.l2Cache.Del(ctx, cacheKey)
	return nil, false
}

func (h *CachedAuthClient) storeInL1(cacheKey string, cachedToken *CachedToken) {
	ttl := h.calculateTTL(cachedToken.ExpiresAt, h.l1TTL)
	h.l1Cache.Set(cacheKey, cachedToken, ttl)
}

func (h *CachedAuthClient) storeInL2(ctx context.Context, cacheKey string, cachedToken *CachedToken) {
	data, err := json.Marshal(cachedToken)
	if err != nil {
		log.Printf("Failed to marshal token for L2 cache: %v", err)
		return
	}

	ttl := h.calculateTTL(cachedToken.ExpiresAt, h.l2TTL)
	if err := h.l2Cache.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		h.metrics.L2Errors++
		log.Printf("Failed to store in L2 cache: %v", err)
	}
}

func (h *CachedAuthClient) generateCacheKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("auth:token:%x", hash)
}

func (h *CachedAuthClient) calculateTTL(expiresAt int64, maxTTL time.Duration) time.Duration {
	timeUntilExpiry := time.Until(time.Unix(expiresAt, 0))

	if timeUntilExpiry <= 0 {
		return time.Minute // Minimum TTL for expired tokens (edge case)
	}

	// Use 80% of remaining time, but not more than maxTTL
	calculatedTTL := time.Duration(float64(timeUntilExpiry) * 0.8)
	if calculatedTTL > maxTTL {
		calculatedTTL = maxTTL
	}

	// Ensure minimum TTL of 30 seconds
	if calculatedTTL < 30*time.Second {
		calculatedTTL = 30 * time.Second
	}

	return calculatedTTL
}

func (h *CachedAuthClient) convertToCachedToken(claims *auth_service.UserClaims) *CachedToken {
	return &CachedToken{
		UserID:    claims.UserId,
		Email:     claims.Email,
		Issuer:    claims.Issuer,
		Subject:   claims.Subject,
		ExpiresAt: claims.ExpiresAt,
		IssuedAt:  claims.IssuedAt,
		CachedAt:  time.Now().Unix(),
	}
}

func (h *CachedAuthClient) buildResponse(cachedToken *CachedToken) *auth_service.ValidateTokenResponse {
	return &auth_service.ValidateTokenResponse{
		IsValid:      true,
		ErrorMessage: "",
		Claims: &auth_service.UserClaims{
			Email:     cachedToken.Email,
			UserId:    cachedToken.UserID,
			Issuer:    cachedToken.Issuer,
			Subject:   cachedToken.Subject,
			ExpiresAt: cachedToken.ExpiresAt,
			IssuedAt:  cachedToken.IssuedAt,
		},
	}
}

func (h *CachedAuthClient) GetMetrics() *CacheMetrics {
	return h.metrics
}

func (h *CachedAuthClient) GetCacheStats() map[string]interface{} {
	l1Items := len(h.l1Cache.Items())

	var l2Items int64
	if h.redisEnabled {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Get approximate count of auth tokens in Redis
		keys, _ := h.l2Cache.Keys(ctx, "auth:token:*").Result()
		l2Items = int64(len(keys))
	}

	hitRate := float64(0)
	if h.metrics.TotalRequests > 0 {
		totalHits := h.metrics.L1Hits + h.metrics.L2Hits
		hitRate = float64(totalHits) / float64(h.metrics.TotalRequests) * 100
	}

	return map[string]interface{}{
		"cache_enabled":    h.cacheEnabled,
		"redis_enabled":    h.redisEnabled,
		"l1_items":         l1Items,
		"l2_items":         l2Items,
		"l1_ttl":           h.l1TTL,
		"l2_ttl":           h.l2TTL,
		"hit_rate_percent": fmt.Sprintf("%.2f", hitRate),
		"metrics":          h.metrics,
	}
}

func (h *CachedAuthClient) ClearCache() {
	// Clear L1
	h.l1Cache.Flush()

	// Clear L2
	if h.redisEnabled {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		keys, err := h.l2Cache.Keys(ctx, "auth:token:*").Result()
		if err == nil && len(keys) > 0 {
			h.l2Cache.Del(ctx, keys...)
		}
	}

	log.Println("Multi-tier auth cache cleared (L1 + L2)")
}

func (h *CachedAuthClient) Close() error {
	h.ClearCache()

	if h.l2Cache != nil {
		h.l2Cache.Close()
	}

	return h.authClient.Close()
}
