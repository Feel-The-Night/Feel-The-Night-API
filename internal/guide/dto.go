package guide

import (
	"time"
)

// CreateInput representa os dados para criação de um guia.
// O autor (UserID) vem do token, não do body.
type CreateInput struct {
	Title       string `json:"title" binding:"required,max=254"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"omitempty,max=15"`
	Description string `json:"description" binding:"required,max=2083"`
	Link        string `json:"link" binding:"omitempty,url"`
	CharacterID uint   `json:"character_id" binding:"required,min=1"`
}

// IDInput valida o parâmetro :id na URL.
type IDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// ListQuery são os filtros e a paginação de GET /guides.
type ListQuery struct {
	Title string `form:"title" binding:"omitempty,max=100"`
	Page  int    `form:"page,default=1" binding:"min=1"`
	Limit int    `form:"limit,default=20" binding:"min=1,max=100"`
}

// UpdateInput representa os campos opcionais para atualização (PATCH) de guia.
type UpdateInput struct {
	Title       string `json:"title" binding:"omitempty,max=254"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"omitempty,max=15"`
	Description string `json:"description" binding:"omitempty,max=2083"`
	Link        string `json:"link" binding:"omitempty,url"`
}

// Response representa os dados do guia retornados nas respostas HTTP.
type Response struct {
	GuideID      uint      `json:"guide_id"`
	Title        string    `json:"title"`
	BannerURL    string    `json:"banner_url,omitempty"`
	Type         string    `json:"type,omitempty"`
	Description  string    `json:"description"`
	Link         string    `json:"link,omitempty"`
	CreationDate time.Time `json:"creation_date"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	UserID       uint      `json:"user_id"`
	CharacterID  uint      `json:"character_id"`
}
