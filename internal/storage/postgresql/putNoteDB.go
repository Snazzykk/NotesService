package postgresql

import (
	"NotesService/internal/models"
	"NotesService/internal/storage/storageErr"
	"database/sql"
	"fmt"
	"time"
)

func (s *Storage) PutNote(idUser int64, idNote int64, title string, content string) (*models.Note, error) {
	const op = "storage.postgresql.PutNote"

	note := &models.Note{
		UserID:    idUser,
		ID:        idNote,
		Title:     title,
		Content:   content,
		UpdatedAt: time.Now(),
	}

	err := s.db.QueryRow(`UPDATE notes 
								SET title=$3,
								    content=$4,
								    updated_at=CURRENT_TIMESTAMP 
								WHERE user_id = $1 AND id = $2
								RETURNING id,user_id,title,content,created_at,updated_at`, idUser, idNote, title, content).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, fmt.Errorf("%s: %w", op, storageErr.ErrNoteNotFound)
		}
	}

	return note, nil

}
