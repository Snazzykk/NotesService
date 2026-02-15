package postgresql

import (
	"NotesService/internal/models"
	"NotesService/internal/storage/storageErr"
	"database/sql"
	"fmt"
)

func (s *Storage) GetOneNote(idUser int64, idNote int64) (*models.Note, error) {
	const op = "storage.postgresql.GetOneNote"

	row := s.db.QueryRow(`SELECT id, user_id, title, content, created_at, updated_at
									  FROM notes
									  Where user_id = $1 AND id = $2`, idUser, idNote)

	note := &models.Note{}

	err := row.Scan(
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
