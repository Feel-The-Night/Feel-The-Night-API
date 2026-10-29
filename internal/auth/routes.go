package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r gin.IRouter, h *Handler, limiter gin.HandlerFunc) {
	r.POST("/auth/login", limiter, h.Login)
	// r.POST("/auth/refresh", limiter, h.Refresh)
}
