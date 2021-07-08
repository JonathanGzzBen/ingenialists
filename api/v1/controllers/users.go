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

type UpdateUserDTO struct {
	Name              string    `json:"name" binding:"required"`
	Birthdate         time.Time `json:"birthdate" example:"2006-01-02T15:04:05Z"`
	Gender            string    `json:"gender"`
	ProfilePictureURL string    `json:"profilePictureUrl"`
	Description       string    `json:"description"`
	ShortDescription  string    `json:"shortDescription"`
	Role              string    `json:"role" example:"User"`
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

// UpdateUser is the handler for PUT requests to /users/:id
// 	@ID UpdateUser
// 	@Summary Update user
// 	@Description Update matching user with provided data.
// 	@Tags users
// 	@Security AccessToken
// 	@Param id path int true "User ID"
// 	@Param user body UpdateUserDTO true "User"
// 	@Success 200 {object} models.User
// 	@Failure 400 {object} models.APIError
// 	@Router /users/{id} [put]
func (uc *UsersController) UpdateUser(c *gin.Context) {
	at := c.GetHeader(accessTokenName)
	au, err := getAuthenticatedUser(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "not authenticated: " + err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "id is not a valid"})
		return
	}
	if au.ID != uint(id) {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "id does not match authenticated user"})
		return
	}

	var uuDTO *UpdateUserDTO
	if err := c.ShouldBindJSON(uuDTO); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "invalid user"})
		return
	}
	uu, err := parseUpdateUserDTO(uuDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "invalid user"})
		return
	}
	var u models.User
	uc.db.First(&u, au.ID)
	u.Name = uu.Name
	u.Birthdate = uu.Birthdate
	u.Gender = uu.Gender
	u.ProfilePictureURL = uu.ProfilePictureURL
	u.Description = uu.Description
	u.ShortDescription = uu.ShortDescription
	u.Role = uu.Role
	res := uc.db.Save(&u)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "invalid user"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func parseUpdateUserDTO(uc *UpdateUserDTO) (*models.User, error) {
	ucJSON, err := json.Marshal(uc)
	if err != nil {
		return nil, err
	}
	var u *models.User
	err = json.Unmarshal(ucJSON, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
