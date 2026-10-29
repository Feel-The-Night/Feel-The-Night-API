package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleBeginner Role = "iniciante"
	RoleVeteran  Role = "veterano"
	RoleAdmin    Role = "admin"
)

// User table
type User struct {
	UserID     uint   `gorm:"primaryKey"`
	Nickname   string `gorm:"size:30;not null"`
	Email      string `gorm:"size:254;not null;unique:idx_users_email,where:deleted_at IS NULL"`
	DiscordID  string `gorm:"size:25"`
	Password   string `gorm:"size:100;not null" json:"-"`
	EventCount int    `gorm:"default:0"`
	GuideCount int    `gorm:"default:0"`
	//RegistrationDate time.Time `gorm:"default:CURRENT_DATE"`
	// IsVeteran   bool   `gorm:"default:false"`
	// IsModerator bool   `gorm:"default:false"`

	// Campos de auditoria (substituindo gorm.Model)
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Has many Events and Guides
	Role   Role    `gorm:"type:varchar(15);not null;default:'iniciante';index;check:chk_users_role,role IN ('iniciante','veterano','admin')"`
	Events []Event `gorm:"constraint:OnUpdate:CASCADE"`
	Guides []Guide `gorm:"constraint:OnUpdate:CASCADE"`
}

func (User) TableName() string { return "users" }

// Character table
type Character struct {
	CharacterID uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"size:255;not null"`
	ImageURL    string `gorm:"size:255"`
	Type        string `gorm:"size:30;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Has many guides
	Guides []Guide `gorm:"constraint:OnUpdate:CASCADE"`
}

func (Character) TableName() string { return "characters" }

// Guide table
type Guide struct {
	GuideID     uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:254;not null"`
	BannerURL   string `gorm:"size:255"`
	Type        string `gorm:"size:15"`
	Description string `gorm:"size:2083;not null"`
	Link        string `gorm:"size:2083"`
	Likes       int    `gorm:"default:0"`
	Dislikes    int    `gorm:"default:0"`
	//CreationDate time.Time `gorm:"default:CURRENT_DATE"`

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
