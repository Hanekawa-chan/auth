package app

import "errors"

var (
	ErrInternal      = errors.New("internal error")
	ErrEmptyRequired = errors.New("required value is empty")
	ErrValidation    = errors.New("variable didn't pass validation")
	ErrType          = errors.New("wrong type")
	ErrNotFound      = errors.New("rows not found")
)
