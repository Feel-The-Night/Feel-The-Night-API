package dto

// ErrorResponse é o formato padrão para respostas de erro da API.
// @Title ErrorResponse
type ErrorResponse struct {
	ErrorName string `json:"error_name" example:"recurso não encontrado"`
	Error     string `json:"error" example:"o recurso solicitado não foi encontrado"`
}
