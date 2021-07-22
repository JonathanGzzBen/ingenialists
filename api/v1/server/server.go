package server

import (
	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type Server struct {
	db           *gorm.DB
	googleConfig oauth2.Config
	Router       *gin.Engine
}

type ServerConfig struct {
	DB                 *gorm.DB
	GoogleClientID     string
	GoogleClientSecret string
	Hostname           string
}

func NewServer(sc ServerConfig) *Server {
	server := &Server{
		db: sc.DB,
	}
	server.db.AutoMigrate(&models.Article{})
	server.db.AutoMigrate(&models.Category{})
	server.db.AutoMigrate(&models.User{})
	server.SetupGoogleOAuth2(GoogleClientConfig{
		ClientID:     sc.GoogleClientID,
		ClientSecret: sc.GoogleClientSecret,
	})

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		ur := v1.Group("/users")
		{
			ur.GET("/", server.GetAllUsers)
			ur.GET("/:id", server.GetUser)
			ur.PUT("/:id", server.UpdateUser)
		}
		ar := v1.Group("/auth")
		{
			ar.GET("/", server.GetCurrentUser)
			ar.GET("/google-login", server.LoginGoogle)
			ar.GET("/google-callback", server.GoogleCallback)
		}
		cr := v1.Group("/categories")
		{
			cr.GET("/", server.GetAllCategories)
			cr.GET("/:id", server.GetCategory)
			cr.POST("/", server.CreateCategory)
			cr.PUT("/:id", server.UpdateCategory)
			cr.DELETE("/:id", server.DeleteCategory)
		}
		arr := v1.Group("/articles")
		{
			arr.GET("/", server.GetAllArticles)
			arr.GET("/:id", server.GetArticle)
			arr.POST("/", server.CreateArticle)
			arr.PUT("/:id", server.UpdateArticle)
			arr.DELETE("/:id", server.DeleteArticle)
		}
	}

	swaggerUrl := ginSwagger.URL(sc.Hostname + "/v1/swagger/doc.json")
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))
	server.Router = router
	return server
}

func (s *Server) Run(port ...string) {
	s.Router.Run(port[0])
}