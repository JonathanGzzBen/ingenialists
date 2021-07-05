package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
)

type CreateUserDTO struct {
	GoogleSub         string    `form:"googleSub" json:"googleSub"`
	Name              string    `form:"name" json:"name" binding:"required"`
	Birthdate         time.Time `form:"birthdate" json:"birthdate" binding:"required" example:"2006-01-02T15:04:05Z"`
	Gender            string    `form:"gender" json:"gender" binding:"required"`
	ProfilePictureURL string    `form:"profilePictureUrl" json:"profilePictureUrl"`
	Description       string    `form:"description" json:"description"`
	ShortDescription  string    `form:"shortDescription" json:"shortDescription"`
	Role              string    `form:"role" json:"role" binding:"required" example:"User"`
}

var mockUsers = []models.User{
	{
		GoogleSub:        "21231231234",
		Name:             "First Mock User",
		Gender:           "male",
		Description:      "I am the first user of this website",
		ShortDescription: "First user",
		Role:             "user",
	},
	{
		GoogleSub:        "59343786781",
		Name:             "Second Mock User",
		Gender:           "male",
		Description:      "I am the second user of this website",
		ShortDescription: "Second user",
		Role:             "user",
	},
}

var db, _ = models.DB()

// GetAllUsers is the handler for GET requests to /users
// 	@ID GetAllUsers
// 	@Summary Get all users
// 	@Description Get all registered users.
// 	@Tags users
// 	@Success 200 {array} models.User
// 	@Failure 500 {object} models.APIError
// 	@Router /users [get]
func GetAllUsers(c *gin.Context) {
	var users models.User
	result := db.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not connect to database"})
		return
	}
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

// CreateUser is the handler for POST requests to /users
// 	@ID CreateUser
// 	@Summary Create user
// 	@Description Creates a new user.
// 	@Tags users
// 	@Param user body CreateUserDTO true "User"
// 	@Success 200 {object} models.User
// 	@Failure 400 {object} models.APIError
// 	@Router /users [post]
func CreateUser(c *gin.Context) {
	var cu CreateUserDTO

	if err := c.ShouldBindJSON(&cu); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid create user request: " + err.Error()})
		return
	}

	u, err := parseCreateUserDTO(cu)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "could not parse user: " + err.Error()})
		return
	}

	db.Create(&u)
	c.JSON(http.StatusOK, u)
}

func parseCreateUserDTO(cu CreateUserDTO) (*models.User, error) {
	cuJSON, err := json.Marshal(cu)
	if err != nil {
		return nil, err
	}
	var u models.User
	err = json.Unmarshal(cuJSON, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
