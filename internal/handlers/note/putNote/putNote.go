package putNote

import (
	resp "NotesService/internal/api/response"
	"NotesService/internal/auth"
	"NotesService/internal/models"
	"NotesService/internal/storage"
	"NotesService/internal/storage/storageErr"
	sl "NotesService/pkg/logger/logSlog"
	"errors"
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

// PutNote godoc
// @Summary Update a note by ID
// @Description Updates the title and/or content of a note for a specific user. Requires JWT authentication.
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Param note_id path int true "Note ID" minimum(1)
// @Param request body models.PutNoteRequest true "Note update payload"
// @Success 200 {object} models.NoteResponse "Updated note"
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Security ApiKeyAuth
// @Router /users/{id}/notes/{note_id} [put]
func New(log *slog.Logger, putNote NoteStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.putNote.New"

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

		var req models.PutNoteRequest
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

		idUserStr := chi.URLParam(r, "id")
		if idUserStr == "" {
			log.Info("Id is empty")
			render.JSON(w, r, resp.Error("Id is empty"))
			render.Status(r, http.StatusBadRequest)
			return
		}

		idUser, err := strconv.ParseInt(idUserStr, 10, 64)
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

		idNoteStr := chi.URLParam(r, "note_id")
		if idNoteStr == "" {
			log.Info("Note id is empty")
			render.JSON(w, r, resp.Error("Note id is empty"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		idNote, err := strconv.ParseInt(idNoteStr, 10, 64)
		if err != nil {
			log.Error("Failed to convert id to int64", "error", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Invalid id format: must be integer"))
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

		note, err := putNote.PutNote(idUser, idNote, Title, Content)
		if err != nil {
			if errors.Is(err, storageErr.ErrNoteNotFound) {
				log.Error("Note not found", "error", sl.Err(err))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("Note not found"))
				return
			} else {
				log.Error("Failed to put note", "error", sl.Err(err))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("Failed to put note"))
				return
			}
		}
		log.Info("Success", slog.Int64("idUser", idUser), slog.Int64("idNote", idNote))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, models.NoteResponse{
			Response:  resp.OK("Success"),
			NoteID:    note.ID,
			UserId:    note.UserID,
			Title:     note.Title,
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})

	}

}
