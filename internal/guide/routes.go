package guide

import "github.com/gin-gonic/gin"

func RegisterRoutes(r gin.IRouter, h *Handler, logged, creators, admin gin.HandlerFunc) {
	guides := r.Group("/guides")
	{
		guides.GET("", h.GetGuides)
		guides.GET("/:id", h.GetGuideById)
		guides.POST("", logged, creators, h.CreateGuide)
		guides.PATCH("/:id", logged, creators, h.UpdateGuide)
		guides.DELETE("/:id", logged, admin, h.DeleteGuideById)
	}
}
