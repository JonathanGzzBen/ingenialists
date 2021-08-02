package server

import (
	"net/http"
	"strconv"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type CreateArticleDTO struct {
	CategoryID uint   `json:"categoryId"`
	Body       string `json:"body"`
	Title      string `json:"title"`
	ImageURL   string `json:"imageUrl"`
	Tags       string `json:"tags"`
}

type UpdateArticleDTO struct {
	CategoryID uint   `json:"categoryId"`
	Body       string `json:"body"`
	Title      string `json:"title"`
	ImageURL   string `json:"imageUrl"`
	Tags       string `json:"tags"`
}

// GetAllArticles is the handler for GET requests to /articles
// 	@ID GetAllArticles
// 	@Summary Get all articles
// 	@Description Get all registered articles.
// 	@Tags articles
// 	@Success 200 {array} models.Article
// 	@Failure 500 {object} models.APIError
// 	@Router /articles [get]
func (s *Server) GetAllArticles(c *gin.Context) {
	var a []models.Article
	r := s.db.Preload(clause.Associations).Find(&a)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not get articles"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// GetArticle is the handler for GET requests to /article/:id
// 	@ID GetArticle
// 	@Summary Get article
// 	@Description Get article with matching ID.
// 	@Tags articles
// 	@Param id path int true "Article ID"
// 	@Success 200 {object} models.Article
// 	@Failure 404 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /articles/{id} [get]
func (s *Server) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid id: " + err.Error()})
		return
	}
	var article models.Article
	res := s.db.Find(&article, id)
	if res.Error == nil && res.RowsAffected != 1 {
		c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "category not found"})
		return
	}
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not find category"})
		return
	}
	c.JSON(http.StatusOK, article)
}

// CreateArticles is the handler for POST requests to /articles
// 	@ID CreateArticle
// 	@Summary Create article
// 	@Description Register a new article.
// 	@Tags articles
// 	@Param article body CreateArticleDTO true "Article"
// 	@Security AccessToken
// 	@Success 200 {object} models.Category
// 	@Failure 400 {object} models.APIError
// 	@Failure 403 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /articles [post]
func (s *Server) CreateArticle(c *gin.Context) {
	at := c.GetHeader(AccessTokenName)
	au, err := s.userByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you must be authenticated to create an article"})
		return
	}
	if !(au.Role == models.RoleWriter || au.Role == models.RoleAdministrator) {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "only Writers and Administrators can create articles"})
		return
	}
	var ca CreateArticleDTO
	if err := c.ShouldBindJSON(&ca); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusInternalServerError, Message: "invalid article"})
		return
	}

	// Verify that a category with matching ID exists
	var category models.Category
	res := s.db.Find(&category, ca.CategoryID)
	if res.Error != nil || res.RowsAffected != 1 {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "category with provided id could not be retrieved"})
		return
	}

	article := models.Article{
		UserID:     au.ID,
		CategoryID: ca.CategoryID,
		Body:       ca.Body,
		Title:      ca.Title,
		ImageURL:   ca.ImageURL,
		Tags:       ca.Tags,
	}
	res = s.db.Create(&article)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not create article: " + res.Error.Error()})
		return
	}
	res = s.db.Preload(clause.Associations).Find(&article, article.ID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not retrieve created article: " + res.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}

// UpdateArticle is the handler for PUT requests to /articles
// 	@ID UpdateArticle
// 	@Summary Update article
// 	@Description Updates a registered article.
// 	@Tags articles
// 	@Param id path int true "Article ID"
// 	@Param article body UpdateArticleDTO true "Article"
// 	@Security AccessToken
// 	@Success 200 {object} models.Article
// 	@Failure 400 {object} models.APIError
// 	@Failure 403 {object} models.APIError
// 	@Failure 404 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /articles/{id} [put]
func (s *Server) UpdateArticle(c *gin.Context) {
	at := c.GetHeader(AccessTokenName)
	au, err := s.userByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you must be authenticated to update an article"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid id: " + err.Error()})
		return
	}
	var article models.Article
	res := s.db.Find(&article, id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	if res.RowsAffected != 1 {
		c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "article with provided id not found"})
		return
	}

	if article.UserID != au.ID {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you can only modify articles created by you"})
		return
	}

	var ua UpdateArticleDTO
	if err := c.ShouldBindJSON(&ua); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid article: " + err.Error()})
		return
	}

	article.CategoryID = ua.CategoryID
	article.Body = ua.Body
	article.Title = ua.Title
	article.ImageURL = ua.ImageURL
	article.Tags = ua.Tags

	res = s.db.Save(&article)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusBadRequest, Message: "could not save updated article: " + err.Error()})
		return
	}
	res = s.db.Preload(clause.Associations).Find(&article, article.ID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusBadRequest, Message: "could not retrieve updated article: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}

// DeleteArticle is the handler for DELETE requests to /articles/:id
// 	@ID DeleteArticle
// 	@Summary Delete article
// 	@Description Delete article with matching ID.
// 	@Tags articles
// 	@Param id path int true "Article ID"
// 	@Security AccessToken
// 	@Success 204 {object} string
// 	@Failure 403 {object} models.APIError
// 	@Failure 404 {object} models.APIError
// 	@Failure 500 {object} models.APIError
// 	@Router /articles/{id} [delete]
func (s *Server) DeleteArticle(c *gin.Context) {
	at := c.GetHeader(AccessTokenName)
	au, err := s.userByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you must be authenticated to delete an article"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid id: " + err.Error()})
		return
	}

	var article models.Article
	res := s.db.Find(&article, id)
	if res.Error != nil || res.RowsAffected != 1 {
		c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "article not found"})
		return
	}

	// If article doens't belong to authenticated user
	// and authenticated user is not administrator
	if !(article.UserID == au.ID || au.Role == models.RoleAdministrator) {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "you are not authenticated as administrator or this article doesn't belong to you"})
		return
	}

	res = s.db.Delete(&article)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not delete article: " + err.Error()})
		return
	}
	c.String(http.StatusNoContent, "deleted")
}
