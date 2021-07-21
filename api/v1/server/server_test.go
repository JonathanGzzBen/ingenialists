package server_test

import (
	"os"

	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestEnvironment struct {
	Server *server.Server
	DB     *gorm.DB
}

func NewTestEnvironment() *TestEnvironment {
	os.Remove("test.db")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	}
	server := server.NewServer(
		server.ServerConfig{
			DB:                 db,
			GoogleClientID:     os.Getenv("ING_GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("ING_GOOGLE_CLIENT_SECRET"),
			Hostname:           "http://localhost:8080",
		},
	)
	ts := &TestEnvironment{
		Server: server,
		DB:     db,
	}
	return ts
}
