package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func V1Router() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		ur := v1.Group("/users")
		{
			ur.GET("/", GetAllUsers)
			ur.GET("/:id", GetUser)
			ur.PUT("/:id", UpdateUser)
		}
		ar := v1.Group("/auth")
		{
			ar.GET("/", GetCurrentUser)
			ar.GET("/google-login", LoginGoogle)
			ar.GET("/google-callbcontrollers.", GoogleCallback)
		}
		cr := v1.Group("/categories")
		{
			cr.GET("/", GetAllCategories)
			cr.GET("/:id", GetCategory)
			cr.POST("/", CreateCategory)
			cr.PUT("/:id", UpdateCategory)
			cr.DELETE("/:id", DeleteCategory)
		}
		arr := v1.Group("/articles")
		{
			arr.GET("/", GetAllArticles)
			arr.GET("/:id", GetArticle)
			arr.POST("/", CreateArticle)
			arr.PUT("/:id", UpdateArticle)
			arr.DELETE("/:id", DeleteArticle)
		}
	}

	// hostname is used by multiple controllers
	// to make requests to authentication controller
	hostname := os.Getenv("ING_HOSTNAME")

	if len(hostname) == 0 {
		panic("Environment variable ING_HOSTNAME missing")
	}

	swaggerUrl := ginSwagger.URL(hostname + "/v1/swagger/doc.json")
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))
	return r
}
