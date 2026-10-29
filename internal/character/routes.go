package character

import "github.com/gin-gonic/gin"

func RegisterRoutes(r gin.IRouter, h *Handler, logged, creators, admin gin.HandlerFunc) {
	characters := r.Group("/characters")
	{
		characters.GET("", h.GetCharacters)
		characters.GET("/:id", h.GetCharacterById)
		characters.POST("", logged, creators, h.CreateCharacter)
		characters.PATCH("/:id", logged, creators, h.UpdateCharacter)
		characters.DELETE("/:id", logged, admin, h.DeleteCharacterById)
	}
}
