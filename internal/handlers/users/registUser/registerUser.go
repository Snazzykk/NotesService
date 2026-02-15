package registUser

import (
	resp "NotesService/internal/api/response"
	"NotesService/internal/models"
	"NotesService/internal/storage"
	sl "NotesService/pkg/logger/logSlog"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type JWTManager interface {
	GenerateToken(userID int64, username string) (string, error)
}

type UserStorage interface {
	storage.UserStorage
}

// RegisterUser godoc
// @Summary Register new user
// @Description Creates a new user and returns user info with JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserRequest true "User registration data"
// @Success 201 {object} models.UserResponse
// @Failure 400
// @Failure 500
// @Router /users [post]
func New(log *slog.Logger, userStorage UserStorage, jwtManager JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.registerUser.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		//1.Read body request
		var req models.UserRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			if err == io.EOF {
				log.Info("Request body is empty (EOF)")
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("Request body cannot be empty"))
				return
			}

			if strings.Contains(err.Error(), "invalid character") {
				log.Info("Invalid JSON format", slog.String("error", err.Error()))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("Invalid JSON format"))
				return
			}

			if strings.Contains(err.Error(), "syntax error") {
				log.Info("JSON syntax error", slog.String("error", err.Error()))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("JSON syntax error"))
				return
			}

			log.Error("Failed to decode request body", "error", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to decode request body"))
			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Failed to validate request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		if req.Username == "" {
			log.Error("The fields cannot be empty.")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("The fields cannot be empty."))
			return
		}

		UserName := strings.TrimSpace(req.Username)

		user, err := userStorage.RegisterUser(UserName)
		if err != nil {
			log.Info("Failed to save user", "error", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to save user"))
			return
		}

		// 4. Генерация JWT токена
		token, err := jwtManager.GenerateToken(user.ID, user.Username)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to generate token"))
			return
		}

		log.Info("Success", slog.Int64("id", user.ID))

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, models.UserResponse{
			Response:  resp.Created("Success"),
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			Token:     token,
		})

	}
}
