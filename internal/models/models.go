package models

import (
	"time"

	"gorm.io/gorm"
)

// User table
type User struct {
	UserID           uint      `gorm:"primaryKey"`
	Nickname         string    `gorm:"size:30;not null"`
	Email            string    `gorm:"size:50;not null"`
	DiscordID        string    `gorm:"size:25"`
	Password         string    `gorm:"size:30;not null"`
	RegistrationDate time.Time `gorm:"default:CURRENT_DATE"`
	EventCount       int       `gorm:"default:0"`
	GuideCount       int       `gorm:"default:0"`
	// IsVeteran   bool   `gorm:"default:false"`
	// IsModerator bool   `gorm:"default:false"`

	// Campos de auditoria (substituindo gorm.Model)
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Has many Events and Guides
	Roles  []Role  `gorm:"many2many:user_role"`
	Events []Event `gorm:"constraint:OnUpdate:CASCADE"`
	Guides []Guide `gorm:"constraint:OnUpdate:CASCADE"`
}

func (User) TableName() string { return "users" }

// Character table
type Character struct {
	CharacterID uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:20;not null"`
	Description string `gorm:"size:255;not null"`
	ImageURL    string `gorm:"size:255"`
	Type        string `gorm:"size:15;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Has many guides
	Guides []Guide `gorm:"constraint:OnUpdate:CASCADE"`
}

func (Character) TableName() string { return "characters" }

// Guide table
type Guide struct {
	GuideID      uint      `gorm:"primaryKey"`
	Title        string    `gorm:"size:20;not null"`
	BannerURL    string    `gorm:"size:255"`
	Type         string    `gorm:"size:15"`
	Description  string    `gorm:"size:50;not null"`
	Link         string    `gorm:"size:2083"`
	CreationDate time.Time `gorm:"default:CURRENT_DATE"`
	Likes        int       `gorm:"default:0"`
	Dislikes     int       `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Foreign Keys (Belongs to User and Character)
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CharacterID uint      `gorm:"not null"`
	Character   Character `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Guide) TableName() string { return "guides" }

// Event table
type Event struct {
	EventID     uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:75;not null"`
	Description string    `gorm:"size:255;not null"`
	BannerURL   string    `gorm:"size:255"`
	Day         time.Time `gorm:"type:date;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Foreign Key (Belongs to User)
	UserID uint `gorm:"not null"`
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Has many LastEvents
	LastEvents []LastEvent `gorm:"constraint:OnUpdate:CASCADE"`
}

func (Event) TableName() string { return "events" }

// LastEvent table
type LastEvent struct {
	LastEventID uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:75;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Foreign Key (Belongs to Event)
	EventID uint  `gorm:"not null"`
	Event   Event `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (LastEvent) TableName() string { return "last_events" }

// Access roles
type Role struct {
	ID       int      `gorm:"primaryKey"`
	Accesses []Access `gorm:"many2many:role_access;"`
}

type Access struct {
	AccessID   int `gorm:"primaryKey"`
	AccessName string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"senha"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type ErrorReponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
