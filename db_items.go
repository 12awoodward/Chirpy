package main

import (
	"time"

	"github.com/12awoodward/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func toJSONUser(old database.User) User {
	return User{
		ID: old.ID,
		CreatedAt: old.CreatedAt,
		UpdatedAt: old.UpdatedAt,
		Email: old.Email,
	}
}

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func toJSONChirp(old database.Chirp) Chirp {
	return Chirp{
		ID: old.ID,
		CreatedAt: old.CreatedAt,
		UpdatedAt: old.UpdatedAt,
		Body: old.Body,
		UserID: old.UserID,
	}
}