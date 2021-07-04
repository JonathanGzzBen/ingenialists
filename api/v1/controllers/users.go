package controllers

import (
	"net/http"
	"strconv"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
)

var mockUsers = []models.User{
	{
		ID:               1,
		GoogleSub:        "21231231234",
		Name:             "First Mock User",
		Gender:           "male",
		Description:      "I am the first user of this website",
		ShortDescription: "First user",
		Role:             "user",
	},
	{
		ID:               2,
		GoogleSub:        "59343786781",
		Name:             "Second Mock User",
		Gender:           "male",
		Description:      "I am the second user of this website",
		ShortDescription: "Second user",
		Role:             "user",
	},
}

// GetAllUsers is the handler for GET requests to /users
// 	@ID GetAllUsers
// 	@Summary Get all users
// 	@Description Get all registered users.
// 	@Tags users
// 	@Success 200 {array} models.User
// 	@Failure 500 {object} models.APIError
// 	@Router /users [get]
func GetAllUsers(c *gin.Context) {
	users := mockUsers
	c.JSON(http.StatusOK, users)
}

// GetUser is the handler for GET requests to /users/:id
// 	@ID GetUser
// 	@Summary Get user
// 	@Description Get user with matching ID.
// 	@Tags users
// 	@Param id path int true "User ID"
// 	@Success 200 {object} models.User
// 	@Failure 404 {object} models.APIError
// 	@Router /users/{id} [get]
func GetUser(c *gin.Context) {
	for _, u := range mockUsers {
		if strconv.Itoa(int(u.ID)) == c.Param("id") {
			c.JSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "user not found"})
}
