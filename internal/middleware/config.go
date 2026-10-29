package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:5174,http://localhost:5175"

func CorsMiddleware() gin.HandlerFunc {
	// CORS_ORIGINS=https://meusite.com,http://localhost:5173
	originString := os.Getenv("CORS_ORIGINS")
	if originString == "" {
		originString = defaultOrigins
	}

	allowed := make(map[string]bool)
	for _, o := range strings.Split(originString, ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowed[o] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" && allowed[origin] {
			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			h.Set("Access-Control-Max-Age", "7200") // 2h: cacheia o preflight
			h.Add("Vary", "Origin")

			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
