package controllers

import (
	"net/http"
	"strconv"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
)

type CreateCategoryDTO struct {
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

type UpdateCategoryDTO struct {
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

// GetAllCategories is the handler for GET requests to /categories
// 	@ID GetAllCategories
// 	@Summary Get all categories
// 	@Description Get all registered categories.
// 	@Tags categories
// 	@Success 200 {array} models.Category
// 	@Failure 500 {object} models.APIError
// 	@Router /categories [get]
func GetAllCategories(c *gin.Context) {
	var categories models.Category
	r := models.DB.Find(&categories)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not get categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetCategory is the handler for GET requests to /categories/:id
// 	@ID GetCategory
// 	@Summary Get category
// 	@Description Get category with matching ID.
// 	@Tags categories
// 	@Param id path int true "Category ID"
// 	@Success 200 {object} models.Category
// 	@Failure 404 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /categories/{id} [get]
func GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid id: " + err.Error()})
		return
	}
	var category models.Category
	res := models.DB.Find(&category, id)
	if res.Error == nil && res.RowsAffected != 1 {
		c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "category not found"})
		return
	}
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not find category"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// CreateCategory is the handler for POST requests to /categories
// 	@ID CreateCategory
// 	@Summary Create category
// 	@Description Register a new category.
// 	@Tags categories
// 	@Param category body CreateCategoryDTO true "Category"
// 	@Security AccessToken
// 	@Success 200 {object} models.Category
// 	@Failure 400 {object} models.APIError
// 	@Failure 403 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /categories [post]
func CreateCategory(c *gin.Context) {
	at := c.GetHeader(accessTokenName)
	u, err := getAuthenticatedUser(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you must be authenticated to create a category"})
		return
	}
	if u.Role != models.RoleAdministrator {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "only users with role Administrator can create categories"})
		return
	}
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusInternalServerError, Message: "invalid category"})
		return
	}
	result := models.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not create category:" + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// UpdateCategory is the handler for PUT requests to /categories
// 	@ID UpdateCategory
// 	@Summary Update category
// 	@Description Updates a registered category.
// 	@Tags categories
// 	@Param id path int true "Category ID"
// 	@Param category body UpdateCategoryDTO true "Category"
// 	@Security AccessToken
// 	@Success 200 {object} models.Category
// 	@Failure 400 {object} models.APIError
// 	@Failure 403 {object} models.APIError
// 	@Failure 404 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	at := c.GetHeader(accessTokenName)
	u, err := getAuthenticatedUser(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you must be authenticated to update a category"})
		return
	}
	if u.Role != models.RoleAdministrator {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "only users with role Administrator can update categories"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid id: " + err.Error()})
		return
	}
	var category models.Category
	res := models.DB.Find(&category, id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	if res.RowsAffected != 1 {
		c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "category with provided id not found"})
		return
	}

	var cu UpdateCategoryDTO
	if err := c.ShouldBindJSON(&cu); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid category: " + err.Error()})
		return
	}

	category.Name = cu.Name
	category.ImageURL = cu.ImageURL
	res = models.DB.Save(&category)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusBadRequest, Message: "could not save updated category: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// DeleteCategory is the handler for DELETE requests to /categories/:id
// 	@ID DeleteCategory
// 	@Summary Delete category
// 	@Description Delete category with matching ID.
// 	@Tags categories
// 	@Param id path int true "Category ID"
// 	@Security AccessToken
// 	@Success 204 {object} string
// 	@Failure 403 {object} models.APIError
// 	@Failure 404 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	at := c.GetHeader(accessTokenName)
	au, err := getAuthenticatedUser(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you must be authenticated to delete a category"})
		return
	}
	// If authenticated user is not Administrator
	if au.Role != models.RoleAdministrator {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you are not authenticated as administrator"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid id: " + err.Error()})
		return
	}

	var category models.Category
	res := models.DB.Find(&category, id)
	if res.Error != nil || res.RowsAffected != 1 {
		c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "category not found"})
		return
	}

	res = models.DB.Delete(&category)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not delete article: " + err.Error()})
		return
	}
	c.String(http.StatusNoContent, "deleted")
}
