package controllers

import (
	"net/http"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
)

// GetAllUsers is the handler for GET requests to /users
// 	@ID GetAllUsers
// 	@Summary Get all users
// 	@Description Get all registered users.
// 	@Tags users
// 	@Success 200 {array} models.User
// 	@Failure 500 {object} models.APIError
// 	@Router /users [get]
func GetAllUsers(c *gin.Context) {
	users := []models.User{
		{
			ID:               1,
			Name:             "First Mock User",
			Gender:           "male",
			Description:      "I am the first user of this website",
			ShortDescription: "First user",
			Role:             "user",
		},
	}
	c.JSON(http.StatusOK, users)
}
