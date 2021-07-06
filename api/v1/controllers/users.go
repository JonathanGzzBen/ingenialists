package controllers

import (
	"net/http"
	"strconv"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UsersController struct{ db *gorm.DB }

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
