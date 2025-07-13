package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// GetErrorResponse converts an error to a standardized response format
func GetErrorResponse(err error) (int, ErrorResponse) {
	if err == nil {
		return http.StatusOK, ErrorResponse{}
	}

	// Check if it's our custom error
	var customErr *CustomError
	if errors.As(err, &customErr) {
		return customErr.StatusCode(), ErrorResponse{
			Error: customErr.Message,
			Code:  customErr.Code,
		}
	}

	// Handle specific known errors
	switch err.Error() {
	case "record not found":
		return http.StatusNotFound, ErrorResponse{
			Error: "Resource not found",
			Code:  http.StatusNotFound,
		}
	}

	// Default to internal server error
	return http.StatusInternalServerError, ErrorResponse{
		Error: "Internal server error",
		Code:  http.StatusInternalServerError,
	}
}

// HandleError handles errors in Gin handlers consistently
func HandleError(c *gin.Context, err error) {
	statusCode, response := GetErrorResponse(err)
	c.JSON(statusCode, response)
}
