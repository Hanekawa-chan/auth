package app

import "errors"

// TODO стоит добавить grpc ошибки наверное как-то
var (
	ErrInternal      = errors.New("internal error")
	ErrEmptyRequired = errors.New("required value is empty")
	ErrValidation    = errors.New("variable didn't pass validation")
	ErrType          = errors.New("wrong type")
	ErrNotFound      = errors.New("rows not found")
)
