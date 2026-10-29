package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/httpx"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/jwt"
	"go.uber.org/zap"
)

const claimsKey = "claims"

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))

		tokenString := header
		if len(header) > 7 && strings.EqualFold(header[:7], "Bearer ") {
			tokenString = strings.TrimSpace(header[7:])
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: "token ausente"})
			return
		}

		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			msg := "token invalido"
			if errors.Is(err, jwt.ErrExpiredToken) {
				msg = "token expirado"
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: msg})
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

func RequireRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			// RequireAuth não rodou antes (erro de configuração da rota)
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: "missing or invalid token"})
			return
		}

		if !slices.Contains(roles, models.Role(claims.Role)) {
			zap.L().Warn("acesso negado",
				zap.Int64("user_id", claims.UserID),
				zap.String("role", claims.Role),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, httpx.ErrorResponse{Error: "access denied"})
			return
		}

		c.Next()
	}
}

func GetClaims(c *gin.Context) *jwt.Claims {
	v, ok := c.Get(claimsKey)
	if !ok {
		return nil
	}
	claims, _ := v.(*jwt.Claims)
	return claims
}
