package event

// CreateInput holds the data to create an event.
// The author (UserID) comes from the token, not from the body.
type CreateInput struct {
	Title       string `json:"title" binding:"required,max=75"`
	Description string `json:"description" binding:"required,max=255"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url,max=255"`
	Day         string `json:"day" binding:"required" example:"2026-11-15"` // YYYY-MM-DD
}

// IDInput validates the :id URL parameter.
type IDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UpdateInput holds the optional fields for an event update (PATCH).
type UpdateInput struct {
	Title       string `json:"title" binding:"omitempty,max=75"`
	Description string `json:"description" binding:"omitempty,max=255"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url,max=255"`
	Day         string `json:"day" binding:"omitempty" example:"2026-11-15"` // YYYY-MM-DD
}

// ListQuery holds the filters and pagination for GET /events.
type ListQuery struct {
	Title    string `form:"title" binding:"omitempty,max=100"`
	Upcoming bool   `form:"upcoming"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	Limit    int    `form:"limit,default=20" binding:"min=1,max=100"`
}

// Response is the event representation returned by the API.
type Response struct {
	EventID     uint   `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BannerURL   string `json:"banner_url,omitempty"`
	Day         string `json:"day" example:"2026-11-15"` // YYYY-MM-DD
	UserID      uint   `json:"user_id"`
}
