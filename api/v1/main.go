package main

import (
	"log"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/JonathanGzzBen/ingenialists/api/v1/controllers"
	_ "github.com/JonathanGzzBen/ingenialists/api/v1/docs"
	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
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

	// Initialize database
	db, err := models.DB()
	if err != nil {
		log.Panic("could not open database file")
	}
	db.AutoMigrate(&models.User{})

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		uc := controllers.NewUsersController(db)
		ur := v1.Group("/users")
		{
			ur.GET("/", uc.GetAllUsers)
			ur.GET("/:id", uc.GetUser)
			ur.POST("/", uc.CreateUser)
		}
		ac := controllers.NewAuthController(db)
		ar := v1.Group("/auth")
		{
			ar.GET("/", ac.GetCurrentUser)
			ar.GET("/google-login", ac.LoginGoogle)
			ar.GET("/google-callback", ac.GoogleCallback)
		}
	}

	swaggerUrl := ginSwagger.URL("http://localhost:8080/v1/swagger/doc.json")
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))

	r.Run()
}
