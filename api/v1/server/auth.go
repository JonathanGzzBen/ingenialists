package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var (
	state = "ingenialists"

	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
	googleCallbackURL = "http://127.0.0.1:8080/v1/auth/google-callback"
	gc                oauth2.Config

	accessTokenName = "AccessToken"
)

type googleUserInfoResponse struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

type GoogleClientConfig struct {
	ClientID     string
	ClientSecret string
}

// SetupGoogleOAuth2 initializes GoogleOAuth2 service
// and is necessary to use auth service
func (s *Server) SetupGoogleOAuth2(gcf GoogleClientConfig) {
	gc = oauth2.Config{
		ClientID:     gcf.ClientID,
		ClientSecret: gcf.ClientSecret,
		Endpoint:     endpoints.Google,
		RedirectURL:  googleCallbackURL,
		Scopes:       []string{"openid", "profile", "email"},
	}
	s.googleConfig = gc
}

// CurrentUser is the handler for GET requests to /auth
// 	@ID GetCurrentUser
// 	@Tags auth
// 	@Success 200 {object} string
// 	@Failure 403 {object} models.APIError
// 	@Security AccessToken
// 	@Router /auth [get]
func (s *Server) GetCurrentUser(c *gin.Context) {
	at := c.GetHeader(accessTokenName)
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
	url := gc.AuthCodeURL(state, oauth2.AccessTypeOffline)
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
	token, err := gc.Exchange(ctx, authCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "failed to exchange token: " + err.Error()})
		return
	}

	uinfo, err := s.userInfoByAccessToken(token.AccessToken)
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
	ui, err := s.userInfoByAccessToken(at)
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

// userInfoByAccessToken returns userInfo
func (s *Server) userInfoByAccessToken(at string) (*googleUserInfoResponse, error) {
	response, err := http.Get(googleUserInfoURL + "?access_token=" + at)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("invalid access token")
	}
	defer response.Body.Close()
	var uinfo *googleUserInfoResponse
	json.NewDecoder(response.Body).Decode(&uinfo)
	return uinfo, nil
}
