package postgresql

import (
	"NotesService/internal/models"
	"fmt"
	"strconv"
	"strings"
)

func (s *Storage) GetAllNotes(idUser int64, limit string, offset string, sort string) ([]*models.Note, error) {
	const op = "storage.postgresql.GetAllNotes"

	limitDefault := 10
	offsetDefault := 0

	notes := []*models.Note{}

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err == nil && l > 0 {
			limitDefault = l
		}
	}

	if offset != "" {
		of, err := strconv.Atoi(offset)
		if err == nil && of >= 0 {
			offsetDefault = of
		}
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	query := fmt.Sprintf(`
	SELECT id, user_id, title, content, created_at, updated_at
    FROM notes
    WHERE user_id = $1
    ORDER BY created_at %s
    LIMIT $2
    OFFSET $3
`, strings.ToUpper(sort))

	rows, err := s.db.Query(query, idUser, limitDefault, offsetDefault)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {

		note := &models.Note{}

		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration: %w", op, err)
	}

	return notes, nil
}
