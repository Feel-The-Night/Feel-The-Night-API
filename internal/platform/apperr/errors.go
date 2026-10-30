package apperr

import "errors"

var (
	ErrNoUpdateData = errors.New("no data provided for update")
)
