package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
)

// This data should not be modified, its purpose
// is to be used to initialize database.
var mockCategories = []models.Category{
	{ID: 1231, Name: "First Category", ImageURL: "https://i.imgur.com/oCsJWt7.jpeg"},
	{ID: 2131, Name: "Second Name", ImageURL: "https://i.imgur.com/oCsJWt7.jpeg"},
	{ID: 56232, Name: "Third Name", ImageURL: "https://i.imgur.com/oCsJWt7.jpeg"},
}

func TestGetAllCategories(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockCategories)

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

	var resCategories []models.Category
	err = json.NewDecoder(res.Body).Decode(&resCategories)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(mockCategories) != len(resCategories) {
		t.Fatalf("Expected %v, got %v", len(mockCategories), len(resCategories))
	}
}

func TestGetCategory(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockCategories)
	// Take second category to make sure it's finding
	// it by the ID and not the first item in DB
	mockCategory := mockCategories[1]

	res, err := http.Get(fmt.Sprintf("%s/v1/categories/"+strconv.Itoa(int(mockCategory.ID)), ts.URL))
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
	var resCategory models.Category
	err = json.NewDecoder(res.Body).Decode(&resCategory)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(resCategory, mockCategory) {
		t.Fatalf("Expected %v, got %v", mockCategory, resCategory)
	}

}

func TestCreateCategoryAsRegularUserReturnForbidden(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockCategories)
	mockCategory := mockCategories[1]

	mcJSONBytes, err := json.Marshal(mockCategory)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/categories", ts.URL), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "AccessToken")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %v, got %v", http.StatusForbidden, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}
