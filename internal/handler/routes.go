package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kisalto/Feel-The-Night-Swagger/docs"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/database"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/middleware"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

// SetupRoutes registra todas as rotas da sua API
func SetupRoutes(
	userHandler *UserHandler,
	characterHandler *CharacterHandler,
	guideHandler *GuideHandler,
	eventHandler *EventHandler,
	authHandler *AuthHandler,
) *gin.Engine {
	router := gin.Default()

	// Permissions Middleware
	logged := middleware.RequireAuth()
	creators := middleware.RequireRole(models.RoleVeteran, models.RoleAdmin)
	admin := middleware.RequireRole(models.RoleAdmin)

	// Swagger routes
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
		users.GET("", logged, admin, userHandler.GetUserByEmail)
		// GetUserByID vai desaparecer, mas fica ai por enquanto
		users.GET("/:id", userHandler.GetUserById)
		users.PATCH("/:id", logged, userHandler.UpdateUser)
		users.DELETE("/:id", logged, admin, userHandler.DeleteUserById)
	}

	// Rotas de Personagens
	characters := router.Group("/characters")
	{
		characters.POST("", logged, creators, characterHandler.CreateCharacter)
		characters.GET("/:id", characterHandler.GetCharacterById)
		characters.PATCH("/:id", logged, creators, characterHandler.UpdateCharacter)
		characters.DELETE("/:id", logged, admin, characterHandler.DeleteCharacterById)
	}

	// Rotas de Guias
	guides := router.Group("/guides")
	{
		guides.POST("", logged, creators, guideHandler.CreateGuide)
		guides.GET("", guideHandler.GetGuides)
		guides.GET("/:id", guideHandler.GetGuideById)
		guides.PATCH("/:id", logged, creators, guideHandler.UpdateGuide)
		guides.DELETE("/:id", logged, admin, guideHandler.DeleteGuideById)
	}

	// Rotas de Eventos
	events := router.Group("/events")
	{
		events.POST("", logged, creators, eventHandler.CreateEvent)
		events.GET("/:id", eventHandler.GetEventById)
		events.PATCH("/:id", logged, creators, eventHandler.UpdateEvent)
		events.DELETE("/:id", logged, admin, eventHandler.DeleteEventById)
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
