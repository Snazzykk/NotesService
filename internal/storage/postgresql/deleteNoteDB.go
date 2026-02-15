package postgresql

import (
	"NotesService/internal/storage/storageErr"
	"fmt"
)

func (s *Storage) DeleteNote(idUser int64, idNote int64) error {
	const op = "storage.postgresql.DeleteNote"

	res, err := s.db.Exec(`DELETE FROM notes 
								WHERE user_id = $1 AND id = $2`, idUser, idNote)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	//res.RowsAffected() - возвращает количество затронутых строк
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	// если rowsAffected = 1, то было что то удалено, если 0, то ошибка
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storageErr.ErrNoteNotFound)
	}

	return nil

}
