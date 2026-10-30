package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/apperr"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailTaken = errors.New("email already registered")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashed), nil
}

// emailTaken reports whether the email is already used by another user (other than exceptID).
func (s *Service) emailTaken(ctx context.Context, email string, exceptID uint) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ? AND user_id <> ?", email, exceptID).
		Count(&count).Error
	return count > 0, err
}

// CreateUser creates a new user.
// The client only controls nickname, email, discord_id and password;
// ID, role, counters and timestamps are always set here or by the database.
func (s *Service) CreateUser(ctx context.Context, user *models.User) error {
	user.Email = normalizeEmail(user.Email)

	taken, err := s.emailTaken(ctx, user.Email, 0)
	if err != nil {
		return err
	}
	if taken {
		return ErrEmailTaken
	}

	hashed, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	user.Role = models.RoleBeginner
	user.EventCount = 0
	user.GuideCount = 0

	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		// race condition: two simultaneous sign-ups passed the check
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrEmailTaken
		}
		return err
	}
	return nil
}

func (s *Service) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).
		Where("email = ?", normalizeEmail(email)).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) DeleteUserById(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *Service) UpdateUser(ctx context.Context, id uint, user *models.User) (*models.User, error) {
	existingUser, err := s.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if user.Nickname != "" {
		updates["nickname"] = user.Nickname
	}
	if user.DiscordID != "" {
		updates["discord_id"] = user.DiscordID
	}

	if user.Email != "" {
		email := normalizeEmail(user.Email)
		if email != existingUser.Email {
			taken, err := s.emailTaken(ctx, email, id)
			if err != nil {
				return nil, err
			}
			if taken {
				return nil, ErrEmailTaken
			}
			updates["email"] = email
		}
	}

	if user.Password != "" {
		hashed, err := hashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		updates["password"] = hashed
	}

	if len(updates) == 0 {
		return nil, apperr.ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).Model(existingUser).Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	return existingUser, nil
}
