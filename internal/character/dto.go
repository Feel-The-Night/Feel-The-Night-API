package character

// CreateInput representa os dados para criação de um personagem.
// @Title CreateInput
type CreateInput struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=255"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"required,max=30"`
}

// IDInput valida o parâmetro :id na URL.
// @Title IDInput
type IDInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// UpdateInput representa os campos opcionais para atualização (PATCH) de personagem.
// @Title UpdateInput
type UpdateInput struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	Type        string `json:"type" binding:"omitempty,max=30"`
}

// Response representa os dados do personagem retornados nas respostas HTTP.
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
