package services

import (
	"context"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/database"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"gorm.io/gorm"
)

type EventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

func (s *EventService) CreateEvent(ctx context.Context, event *models.Event) error {
	return s.db.WithContext(ctx).Create(event).Error
}

func (s *EventService) GetEventById(ctx context.Context, id uint) (*models.Event, error) {
	var event models.Event
	if err := s.db.WithContext(ctx).First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *EventService) DeleteEventById(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Event{EventID: id})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *EventService) UpdateEventById(ctx context.Context, id uint, event *models.Event) (*models.Event, error) {
	existingEvent, err := s.GetEventById(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if event.Title != "" {
		updates["title"] = event.Title
	}
	if event.BannerURL != "" {
		updates["banner_url"] = event.BannerURL
	}
	if !event.Day.IsZero() {
		updates["day"] = event.Day
	}
	if event.Description != "" {
		updates["description"] = event.Description
	}

	if len(updates) == 0 {
		return nil, ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).Model(existingEvent).Updates(updates).Error; err != nil {
		return nil, err
	}

	return existingEvent, nil
}

func (s *EventService) GetEvents(ctx context.Context, title string, upcoming bool, page, limit int) ([]models.Event, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.Event{})
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+database.EscapeLike(title)+"%")
	}
	if upcoming {
		query = query.Where("day >= CURRENT_DATE")
	}
	query = query.Session(&gorm.Session{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var events []models.Event
	if err := query.
		Order("day ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}
