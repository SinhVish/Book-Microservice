package handlers

import (
	"net/http"

	"user-service/internal/clients"

	"github.com/gin-gonic/gin"
)

type CacheHandler struct {
	authClient *clients.CachedAuthClient
}

func NewCacheHandler(authClient *clients.CachedAuthClient) *CacheHandler {
	return &CacheHandler{
		authClient: authClient,
	}
}

// GetCacheStats returns cache statistics
func (h *CacheHandler) GetCacheStats(c *gin.Context) {
	stats := h.authClient.GetCacheStats()

	c.JSON(http.StatusOK, gin.H{
		"service": "user-service",
		"cache":   stats,
	})
}

// GetCacheMetrics returns detailed cache metrics
func (h *CacheHandler) GetCacheMetrics(c *gin.Context) {
	metrics := h.authClient.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"service": "user-service",
		"metrics": metrics,
	})
}

// ClearCache clears both L1 and L2 cache
func (h *CacheHandler) ClearCache(c *gin.Context) {
	h.authClient.ClearCache()

	c.JSON(http.StatusOK, gin.H{
		"service": "user-service",
		"message": "Cache cleared successfully",
	})
}
