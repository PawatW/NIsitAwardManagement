package exceptions

import "errors"

var (
	ErrAwardCategoryNotFound      = errors.New("award category not found")
	ErrAwardCategoryAlreadyExists = errors.New("award category already exists")
	ErrInvalidAwardCategoryName   = errors.New("invalid award category name")
	ErrAwardCategoryInUse         = errors.New("award category is in use and cannot be deleted")
)