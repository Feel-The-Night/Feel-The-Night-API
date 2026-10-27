package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
)

var DB *gorm.DB

func Connect() error {
	// Local
	// err := godotenv.Load()
	// if err != nil {
	// 	return fmt.Errorf("erro ao carregar o arquivo .env: %w", err)
	// }

	// Cloud render
	err := godotenv.Load()
	if err != nil {
		// Em vez de retornar erro e derrubar o app, apenas avisa no terminal
		fmt.Println("Aviso: Arquivo .env não encontrado. Lendo variáveis do sistema (Produção).")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		TranslateError:                           true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Character{},
		&models.Event{},
		&models.Guide{},
		&models.LastEvent{},
	)
	if err != nil {
		return fmt.Errorf("falha ao executar o AutoMigrate: %w", err)
	}

	if err := createConstraints(); err != nil {
		return fmt.Errorf("falha ao criar constraints: %w", err)
	}

	fmt.Println("Conexão com o banco de dados e migrações concluídas com sucesso!")
	return nil
}

func createConstraints() error {
	m := DB.Migrator()

	type fk struct {
		model any
		field string
	}

	fks := []fk{
		{&models.Guide{}, "User"},
		{&models.Guide{}, "Character"},
		{&models.Event{}, "User"},
		{&models.LastEvent{}, "Event"},
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

func HasAccess(userID int, accessName string) bool {
	var count int64

	err := DB.Model(&models.User{}).
		Joins("JOIN user_role ur ON ur.user_id = users.id").
		Joins("JOIN role_access ra ON ra.role_id = ur.role_id").
		Joins("JOIN access a ON a.access_id = ra.access_id").
		Where("users.id = ? AND a.access_name = ?", userID, accessName).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

// EscapeLike escapa \, % e _ para usar texto do usuário com LIKE/ILIKE
// sem que esses caracteres virem curingas.
func EscapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
