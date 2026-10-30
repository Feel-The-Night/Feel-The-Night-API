package database

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		TranslateError:                           true,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(
			log.New(os.Stdout, "", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := createConstraints(db); err != nil {
		return nil, fmt.Errorf("failed to create constraints: %w", err)
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, "migrations")

}

func createConstraints(db *gorm.DB) error {
	m := db.Migrator()

	type fk struct {
		model any
		field string
	}

	fks := []fk{
		{&models.User{}, "Guides"},      // guides.user_id → users
		{&models.Character{}, "Guides"}, // guides.character_id → characters
		{&models.User{}, "Events"},      // events.user_id → users
		{&models.Event{}, "LastEvents"}, // last_events.event_id → events
	}

	for _, f := range fks {
		if !m.HasConstraint(f.model, f.field) {
			if err := m.CreateConstraint(f.model, f.field); err != nil {
				return fmt.Errorf("failed to create constraint %s.%s: %w", fmt.Sprintf("%T", f.model), f.field, err)
			}
		}
	}
	return nil
}

// EscapeLike escapes \, % and _ so user input can be used with LIKE/ILIKE
// without those characters acting as wildcards.
func EscapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
