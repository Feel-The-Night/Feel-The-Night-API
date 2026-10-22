package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kisalto/Feel-The-Night-Swagger/docs"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/database"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

// SetupRoutes registra todas as rotas da sua API
func SetupRoutes(
	userHandler *UserHandler,
	characterHandler *CharacterHandler,
	guideHandler *GuideHandler, eventHandler *EventHandler,
	authHandler *AuthHandler,
) *gin.Engine {
	router := gin.Default()

	// Rota do Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas gerais
	router.GET("/health", HealthCheckHandler)
	router.POST("/auth/login", authHandler.Login)
	// router.POST("/refresh")

	// Rotas de Usuários
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.GetUserById)
	router.DELETE("/users/:id", userHandler.DeleteUserById)
	router.PATCH("/users/:id", userHandler.UpdateUser)

	// Rotas de Personagens
	router.POST("/characters", characterHandler.CreateCharacter)
	router.GET("/characters/:id", characterHandler.GetCharacterById)
	router.DELETE("/characters/:id", characterHandler.DeleteCharacterById)
	router.PATCH("/characters/:id", characterHandler.UpdateCharacter)

	// Rotas de Guias
	router.POST("/guides", guideHandler.CreateGuide)
	router.GET("/guides/:id", guideHandler.GetGuideById)
	router.DELETE("/guides/:id", guideHandler.DeleteGuideById)
	router.PATCH("/guides/:id", guideHandler.UpdateGuide)

	// Rotas de Eventos
	router.POST("/events", eventHandler.CreateEvent)
	router.GET("/events/:id", eventHandler.GetEventById)
	router.DELETE("/events/:id", eventHandler.DeleteEventById)
	router.PATCH("/events/:id", eventHandler.UpdateEvent)

	return router
}

func Setup() *gin.Engine {
	// Inicialização de serviços e handlers
	// User
	userService := services.NewUserService(database.DB)
	userHandler := NewUserHandler(userService)

	// Character
	characterService := services.NewCharacterService(database.DB)
	characterHandler := NewCharacterHandler(characterService)

	// Guides
	guideService := services.NewGuideService(database.DB)
	guideHandler := NewGuideHandler(guideService)

	// Eventos
	eventService := services.NewEventService(database.DB)
	eventHandler := NewEventHandler(eventService)

	// Auth
	authHandler := NewAuthHandler(userService)

	return SetupRoutes(userHandler, characterHandler, guideHandler, eventHandler, authHandler)
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
