//go:build integration

package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/app"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/database"
)

var (
	router *gin.Engine
	db     *gorm.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	pg, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("ftn_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("starting postgres: %v", err)
	}

	host, _ := pg.Host(ctx)
	port, _ := pg.MappedPort(ctx, "5432/tcp")

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_NAME", "ftn_test")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("JWT_SECRET", "a-test-secret-with-more-than-32-chars")

	db, err = database.Connect() // runs the Goose migrations inside the container
	if err != nil {
		log.Fatalf("connecting: %v", err)
	}
	router = app.NewRouter(db)

	code := m.Run()
	_ = pg.Terminate(ctx)
	os.Exit(code)
}
