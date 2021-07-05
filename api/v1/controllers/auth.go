package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"gorm.io/gorm"
)

var (
	state = "ingenialists"

	googleClientID     = "508942453082-n2hslfnvv37kvebfcqp0ii33idc7tv4s.apps.googleusercontent.com"
	googleClientSecret = "MB5vDo99iMvLv3gRU4xLfi1C"
	googleUserInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"

	googleCallbackURL = "http://127.0.0.1:8080/v1/auth/google-callback"

	googleConfig = oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Endpoint:     endpoints.Google,
		RedirectURL:  googleCallbackURL,
		Scopes:       []string{"openid", "profile", "email"},
	}
)

type AuthController struct {
	db *gorm.DB
}

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

func NewAuthController(db *gorm.DB) AuthController {
	ac := AuthController{
		db: db,
	}
	return ac
}

// LoginGoogle is the handler for GET requests to /auth/google-login
// 	@ID LoginGoogle
// 	@Summary Login with Google
// 	@Description Logins with Google Oauth2
// 	@Tags auth
// 	@Success 302 {object} string
// 	@Router /auth/google-login [get]
func (ac *AuthController) LoginGoogle(c *gin.Context) {
	url := googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// LoginGoogle is the handler for GET requests to /auth/google-login
func (ac *AuthController) GoogleCallback(c *gin.Context) {
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
	response, err := http.Get(googleUserInfoURL + "?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "failed to get user info: " + err.Error()})
		return
	}
	defer response.Body.Close()
	var uinfo googleUserInfoResponse
	json.NewDecoder(response.Body).Decode(&uinfo)

	var u models.User
	u.GoogleAccessToken = token.AccessToken
	res := ac.db.Where("google_sub = ? ", uinfo.Sub).First(&u)
	// If there is no user with that sub, create one
	if res.Error != nil {
		u = models.User{
			GoogleSub:          uinfo.Sub,
			GoogleRefreshToken: token.RefreshToken,
			GoogleAccessToken:  token.AccessToken,
			Token:              uuid.New(),
			ProfilePictureURL:  uinfo.Picture,
			Name:               uinfo.Name,
		}
		ac.db.Save(&u)
		c.JSON(http.StatusOK, u)
		return
	}
	// If user found with that sub, return it
	c.JSON(http.StatusOK, u)
}
