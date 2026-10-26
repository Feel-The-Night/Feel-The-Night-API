package dto

// ErrorResponse é o formato padrão para respostas de erro da API.
// @Title ErrorResponse
type ErrorResponse struct {
	ErrorName string `json:"error_name" example:"recurso não encontrado"`
	Error     string `json:"error" example:"o recurso solicitado não foi encontrado"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
