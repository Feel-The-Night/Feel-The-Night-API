package character

import (
	"context"
	"errors"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/apperr"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/database"
	"gorm.io/gorm"
)

var (
	ErrExists    = errors.New("character already exists")
	ErrHasGuides = errors.New("character has linked guides")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// nameTaken reports whether the name is already used by another character (other than exceptID).
// Case-insensitive: "Ryu" and "ryu" are the same character.
func (s *Service) nameTaken(ctx context.Context, name string, exceptID uint) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.Character{}).
		Where("LOWER(name) = LOWER(?) AND character_id <> ?", name, exceptID).
		Count(&count).Error
	return count > 0, err
}

func (s *Service) CreateCharacter(ctx context.Context, character *models.Character) error {
	taken, err := s.nameTaken(ctx, character.Name, 0)
	if err != nil {
		return err
	}
	if taken {
		return ErrExists
	}

	if err := s.db.WithContext(ctx).Create(character).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrExists
		}
		return err
	}
	return nil
}

func (s *Service) GetCharacterById(ctx context.Context, id uint) (*models.Character, error) {
	var character models.Character
	if err := s.db.WithContext(ctx).First(&character, id).Error; err != nil {
		return nil, err
	}
	return &character, nil
}

func (s *Service) DeleteCharacterById(ctx context.Context, id uint) error {
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&models.Guide{}).
		Where("character_id = ?", id).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrHasGuides
	}

	result := s.db.WithContext(ctx).Delete(&models.Character{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *Service) UpdateCharacterById(ctx context.Context, id uint, character *models.Character) (*models.Character, error) {
	existingCharacter, err := s.GetCharacterById(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if character.Name != "" && character.Name != existingCharacter.Name {
		taken, err := s.nameTaken(ctx, character.Name, id)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, ErrExists
		}
		updates["name"] = character.Name
	}
	if character.Description != "" {
		updates["description"] = character.Description
	}
	if character.Type != "" {
		updates["type"] = character.Type
	}
	if character.ImageURL != "" {
		updates["image_url"] = character.ImageURL
	}

	if len(updates) == 0 {
		return nil, apperr.ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).Model(existingCharacter).Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrExists
		}
		return nil, err
	}

	return existingCharacter, nil
}

func (s *Service) GetCharacters(ctx context.Context, name, charType string) ([]models.Character, error) {
	query := s.db.WithContext(ctx).Model(&models.Character{})
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+database.EscapeLike(name)+"%")
	}
	if charType != "" {
		query = query.Where("LOWER(type) = LOWER(?)", charType)
	}

	var characters []models.Character
	err := query.Order("name ASC").Find(&characters).Error
	return characters, err
}
