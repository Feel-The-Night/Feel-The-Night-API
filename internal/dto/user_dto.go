package dto

import "time"

// CreateUserInput representa os dados necessários para criar um usuário.
// @Title CreateUserInput
type CreateUserInput struct {
	Nickname  string `json:"nickname" binding:"required,min=3,max=30"`
	Email     string `json:"email" binding:"required,email"`
	DiscordID string `json:"discord_id"`
	Password  string `json:"password" binding:"required,min=6,max=72"`
}

// UserIDInput valida o parâmetro :id na URL.
// @Title UserIDInput
type UserIDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UserEmailInput valida o body email na requisição.
// @Title UserEmailInput
type UserEmailInput struct {
	Email string `form:"email" binding:"required,email"`
}

// UpdateUserInput permite atualização parcial de dados do usuário.
// @Title UpdateUserInput
type UpdateUserInput struct {
	Nickname  string `json:"nickname" binding:"omitempty,min=3,max=30"`
	Email     string `json:"email" binding:"omitempty,email"`
	DiscordID string `json:"discord_id"`
	Password  string `json:"password" binding:"omitempty,min=6,max=72"`
}

// UserResponse representa a resposta pública do usuário (sem a senha).
// @Title UserResponse
type UserResponse struct {
	UserID           uint      `json:"user_id"`
	Nickname         string    `json:"nickname"`
	Email            string    `json:"email"`
	DiscordID        string    `json:"discord_id,omitempty"`
	RegistrationDate time.Time `json:"registration_date"`
	EventCount       int       `json:"event_count"`
	GuideCount       int       `json:"guide_count"`
}

// PublicUserResponse representa a resposta pública do usuário (sem a senha).
// @Title PublicUserResponse
type PublicUserResponse struct {
	UserID           uint      `json:"user_id"`
	Nickname         string    `json:"nickname"`
	DiscordID        string    `json:"discord_id,omitempty"`
	RegistrationDate time.Time `json:"registration_date"`
	EventCount       int       `json:"event_count"`
	GuideCount       int       `json:"guide_count"`
}
