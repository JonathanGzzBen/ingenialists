package models

import "time"

type User struct {
	ID                uint      `json:"id"`
	GoogleSub         string    `json:"googleSub"`
	Name              string    `json:"name"`
	Birthdate         time.Time `json:"birthdate" example:"2006-01-02T15:04:05Z"`
	Gender            string    `json:"gender"`
	ProfilePictureURL string    `json:"profilePictureUrl"`
	Description       string    `json:"description"`
	ShortDescription  string    `json:"shortDescription"`
	Role              string    `json:"role" example:"User"`
}
