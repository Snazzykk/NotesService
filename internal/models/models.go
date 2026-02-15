package models

import (
	resp "NotesService/internal/api/response"
	"time"
)

type Note struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        int64
	Username  string
	CreatedAt time.Time
}
type UserRequest struct {
	Username string `json:"user_name" validate:"required,min=3" example:"john_doe"`
}

type UserResponse struct {
	resp.Response
	ID        int64     `json:"id" example:"1"`
	Username  string    `json:"user_name" example:"john_doe"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-01T12:00:00Z"`
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
type NoteResponse struct {
	resp.Response
	NoteID    int64     `json:"noteID" example:"1"`
	UserId    int64     `json:"userId" example:"1"`
	Title     string    `json:"title" example:"note title"`
	Content   string    `json:"content" example:"note content"`
	CreatedAt time.Time `json:"createdAt" example:"2026-02-15T18:01:29.342814+02:00"`
	UpdatedAt time.Time `json:"updatedAt" example:"2026-02-15T18:01:29.342814+02:00"`
}

type PutNoteRequest struct {
	TitleNote   string `json:"title" validate:"required" example:"My new title"`
	ContentNote string `json:"content" validate:"required" example:"Updated note content"`
}

type SaveNoteRequest struct {
	TitleNote   string `json:"title" validate:"required" example:"My new title"`
	ContentNote string `json:"content" validate:"required" example:"Updated note content"`
}

type DeleteResponse struct {
	resp.Response
}
