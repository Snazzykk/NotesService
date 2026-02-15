package auth

import (
	"errors"
	"fmt"
	"time"

	"NotesService/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager управляет генерацией и проверкой токенов
type JWTManager struct {
	secret   []byte        // Секретный ключ (хранится в памяти, не в коде!)
	duration time.Duration // Время жизни токена
}

// NewJWTManager создаёт новый менеджер токенов
// Секрет передаётся извне (из main.go) для безопасности
func NewJWTManager(secret string, duration time.Duration) (*JWTManager, error) {
	if secret == "" {
		return nil, errors.New("secret key cannot be empty")
	}

	if duration <= 0 {
		return nil, errors.New("token duration must be positive")
	}

	return &JWTManager{
		secret:   []byte(secret),
		duration: duration,
	}, nil
}

// GenerateToken создаёт новый JWT токен для пользователя
func (m *JWTManager) GenerateToken(userID int64, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["user_name"] = username
	claims["exp"] = time.Now().Add(m.duration).Unix()
	claims["iat"] = time.Now().Unix() // время создания

	// Подписываем токен секретным ключом из памяти
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken проверяет и валидирует токен
func (m *JWTManager) VerifyToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Извлекаем данные из токена
	userID, ok := claims["id"].(float64) // JSON числа → float64 в Go
	if !ok {
		return nil, errors.New("invalid user id in token")
	}

	username, ok := claims["user_name"].(string)
	if !ok {
		return nil, errors.New("invalid username in token")
	}

	return &models.User{
		ID:       int64(userID),
		Username: username,
	}, nil
}
