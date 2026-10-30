package character

// CreateInput holds the data to create a character.
// @Title CreateInput
type CreateInput struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=255"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"required,max=30"`
}

// IDInput validates the :id URL parameter.
// @Title IDInput
type IDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UpdateInput holds the optional fields for a character update (PATCH).
// @Title UpdateInput
type UpdateInput struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"omitempty,max=30"`
}

// Response is the character representation returned by the API.
// @Title Response
type Response struct {
	CharacterID uint   `json:"character_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url,omitempty"`
	Type        string `json:"type"`
}

type ListQuery struct {
	Name string `form:"name" binding:"omitempty,max=100"`
	Type string `form:"type" binding:"omitempty,max=30"`
}
