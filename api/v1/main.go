package main

import (
	_ "github.com/JonathanGzzBen/ingenialists/api/v1/docs"
	"github.com/JonathanGzzBen/ingenialists/api/v1/pkg"
)

// @title Ingenialists API V1
// @version 0.1.0
// @description This is Ingenialist's API
//
// @contact.name JonathanGzzBen
// @contact.url http://www.github.com/JonathanGzzBen
// @contact.email jonathangzzben@gmail.com
// @license.name MIT License
// @license.url https://mit-license.org/
//
// @host localhost:8080
// @BasePath /v1
//
// @securityDefinitions.apikey AccessToken
// @in header
// @name AccessToken
//
// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl /v1/auth/google-callback
// @authorizationUrl /v1/auth/google-login
// @scope.openid Allow identifying account
// @scope.profile Grant access to profile
// @scope.email Grant access to email
func main() {
	api := pkg.IngenialistsAPIV1{}
	api.Run()
}
