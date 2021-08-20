package server

import (
	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	repositories "github.com/JonathanGzzBen/ingenialists/api/v1/repository"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/gorm"
)

type Server struct {
	db             *gorm.DB
	googleClient   IGoogleClient
	googleConfig   IOauthConfig
	development    bool
	Router         *gin.Engine
	CategoriesRepo repositories.CategoriesRepository
}

type ServerConfig struct {
	DB             *gorm.DB
	GoogleConfig   IOauthConfig
	Hostname       string
	Development    bool
	CategoriesRepo repositories.CategoriesRepository
}

func NewServer(sc ServerConfig) *Server {
	server := &Server{
		db:             sc.DB,
		googleConfig:   sc.GoogleConfig,
		development:    sc.Development,
		CategoriesRepo: sc.CategoriesRepo,
	}
	if sc.Development {
		server.googleClient = &GoogleClientMock{}
	} else {
		server.googleClient = &GoogleClient{}
	}
	server.db.AutoMigrate(&models.Article{})
	server.db.AutoMigrate(&models.Category{})
	server.db.AutoMigrate(&models.User{})

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
			if sc.Development {
				ar.GET("/dev-authorize", server.devOAuthAuthorize)
			}
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
