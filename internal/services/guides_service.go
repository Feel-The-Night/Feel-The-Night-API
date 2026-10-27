package services

import (
	"errors"
	"time"

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

func (s *GuideService) CreateGuide(guide *models.Guide) error {

	guide.CreatedAt = time.Now()
	guide.Likes = 0
	guide.Dislikes = 0

	return s.db.Create(guide).Error
}

func (s *GuideService) GetGuideById(id uint) (*models.Guide, error) {
	var guide models.Guide
	if err := s.db.First(&guide, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &guide, nil
}

func (s *GuideService) GetGuides(title string, page, limit int) ([]models.Guide, int64, error) {
	query := s.db.Model(&models.Guide{})
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

func (s *GuideService) DeleteGuideById(id uint) error {
	result := s.db.Delete(&models.Guide{GuideID: id})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *GuideService) UpdateGuide(id uint, guide *models.Guide) (*models.Guide, error) {
	existingGuide, err := s.GetGuideById(id)
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
		return existingGuide, nil
	}

	if err := s.db.Model(existingGuide).Updates(updates).Error; err != nil {
		return nil, err
	}

	return existingGuide, nil
}
