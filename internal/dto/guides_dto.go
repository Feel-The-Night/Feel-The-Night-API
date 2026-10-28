package dto

import (
	"time"
)

// CreateGuideInput representa os dados para criação de um guia.
// O autor (UserID) vem do token, não do body.
type CreateGuideInput struct {
	Title       string `json:"title" binding:"required,max=20"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"omitempty,max=15"`
	Description string `json:"description" binding:"required,max=50"`
	Link        string `json:"link" binding:"omitempty,url"`
	CharacterID uint   `json:"character_id" binding:"required,min=1"`
}

// GuideIDInput valida o parâmetro :id na URL.
type GuideIDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// GuideListQuery são os filtros e a paginação de GET /guides.
type GuideListQuery struct {
	Title string `form:"title" binding:"omitempty,max=100"`
	Page  int    `form:"page,default=1" binding:"min=1"`
	Limit int    `form:"limit,default=20" binding:"min=1,max=100"`
}

// UpdateGuideInput representa os campos opcionais para atualização (PATCH) de guia.
type UpdateGuideInput struct {
	Title       string `json:"title" binding:"omitempty,max=20"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"omitempty,max=15"`
	Description string `json:"description" binding:"omitempty,max=50"`
	Link        string `json:"link" binding:"omitempty,url"`
}

// GuideResponse representa os dados do guia retornados nas respostas HTTP.
type GuideResponse struct {
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

type GuidePageResponse struct {
	Data  []GuideResponse `json:"data"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Total int64           `json:"total"`
}
