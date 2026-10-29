package services

import (
	"context"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/database"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"gorm.io/gorm"
)

type GuideService struct {
	db *gorm.DB
}

func NewGuideService(db *gorm.DB) *GuideService {
	return &GuideService{db: db}
}

func (s *GuideService) CreateGuide(ctx context.Context, guide *models.Guide) error {
	// CreatedAt o GORM preenche sozinho; Likes/Dislikes têm default 0 no model
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&models.Character{}).
		Where("Character_id = ?", guide.CharacterID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrCharacterNotFound
	}

	return s.db.WithContext(ctx).Create(guide).Error

	// return s.db.WithContext(ctx).Create(guide).Error
}

func (s *GuideService) GetGuideById(ctx context.Context, id uint) (*models.Guide, error) {
	var guide models.Guide
	if err := s.db.WithContext(ctx).First(&guide, id).Error; err != nil {
		return nil, err
	}
	return &guide, nil
}

func (s *GuideService) GetGuides(ctx context.Context, title string, page, limit int) ([]models.Guide, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.Guide{})
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+database.EscapeLike(title)+"%")
	}
	// Session deixa a query reutilizável no Count e no Find
	query = query.Session(&gorm.Session{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var guides []models.Guide
	if err := query.
		Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&guides).Error; err != nil {
		return nil, 0, err
	}

	return guides, total, nil
}

func (s *GuideService) DeleteGuideById(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Guide{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *GuideService) UpdateGuide(ctx context.Context, id uint, guide *models.Guide) (*models.Guide, error) {
	existingGuide, err := s.GetGuideById(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)
	if guide.Title != "" {
		updates["title"] = guide.Title
	}
	if guide.BannerURL != "" {
		updates["banner_url"] = guide.BannerURL
	}
	if guide.Type != "" {
		updates["type"] = guide.Type
	}
	if guide.Description != "" {
		updates["description"] = guide.Description
	}
	if guide.Link != "" {
		updates["link"] = guide.Link
	}

	if len(updates) == 0 {
		return nil, ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).Model(existingGuide).Updates(updates).Error; err != nil {
		return nil, err
	}

	return existingGuide, nil
}
