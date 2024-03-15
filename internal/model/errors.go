package model

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input body or query params")
	ErrValidationError = errors.New("given struct is invalid")

	ErrMovieNotExists = errors.New("movie with required id does not exist")
	ErrActorNotExists = errors.New("actor with required id does not exist")

	ErrUserNotExists = errors.New("user with required id does not exist")
	ErrUnauthorized  = errors.New("authorization header with user id is missing")

	ErrPermissionDenied = errors.New("user with required id does not have permission for this operation")

	ErrDatabaseError = errors.New("something wrong with database")
	ErrServiceError  = errors.New("unknown error from the service")
)
