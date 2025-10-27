package domain

import "errors"

var (
	// ErrInnovationNotFound is returned when an innovation is not found
	ErrInnovationNotFound = errors.New("innovation not found")

	// ErrAlreadyVoted is returned when a user has already voted for an innovation
	ErrAlreadyVoted = errors.New("already voted for this innovation")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
)


