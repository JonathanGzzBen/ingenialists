package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetAllCategories(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"))
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
	ts := httptest.NewServer(server.Router)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/v1/categories", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]

	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}

	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}
