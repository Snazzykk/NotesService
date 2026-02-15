package deleteNote

import (
	resp "NotesService/internal/api/response"
	"NotesService/internal/auth"
	"NotesService/internal/models"
	"NotesService/internal/storage"
	"NotesService/internal/storage/storageErr"
	sl "NotesService/pkg/logger/logSlog"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type NoteStorage interface {
	storage.NoteStorage
}

// DeleteNote godoc
// @Summary Delete a note
// @Description Delete a specific note by ID. Requires JWT authentication.
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Param note_id path int true "Note ID" minimum(1)
// @Success 200 {object} models.DeleteResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Security ApiKeyAuth
// @Router /users/{id}/notes/{note_id} [delete]
func New(log *slog.Logger, deleteNote NoteStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.deleteNote.New"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		authorizedUserID, ok := auth.GetUserID(r)
		if !ok {
			log.Error("user_id not found in context")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Unauthorized"))
			return
		}
		idUserStr := chi.URLParam(r, "id")
		if idUserStr == "" {
			log.Info("User id is empty")
			render.JSON(w, r, resp.Error("User id is empty"))
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
		if authorizedUserID != idUser {
			log.Warn("Unauthorized access attempt", slog.Int64("authorized_user_id", authorizedUserID), slog.Int64("requested_user_id", idUser))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Not found"))
			return
		}
		err = deleteNote.DeleteNote(idUser, idNote)
		if err != nil {
			if errors.Is(err, storageErr.ErrNoteNotFound) {
				log.Error("Note not found", "error", sl.Err(err))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("Note not found"))
				return
			} else {
				log.Error("Failed to get Note", "error", sl.Err(err))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("Failed to get Note"))
				return
			}
		}
		log.Info("Success", slog.Int64("idUser", idUser), slog.Int64("idNote", idNote))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, models.DeleteResponse{Response: resp.OK("Success Delete")})
	}
}
