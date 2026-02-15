package getAllNotes

import (
	resp "NotesService/internal/api/response"
	"NotesService/internal/auth"
	"NotesService/internal/models"
	"NotesService/internal/storage"
	sl "NotesService/pkg/logger/logSlog"
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

// GetAllNotes godoc
// @Summary Get all notes for a user
// @Description Returns all notes for a specific user. Requires JWT authentication.
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Param limit query int false "Limit number of notes" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Param sort query string false "Sort by field (createdAt)"
// @Success 200 {array} models.NoteResponse "List of notes"
// @Failure 400
// @Failure 401
// @Failure 500
// @Security ApiKeyAuth
// @Router /users/{id}/notes [get]
func New(log *slog.Logger, getAllNotes NoteStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.getAllNotes.New"

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

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		sort := r.URL.Query().Get("sort")

		notes, err := getAllNotes.GetAllNotes(idUser, limit, offset, sort)
		if err != nil {
			log.Error("Failed to get all notes", "error", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Failed to get all notes"))
			return
		}

		log.Info("Success", slog.Int64("id", idUser))

		for _, note := range notes {
			render.Status(r, http.StatusOK)
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

}
