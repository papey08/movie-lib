package model

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input body or query params")
	ErrValidationError = errors.New("given struct is invalid")

	ErrMovieNotExists = errors.New("movie with required id does not exist")
	ErrActorNotExists = errors.New("actor with required id does not exist")

	ErrUserNotExists = errors.New("user with required id does not exist")

	ErrPermissionDenied = errors.New("user with required id does not have permission for this operation")

	ErrDatabaseError = errors.New("something wrong with database")
)
