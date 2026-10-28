package database

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
)

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
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Character{},
		&models.Event{},
		&models.Guide{},
		&models.LastEvent{},
	); err != nil {
		return nil, fmt.Errorf("falha ao executar o AutoMigrate: %w", err)
	}

	if err := createConstraints(db); err != nil {
		return nil, fmt.Errorf("falha ao criar constraints: %w", err)
	}

	return db, nil
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
				return fmt.Errorf("erro criando constraint %s.%s: %w", fmt.Sprintf("%T", f.model), f.field, err)
			}
		}
	}
	return nil
}

// EscapeLike escapa \, % e _ para usar texto do usuário com LIKE/ILIKE
// sem que esses caracteres virem curingas.
func EscapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
