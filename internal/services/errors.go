package services

import "errors"

var (
	ErrNoUpdateData       = errors.New("nenhum dado recebido para atualização")
	ErrEmailTaken         = errors.New("email já cadastrado")
	ErrCharacterExists    = errors.New("personagem já existe")
	ErrCharacterHasGuides = errors.New("Personagem tem guias")
	ErrCharacterNotFound  = errors.New("Personagem nao encontrado")
)
