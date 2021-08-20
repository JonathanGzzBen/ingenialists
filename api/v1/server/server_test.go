package server_test

import (
	"os"

	repositories "github.com/JonathanGzzBen/ingenialists/api/v1/repository"
	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestEnvironment struct {
	Server *server.Server
	DB     *gorm.DB
}

func (e *TestEnvironment) Close() {
	os.Remove("test.db")
}

func NewTestEnvironment() *TestEnvironment {
	os.Remove("test.db")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	}
	server := server.NewServer(
		server.ServerConfig{
			DB:             db,
			GoogleConfig:   &OAuth2ConfigMock{},
			Hostname:       "http://localhost:8080",
			Development:    true,
			CategoriesRepo: repositories.NewCategoriesGormRepository(db),
		},
	)
	ts := &TestEnvironment{
		Server: server,
		DB:     db,
	}
	return ts
}
