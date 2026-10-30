package app

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"gorm.io/gorm"

	_ "github.com/kisalto/Feel-The-Night-Swagger/docs"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/auth"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/character"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/event"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/guide"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/middleware"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/user"
)

// NewRouter wires the API: global middlewares, services, handlers and routes.
func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.Timeout(10 * time.Second))

	if os.Getenv("APP_ENV") == "production" {
		_ = router.SetTrustedProxies([]string{ /* Render proxy IP ranges */ })
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	// Permission middlewares (created once and passed to the features)
	logged := middleware.RequireAuth()
	creators := middleware.RequireRole(models.RoleVeteran, models.RoleAdmin)
	admin := middleware.RequireRole(models.RoleAdmin)

	// Login and sign-up share the same rate limit bucket
	authLimiter := middleware.NewIPRateLimiter(rate.Every(12*time.Second), 5).Middleware()

	// General routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", HealthCheckHandler)

	// Services
	userService := user.NewService(db)
	characterService := character.NewService(db)
	guideService := guide.NewService(db)
	eventService := event.NewService(db)

	// Features
	auth.RegisterRoutes(router, auth.NewHandler(userService), authLimiter)
	user.RegisterRoutes(router, user.NewHandler(userService), authLimiter, logged, admin)
	character.RegisterRoutes(router, character.NewHandler(characterService), logged, creators, admin)
	guide.RegisterRoutes(router, guide.NewHandler(guideService), logged, creators, admin)
	event.RegisterRoutes(router, event.NewHandler(eventService), logged, creators, admin)

	return router
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
