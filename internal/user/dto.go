package user

import "time"

// CreateInput representa os dados necessários para criar um usuário.
// @Title CreateInput
type CreateInput struct {
	Nickname  string `json:"nickname" binding:"required,min=3,max=30"`
	Email     string `json:"email" binding:"required,email"`
	DiscordID string `json:"discord_id"`
	Password  string `json:"password" binding:"required,min=6,max=72"`
}

// IDInput valida o parâmetro :id na URL.
// @Title IDInput
type IDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// EmailInput valida o body email na requisição.
// @Title EmailInput
type EmailInput struct {
	Email string `form:"email" binding:"required,email"`
}

// UpdateInput permite atualização parcial de dados do usuário.
// @Title UpdateInput
type UpdateInput struct {
	Nickname  string `json:"nickname" binding:"omitempty,min=3,max=30"`
	Email     string `json:"email" binding:"omitempty,email"`
	DiscordID string `json:"discord_id"`
	Password  string `json:"password" binding:"omitempty,min=6,max=72"`
}

// Response representa a resposta pública do usuário (sem a senha).
// @Title Response
type Response struct {
	UserID           uint      `json:"user_id"`
	Nickname         string    `json:"nickname"`
	Email            string    `json:"email"`
	DiscordID        string    `json:"discord_id,omitempty"`
	RegistrationDate time.Time `json:"registration_date"`
	EventCount       int       `json:"event_count"`
	GuideCount       int       `json:"guide_count"`
}

// PublicResponse representa a resposta pública do usuário (sem a senha).
// @Title PublicResponse
type PublicResponse struct {
	UserID           uint      `json:"user_id"`
	Nickname         string    `json:"nickname"`
	DiscordID        string    `json:"discord_id,omitempty"`
	RegistrationDate time.Time `json:"registration_date"`
	EventCount       int       `json:"event_count"`
	GuideCount       int       `json:"guide_count"`
}
