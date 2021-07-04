package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/JonathanGzzBen/ingenialists/api/v1/controllers"
	_ "github.com/JonathanGzzBen/ingenialists/api/v1/docs"
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
func main() {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		ur := v1.Group("/users")
		{
			ur.GET("/", controllers.GetAllUsers)
			ur.GET("/:id", controllers.GetUser)
		}
		ar := v1.Group("/auth")
		{
			ar.GET("/google-login", controllers.LoginGoogle)
			ar.GET("/google-callback", controllers.GoogleCallback)
		}
	}

	swaggerUrl := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))

	r.Run()
}
