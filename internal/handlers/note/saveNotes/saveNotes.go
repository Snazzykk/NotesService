package saveNotes

import (
	resp "NotesService/internal/api/response"
	"NotesService/internal/auth"
	"NotesService/internal/models"
	"NotesService/internal/storage"
	sl "NotesService/pkg/logger/logSlog"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type NoteStorage interface {
	storage.NoteStorage
}

// SaveNotes godoc
// @Summary Create a new note
// @Description Saves a new note for a specific user. Requires JWT authentication.
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Param request body models.SaveNoteRequest true "Note payload"
// @Success 201 {object} models.NoteResponse "Created note"
// @Failure 400
// @Failure 401
// @Failure 500
// @Security ApiKeyAuth
// @Router /users/{id}/notes [post]
func New(log *slog.Logger, saveNotes NoteStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.saveNotes.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authorizedUserID, ok := auth.GetUserID(r)
		if !ok {
			log.Error("user_id not found in context")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Unauthorized"))
			return
		}

		//1.Read body request
		var req models.SaveNoteRequest
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

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Info("Id is empty")
			render.JSON(w, r, resp.Error("Id is empty"))
			render.Status(r, http.StatusBadRequest)
			return
		}

		idUser, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("Failed to convert id to int64", "error", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Invalid id format: must be integer"))
			return
		}

		if authorizedUserID != idUser {
			log.Warn("Unauthorized access attempt",
				slog.Int64("authorized_user_id", authorizedUserID),
				slog.Int64("requested_user_id", idUser),
			)

			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Not found"))
			return
		}

		if req.TitleNote == "" && req.ContentNote == "" || req.TitleNote == " " && req.ContentNote == " " {
			log.Error("The fields cannot be empty.")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("The fields cannot be empty."))
			return
		}

		if req.TitleNote == "" || req.TitleNote == " " {
			log.Error("The field title cannot be empty.")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("The field title cannot be empty."))
			return
		}
		if req.ContentNote == "" || req.ContentNote == " " {
			log.Error("The field content cannot be empty.")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("The field content cannot be empty."))
			return
		}

		Title := strings.TrimSpace(req.TitleNote)
		Content := strings.TrimSpace(req.ContentNote)

		note, _, err := saveNotes.SaveNotes(Title, Content, idUser)
		if err != nil {

			log.Info("Failed to save notes", "error", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to save notes"))
			return
		}

		log.Info("Success", slog.Int64("id", note.ID))

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, models.NoteResponse{
			Response:  resp.Created("Success"),
			NoteID:    note.ID,
			UserId:    note.UserID,
			Title:     note.Title,
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}

}
