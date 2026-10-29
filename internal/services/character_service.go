package services

import (
	"context"
	"errors"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/database"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"gorm.io/gorm"
)

type CharacterService struct {
	db *gorm.DB
}

func NewCharacterService(db *gorm.DB) *CharacterService {
	return &CharacterService{db: db}
}

// nameTaken verifica se o nome já é usado por outro personagem (exceto exceptID).
// Compara sem diferenciar maiúsculas: "Ryu" e "ryu" são o mesmo personagem.
func (s *CharacterService) nameTaken(ctx context.Context, name string, exceptID uint) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.Character{}).
		Where("LOWER(name) = LOWER(?) AND character_id <> ?", name, exceptID).
		Count(&count).Error
	return count > 0, err
}

func (s *CharacterService) CreateCharacter(ctx context.Context, character *models.Character) error {
	taken, err := s.nameTaken(ctx, character.Name, 0)
	if err != nil {
		return err
	}
	if taken {
		return ErrCharacterExists
	}

	if err := s.db.WithContext(ctx).Create(character).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrCharacterExists
		}
		return err
	}
	return nil
}

func (s *CharacterService) GetCharacterById(ctx context.Context, id uint) (*models.Character, error) {
	var character models.Character
	if err := s.db.WithContext(ctx).First(&character, id).Error; err != nil {
		return nil, err
	}
	return &character, nil
}

func (s *CharacterService) DeleteCharacterById(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Character{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *CharacterService) UpdateCharacterById(ctx context.Context, id uint, character *models.Character) (*models.Character, error) {
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
			return nil, ErrCharacterExists
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
		return nil, ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).Model(existingCharacter).Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrCharacterExists
		}
		return nil, err
	}

	return existingCharacter, nil
}

func (s *CharacterService) GetCharacters(ctx context.Context, name string, page, limit int) ([]models.Character, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.Character{})
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+database.EscapeLike(name)+"%")
	}

	query = query.Session(&gorm.Session{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var characters []models.Character

	if err := query.
		Order("name").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&characters).Error; err != nil {
		return nil, 0, err
	}

	return characters, total, nil
}
