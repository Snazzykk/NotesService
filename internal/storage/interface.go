package storage

import "NotesService/internal/models"

type NoteStorage interface {
	SaveNotes(title string, content string, idUser int64) (*models.Note, int64, error)
	GetAllNotes(idUser int64, limit, offset, sort string) ([]*models.Note, error)
	GetOneNote(idUser int64, idNote int64) (*models.Note, error)
	PutNote(idUser int64, idNote int64, title string, content string) (*models.Note, error)
	DeleteNote(idUser int64, idNote int64) error
}

type UserStorage interface {
	RegisterUser(userName string) (*models.User, error)
}
