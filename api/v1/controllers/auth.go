package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"gorm.io/gorm"
)

var (
	state = "ingenialists"

	googleClientID     string
	googleClientSecret string
	googleUserInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"
	googleCallbackURL  = "http://127.0.0.1:8080/v1/auth/google-callback"
	googleConfig       oauth2.Config

	accessTokenName = "AccessToken"
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
	googleClientID = os.Getenv("ING_GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("ING_GOOGLE_CLIENT_SECRET")
	googleConfig = oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Endpoint:     endpoints.Google,
		RedirectURL:  googleCallbackURL,
		Scopes:       []string{"openid", "profile", "email"},
	}
	if len(googleClientID) == 0 || len(googleClientSecret) == 0 {
		panic("Environment variables ING_GOOGLE_CLIENT_ID or ING_CLIENT_SECRET missing")
	}
	ac := AuthController{
		db: db,
	}
	return ac
}

// CurrentUser is the handler for GET requests to /auth
// 	@ID GetCurrentUser
// 	@Tags auth
// 	@Success 200 {object} string
// 	@Failure 403 {object} models.APIError
// 	@Security AccessToken
// 	@Router /auth [get]
func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	at := c.GetHeader(accessTokenName)
	u, err := ac.userByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "invalid access token"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// LoginGoogle is the handler for GET requests to /auth/google-login
// it's the entryway for Google OAuth2 flow.
func (ac *AuthController) LoginGoogle(c *gin.Context) {
	log.Println("Client ID: ", googleClientSecret)
	log.Println("Client Secret: ", googleClientSecret)
	url := googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback is the handler for GET requests to /auth/google-callback
// it's part of Google OAuth2 flow.
//
// Returns user's token.
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

	uinfo, err := userInfoByAccessToken(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "failed to get user info: " + err.Error()})
		return
	}

	var u models.User
	res := ac.db.Where("google_sub = ? ", uinfo.Sub).First(&u)
	// If there is no user with that sub, create one
	if res.Error != nil {
		u = models.User{
			GoogleSub:         uinfo.Sub,
			ProfilePictureURL: uinfo.Picture,
			Name:              uinfo.Name,
		}
		ac.db.Save(&u)
	}

	c.JSON(http.StatusOK, token)
}

func (ac *AuthController) userByAccessToken(at string) (*models.User, error) {
	ui, err := userInfoByAccessToken(at)
	if err != nil {
		return nil, err
	}
	var u *models.User
	res := ac.db.Where("google_sub = ? ", ui.Sub).First(&u)
	if res.Error != nil {
		return nil, err
	}
	return u, nil
}

// userInfoByAccessToken returns userInfo
func userInfoByAccessToken(at string) (*googleUserInfoResponse, error) {
	response, err := http.Get(googleUserInfoURL + "?access_token=" + at)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var uinfo *googleUserInfoResponse
	json.NewDecoder(response.Body).Decode(&uinfo)
	return uinfo, nil
}

func getAuthenticatedUser(accessToken string) (*models.User, error) {
	// Get AuthenticatedUser
	var baseURL string
	switch env := os.Getenv("ING_ENVIRONMENT"); env {
	case "development":
		baseURL = os.Getenv("ING_HOSTNAME") + os.Getenv("ING_PORT")
	default:
		baseURL = os.Getenv("ING_HOSTNAME")
	}
	req, err := http.NewRequest("GET", baseURL+"/v1/auth", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(accessTokenName, accessToken)
	log.Println("Access token: ", accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var au models.User
	err = json.NewDecoder(res.Body).Decode(&au)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return &au, nil
}
