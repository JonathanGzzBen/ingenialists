package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	state             = "ingenialists"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
	googleConfig      IGoogleConfig
	AccessTokenName   = "AccessToken"
)

// CurrentUser is the handler for GET requests to /auth
// 	@ID GetCurrentUser
// 	@Tags auth
// 	@Success 200 {object} string
// 	@Failure 403 {object} models.APIError
// 	@Security AccessToken
// 	@Router /auth [get]
func (s *Server) GetCurrentUser(c *gin.Context) {
	at := c.GetHeader(AccessTokenName)
	u, err := s.userByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "invalid access token"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// LoginGoogle is the handler for GET requests to /auth/google-login
// it's the entryway for Google OAuth2 flow.
func (s *Server) LoginGoogle(c *gin.Context) {
	url := googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback is the handler for GET requests to /auth/google-callback
// it's part of Google OAuth2 flow.
//
// Returns user's token.
func (s *Server) GoogleCallback(c *gin.Context) {
	if c.Request.URL.Query().Get("state") != state {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "state did not match"})
		return
	}

	authCode := c.Request.URL.Query().Get("code")
	ctx := context.Background()
	token, err := googleConfig.Exchange(ctx, authCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "failed to exchange token: " + err.Error()})
		return
	}

	uinfo, err := s.googleClient.userInfoByAccessToken(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "failed to get user info: " + err.Error()})
		return
	}

	var u models.User
	res := s.db.Where("google_sub = ? ", uinfo.Sub).First(&u)
	// If there is no user with that sub, create one
	if res.Error != nil {
		u = models.User{
			GoogleSub:         uinfo.Sub,
			ProfilePictureURL: uinfo.Picture,
			Name:              uinfo.Name,
		}
		s.db.Save(&u)
	}

	c.JSON(http.StatusOK, token)
}

func (s *Server) userByAccessToken(at string) (*models.User, error) {
	ui, err := s.googleClient.userInfoByAccessToken(at)
	if s.development {
		return &models.User{ID: 1, GoogleSub: ui.Sub, Name: ui.Name}, nil
	}
	if err != nil {
		return nil, err
	}
	var u *models.User
	res := s.db.Where("google_sub = ? ", ui.Sub).First(&u)
	if res.Error != nil {
		return nil, err
	}
	return u, nil
}

// devOAuthAuthorize handles requests to /auth/authorize
// should only be available during development
func (s *Server) devOAuthAuthorize(c *gin.Context) {
	state := c.Request.URL.Query().Get("state")

	u, err := url.Parse("http:://localhost/callback")
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			models.APIError{
				Code:    http.StatusInternalServerError,
				Message: "could not make url",
			})
	}

	v := url.Values{}
	v.Set("code", "code")
	v.Set("state", state)
	u.RawQuery = v.Encode()

	http.Redirect(c.Writer, c.Request, u.String(), http.StatusTemporaryRedirect)
}
