package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
)

func TestGetAllCategories(t *testing.T) {
	e := NewTestEnvironment()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	categoriesInDB := []models.Category{
		{Name: "First Category"},
		{Name: "Second Category"},
		{Name: "Third Category"},
	}
	e.DB.Create(&categoriesInDB)

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
	if len(categoriesInDB) != len(resCategories) {
		t.Fatalf("Expected %v, got %v", len(categoriesInDB), len(resCategories))
	}
}
