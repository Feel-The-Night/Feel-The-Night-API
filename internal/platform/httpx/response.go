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

// internalError loga o erro real e devolve uma mensagem genérica ao cliente.
func InternalError(c *gin.Context, handler string, action string, err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		zap.L().Warn("timeout", zap.String("handler", handler), zap.String("action", action))
		c.JSON(http.StatusGatewayTimeout, ErrorResponse{Error: "tempo de resposta esgotado"})
		return
	}
	zap.L().Error("internal error",
		zap.String("handler", handler),
		zap.String("action", action),
		zap.String("path", c.Request.URL.Path),
		zap.Error(err),
	)
	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "erro interno"})
}

// PageResponse serve para qualquer listagem paginada (guias, eventos, personagens).
type PageResponse[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// ErrorResponse é o formato padrão para respostas de erro da API.
// @Title ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"o recurso solicitado não foi encontrado"`
}
