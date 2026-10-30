package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r gin.IRouter, h *Handler, limiter, logged, admin gin.HandlerFunc) {
	users := r.Group("/users")
	{
		users.POST("", limiter, h.CreateUser)
		users.GET("", logged, admin, h.GetUserByEmail)
		users.GET("/:id", h.GetUserById)
		users.PATCH("/:id", logged, h.UpdateUser) // "owner or admin" is checked in the handler
		users.DELETE("/:id", logged, admin, h.DeleteUserById)
	}
}
