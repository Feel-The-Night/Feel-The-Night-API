package services

import (
	"errors"
	"time"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Change it so userID, registrationDate, isModerator, isVeteran, eventCount and guideCount
// cant be passed througth the createUser post.
// userID will be added automaticaly  by the DB
// registrationDate will get date.now
// isModerator will always be False
// isVeteran will always be false
// eventCount and guideCount will be 0 when creating account
func (s *UserService) CreateUser(user *models.User) error {
	var existingUser models.User

	if err := s.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("Email ja cadastrado")
	}

	user.RegistrationDate = time.Now()
	user.EventCount = 0
	user.GuideCount = 0
	user.Roles = []models.Role{
		{ID: 1},
	}

	return s.db.Create(user).Error
}

func (s *UserService) GetUserById(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUserById(id uint) error {
	result := s.db.Delete(&models.User{UserID: id})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *UserService) UpdateUser(id uint, user *models.User) (*models.User, error) {
	existingUser, err := s.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user.Email != "" && user.Email != existingUser.Email {
		var checkEmail models.User
		err := s.db.Where("email = ? AND user_id != ?", user.Email, id).First(&checkEmail).Error
		if err == nil {
			return nil, errors.New("email já cadastrado por outro usuário")
		}
	}

	updates := make(map[string]any)

	if user.Nickname != "" {
		updates["nickname"] = user.Nickname
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.DiscordID != "" {
		updates["discord_id"] = user.DiscordID
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}

	if len(updates) == 0 {
		return existingUser, nil
	}

	if err := s.db.Model(existingUser).Updates(updates).Error; err != nil {
		return nil, err
	}

	return existingUser, nil
}
