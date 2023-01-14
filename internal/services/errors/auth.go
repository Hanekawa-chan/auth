package errors

import "errors"

var (
	ErrEmptyRequired = errors.New("required value is empty")
	ErrValidation    = errors.New("variable didn't pass validation")
	ErrType          = errors.New("wrong type")
)
