package postgresql

import (
	"NotesService/internal/models"
	"fmt"
)

func (s *Storage) SaveNotes(title string, content string, idUser int64) (*models.Note, int64, error) {
	const op = "storage.postgresql.SaveNotes"

	if title == "" {
		return nil, 0, fmt.Errorf("%s: title cannot be empty", op)
	}
	if content == "" {
		return nil, 0, fmt.Errorf("%s: content cannot be empty", op)
	}

	stmt, err := s.db.Prepare(`insert into notes (user_id,title,content) values ($1,$2,$3) returning id,user_id,title,content,created_at,updated_at`)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	note := &models.Note{
		UserID:  idUser,
		Title:   title,
		Content: content,
	}
	var id int64

	err = stmt.QueryRow(idUser, title, content).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	return note, id, nil
}
