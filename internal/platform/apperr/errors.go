package apperr

import "errors"

var (
	ErrNoUpdateData = errors.New("nenhum dado recebido para atualização")
)
