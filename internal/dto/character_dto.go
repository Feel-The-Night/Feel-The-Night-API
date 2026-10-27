package dto

// CreateCharacterInput representa os dados para criação de um personagem.
// @Title CreateCharacterInput
type CreateCharacterInput struct {
	Name        string `json:"name" binding:"required,max=20"`
	Description string `json:"description" binding:"required,max=255"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"required,max=15"`
}

// CharacterIDInput valida o parâmetro :id na URL.
// @Title CharacterIDInput
type CharacterIDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UpdateCharacterInput representa os campos opcionais para atualização (PATCH) de personagem.
// @Title UpdateCharacterInput
type UpdateCharacterInput struct {
	Name        string `json:"name" binding:"omitempty,max=20"`
	Description string `json:"description" binding:"omitempty,max=255"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"omitempty,max=15"`
}

// CharacterResponse representa os dados do personagem retornados nas respostas HTTP.
// @Title CharacterResponse
type CharacterResponse struct {
	CharacterID uint   `json:"character_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url,omitempty"`
	Type        string `json:"type"`
}

// CharacterListQuery são os filtros e a paginação de GET /events.
type CharacterListQuery struct {
	Name     string `form:"name" binding:"omitempty,max=100"`
	Upcoming bool   `form:"upcoming"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	Limit    int    `form:"limit,default=20" binding:"min=1,max=100"`
}
