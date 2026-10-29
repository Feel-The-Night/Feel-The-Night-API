package event

import "github.com/gin-gonic/gin"

func RegisterRoutes(r gin.IRouter, h *Handler, logged, creators, admin gin.HandlerFunc) {
	events := r.Group("/events")
	{
		events.GET("", h.GetEvents)
		events.GET("/:id", h.GetEventById)
		events.POST("", logged, creators, h.CreateEvent)
		events.PATCH("/:id", logged, creators, h.UpdateEvent)
		events.DELETE("/:id", logged, admin, h.DeleteEventById)
	}
}
