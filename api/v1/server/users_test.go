package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
)

var mockUsers = []models.User{
	{ID: 123, Name: "First User", Role: models.RoleReader},
	{ID: 456, Name: "Second User", Role: models.RoleReader},
	{ID: 789, Name: "Third User", Role: models.RoleReader},
}

func TestGetAllUsers(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockUsers)

	res, err := http.Get(fmt.Sprintf("%s/v1/users", ts.URL))
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

	var resUsers []models.User
	err = json.NewDecoder(res.Body).Decode(&resUsers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(mockUsers) != len(resUsers) {
		t.Fatalf("Expected %v, got %v", len(mockUsers), len(resUsers))
	}
}

func TestGetUser(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockUsers)
	mockUser := mockUsers[1]

	res, err := http.Get(fmt.Sprintf("%s/v1/users/%d", ts.URL, mockUser.ID))
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

	var resUser models.User
	err = json.NewDecoder(res.Body).Decode(&resUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if mockUser != resUser {
		t.Fatalf("Expected %v, got %v", mockUser, resUser)
	}
}
