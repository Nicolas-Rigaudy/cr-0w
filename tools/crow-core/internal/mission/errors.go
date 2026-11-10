package mission

import (
	"errors"
	"fmt"
)

// Error kinds of the owner actions (item_147). Clients map them to HTTP codes: ErrNotFound
// 404, ErrInvalid 400, ErrConflict 409. The message stays the one the CLI prints.
var (
	ErrInvalid  = errors.New("invalid request")
	ErrConflict = errors.New("conflict")
)

type kindError struct {
	kind error
	msg  string
}

func (e *kindError) Error() string { return e.msg }
func (e *kindError) Unwrap() error { return e.kind }

// failf returns an error of kind with a formatted message.
func failf(kind error, format string, a ...any) error {
	return &kindError{kind: kind, msg: fmt.Sprintf(format, a...)}
}
