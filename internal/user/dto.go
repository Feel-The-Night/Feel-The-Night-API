package user

import "time"

// CreateInput holds the data required to create a user.
// @Title CreateInput
type CreateInput struct {
	Nickname  string `json:"nickname" binding:"required,min=3,max=30"`
	Email     string `json:"email" binding:"required,email"`
	DiscordID string `json:"discord_id"`
	Password  string `json:"password" binding:"required,min=6,max=72"`
}

// IDInput validates the :id URL parameter.
// @Title IDInput
type IDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// EmailInput validates the email query parameter.
// @Title EmailInput
type EmailInput struct {
	Email string `form:"email" binding:"required,email"`
}

// UpdateInput allows a partial update of the user's data.
// @Title UpdateInput
type UpdateInput struct {
	Nickname  string `json:"nickname" binding:"omitempty,min=3,max=30"`
	Email     string `json:"email" binding:"omitempty,email"`
	DiscordID string `json:"discord_id"`
	Password  string `json:"password" binding:"omitempty,min=6,max=72"`
}

// Response is the full user representation (without the password).
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

// PublicResponse is the public user profile (without email or password).
// @Title PublicResponse
type PublicResponse struct {
	UserID           uint      `json:"user_id"`
	Nickname         string    `json:"nickname"`
	DiscordID        string    `json:"discord_id,omitempty"`
	RegistrationDate time.Time `json:"registration_date"`
	EventCount       int       `json:"event_count"`
	GuideCount       int       `json:"guide_count"`
}
