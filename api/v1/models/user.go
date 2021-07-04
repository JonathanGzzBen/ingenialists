package models

type User struct {
	ID                uint   `json:"id"`
	GoogleSub         string `json:"googleSub"`
	Name              string `json:"name"`
	Birthdate         string `json:"birthdate"`
	Gender            string `json:"gender"`
	ProfilePictureURL string `json:"profilePictureUrl"`
	Description       string `json:"description"`
	ShortDescription  string `json:"shortDescription"`
	Role              string `json:"role"`
}
