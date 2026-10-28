package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
	"go.uber.org/zap"
)

func logRequest(c *gin.Context, handler string, action string) {
	zap.L().Info("request",
		zap.String("handler", handler),
		zap.String("action", action),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)
}

// internalError loga o erro real e devolve uma mensagem genérica ao cliente.
func internalError(c *gin.Context, handler string, action string, err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		zap.L().Warn("timeout", zap.String("handler", handler), zap.String("action", action))
		c.JSON(http.StatusGatewayTimeout, dto.ErrorResponse{Error: "tempo de resposta esgotado"})
		return
	}
	zap.L().Error("internal error",
		zap.String("handler", handler),
		zap.String("action", action),
		zap.String("path", c.Request.URL.Path),
		zap.Error(err),
	)
	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "erro interno"})
}
