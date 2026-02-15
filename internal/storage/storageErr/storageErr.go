package storageErr

import "errors"

var (
	ErrNoteNotFound = errors.New("Note not found")
	ErrUserNotFound = errors.New("User not found")
)
