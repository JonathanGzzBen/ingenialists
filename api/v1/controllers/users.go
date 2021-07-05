package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UsersController struct{ db *gorm.DB }

type CreateUserDTO struct {
	GoogleSub         string    `json:"googleSub"`
	Name              string    `json:"name" binding:"required"`
	Birthdate         time.Time `json:"birthdate" binding:"required" example:"2006-01-02T15:04:05Z"`
	Gender            string    `json:"gender" binding:"required"`
	ProfilePictureURL string    `json:"profilePictureUrl"`
	Description       string    `json:"description"`
	ShortDescription  string    `json:"shortDescription"`
	Role              string    `json:"role" binding:"required" example:"User"`
}

func NewUsersController(db *gorm.DB) UsersController {
	return UsersController{
		db: db,
	}
}

// GetAllUsers is the handler for GET requests to /users
// 	@ID GetAllUsers
// 	@Summary Get all users
// 	@Description Get all registered users.
// 	@Tags users
// 	@Success 200 {array} models.User
// 	@Failure 500 {object} models.APIError
// 	@Router /users [get]
func (uc *UsersController) GetAllUsers(c *gin.Context) {
	var users models.User
	result := uc.db.Find(&users)
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
func (uc *UsersController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid id: " + err.Error()})
		return
	}
	var u models.User
	res := uc.db.Find(&u, id)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: res.Error.Error()})
		return
	}
	if res.RowsAffected != 1 {
		c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "user with provided id not found"})
		return
	}
	c.JSON(http.StatusOK, u)
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
func (uc *UsersController) CreateUser(c *gin.Context) {
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

	uc.db.Create(&u)
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
