package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uint      `json:"id"`
	GoogleSub          string    `json:"googleSub"`
	GoogleRefreshToken string    `json:"-"`
	GoogleAccessToken  string    `json:"-"`
	Name               string    `json:"name"`
	Birthdate          time.Time `json:"birthdate" example:"2006-01-02T15:04:05Z"`
	Gender             string    `json:"gender"`
	ProfilePictureURL  string    `json:"profilePictureUrl"`
	Description        string    `json:"description"`
	ShortDescription   string    `json:"shortDescription"`
	Role               string    `json:"role" example:"User"`
	Token              uuid.UUID `json:"token"`
}
