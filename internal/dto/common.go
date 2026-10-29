package dto

// ErrorResponse é o formato padrão para respostas de erro da API.
// @Title ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"o recurso solicitado não foi encontrado"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// PageResponse serve para qualquer listagem paginada (guias, eventos, personagens).
type PageResponse[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
