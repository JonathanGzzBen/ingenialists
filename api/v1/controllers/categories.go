package controllers

import (
	"net/http"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoriesController struct{ db *gorm.DB }

type CreateCategoryDTO struct {
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

// NewCategoriesController returns a new controller for categories
func NewCategoriesController(db *gorm.DB) CategoriesController {
	return CategoriesController{
		db: db,
	}
}

// GetAllCategories is the handler for GET requests to /categories
// 	@ID GetAllCategories
// 	@Summary Get all categories
// 	@Description Get all registered categories.
// 	@Tags categories
// 	@Success 200 {array} models.Category
// 	@Failure 500 {object} models.APIError
// 	@Router /categories [get]
func (cc *CategoriesController) GetAllCategories(c *gin.Context) {
	var categories models.Category
	r := cc.db.Find(&categories)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not get categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// CreateCategory is the handler for POST requests to /categories
// 	@ID CreateCategory
// 	@Summary Create category
// 	@Description Register a new category.
// 	@Tags categories
// 	@Param category body CreateCategoryDTO true "Category"
// 	@Success 200 {object} models.Category
// 	@Failure 400 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /categories [post]
func (cc *CategoriesController) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusInternalServerError, Message: "invalid category"})
		return
	}
	result := cc.db.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not create category:" + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}
