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
	general := router.Group("/")
	{
		general.GET("/health", HealthCheckHandler)
		general.POST("/auth/login", authHandler.Login)
		// router.POST("/refresh")
	}

	// Rotas de Usuários
	users := router.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("/", userHandler.GetUserByEmail)
		users.GET("/:id", userHandler.GetUserById)
		users.DELETE("/:id", userHandler.DeleteUserById)
		users.PATCH("/:id", userHandler.UpdateUser)
	}

	// Rotas de Personagens
	characters := router.Group("/characters")
	{
		characters.POST("", characterHandler.CreateCharacter)
		characters.GET("/:id", characterHandler.GetCharacterById)
		characters.DELETE("/:id", characterHandler.DeleteCharacterById)
		characters.PATCH("/:id", characterHandler.UpdateCharacter)
	}

	// Rotas de Guias
	guides := router.Group("/guides")
	{
		guides.POST("", guideHandler.CreateGuide)
		guides.GET("/:id", guideHandler.GetGuideById)
		guides.DELETE("/:id", guideHandler.DeleteGuideById)
		guides.PATCH("/:id", guideHandler.UpdateGuide)
	}

	// Rotas de Eventos
	events := router.Group("/events")
	{
		events.POST("", eventHandler.CreateEvent)
		events.GET("/:id", eventHandler.GetEventById)
		events.DELETE("/:id", eventHandler.DeleteEventById)
		events.PATCH("/:id", eventHandler.UpdateEvent)
	}

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
