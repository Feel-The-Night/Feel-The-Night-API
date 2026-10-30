package httpx

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogRequest(c *gin.Context, handler string, action string) {
	zap.L().Info("request",
		zap.String("handler", handler),
		zap.String("action", action),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)
}

// InternalError logs the real error and returns a generic message to the client.
func InternalError(c *gin.Context, handler string, action string, err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		zap.L().Warn("timeout", zap.String("handler", handler), zap.String("action", action))
		c.JSON(http.StatusGatewayTimeout, ErrorResponse{Error: "request timed out"})
		return
	}
	zap.L().Error("internal error",
		zap.String("handler", handler),
		zap.String("action", action),
		zap.String("path", c.Request.URL.Path),
		zap.Error(err),
	)
	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
}

// PageResponse is the envelope for any paginated list (guides, events, characters).
type PageResponse[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// ErrorResponse is the standard format for API error responses.
// @Title ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"the requested resource was not found"`
}
