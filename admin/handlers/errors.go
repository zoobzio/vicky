package handlers

import "github.com/zoobzio/rocco"

var (
	// ErrUserNotFound indicates the requested user does not exist.
	ErrUserNotFound = rocco.ErrNotFound.WithMessage("user not found")

	// ErrInvalidEmail indicates the email format is invalid.
	ErrInvalidEmail = rocco.ErrBadRequest.WithMessage("invalid email format")

	// ErrDuplicateLogin indicates the login already exists.
	ErrDuplicateLogin = rocco.ErrConflict.WithMessage("login already exists")

	// ErrInvalidLimit indicates the limit parameter is out of bounds.
	ErrInvalidLimit = rocco.ErrBadRequest.WithMessage("limit must be between 1 and 100")

	// ErrRepositoryNotFound indicates the requested repository does not exist.
	ErrRepositoryNotFound = rocco.ErrNotFound.WithMessage("repository not found")

	// ErrInvalidUserID indicates the user_id filter parameter is invalid.
	ErrInvalidUserID = rocco.ErrBadRequest.WithMessage("invalid user_id parameter")

	// ErrJobNotFound indicates the requested job does not exist.
	ErrJobNotFound = rocco.ErrNotFound.WithMessage("job not found")

	// ErrJobNotCancellable indicates the job cannot be cancelled in its current state.
	ErrJobNotCancellable = rocco.ErrBadRequest.WithMessage("job cannot be cancelled (already completed, failed, or cancelled)")
)
