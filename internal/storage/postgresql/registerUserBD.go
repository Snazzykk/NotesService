package postgresql

import (
	"NotesService/internal/models"
	"fmt"
)

func (s *Storage) RegisterUser(userName string) (*models.User, error) {
	const op = "storage.postgresql.RegisterUser"

	user := &models.User{
		Username: userName,
	}

	err := s.db.QueryRow(`INSERT INTO users (user_name) 
									  values ($1)
									  RETURNING id,user_name,created_at`, userName).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil

}
