package dto

// CreateEventInput representa os dados para criação de um evento.
// O autor (UserID) vem do token, não do body.
type CreateEventInput struct {
	Title       string `json:"title" binding:"required,max=75"`
	Description string `json:"description" binding:"required,max=255"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url,max=255"`
	Day         string `json:"day" binding:"required" example:"2026-11-15"` // AAAA-MM-DD
}

// EventIDInput valida o parâmetro :id na URL.
type EventIDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UpdateEventInput representa os campos opcionais para atualização (PATCH) de evento.
type UpdateEventInput struct {
	Title       string `json:"title" binding:"omitempty,max=75"`
	Description string `json:"description" binding:"omitempty,max=255"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url,max=255"`
	Day         string `json:"day" binding:"omitempty" example:"2026-11-15"` // AAAA-MM-DD
}

// EventListQuery são os filtros e a paginação de GET /events.
type EventListQuery struct {
	Title    string `form:"title" binding:"omitempty,max=100"`
	Upcoming bool   `form:"upcoming"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	Limit    int    `form:"limit,default=20" binding:"min=1,max=100"`
}

// EventResponse representa os dados do evento retornados nas respostas HTTP.
type EventResponse struct {
	EventID     uint   `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BannerURL   string `json:"banner_url,omitempty"`
	Day         string `json:"day" example:"2026-11-15"` // AAAA-MM-DD
	UserID      uint   `json:"user_id"`
}
