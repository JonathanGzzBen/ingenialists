package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
)

var mockArticles = []models.Article{
	{
		ID:         123,
		UserID:     mockUsers[0].ID,
		CategoryID: mockCategories[0].ID,
		Body:       "First article body",
		Title:      "First article title",
	},
	{
		ID:         456,
		UserID:     mockUsers[1].ID,
		CategoryID: mockCategories[1].ID,
		Body:       "Second article body",
		Title:      "Second article title",
	},
	{
		ID:         789,
		UserID:     mockUsers[2].ID,
		CategoryID: mockCategories[2].ID,
		Body:       "Third article body",
		Title:      "THird article title",
	},
}

func TestGetAllArticles(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockArticles)

	res, err := http.Get(fmt.Sprintf("%s/v1/articles", ts.URL))
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

	var resArticles []models.Article
	err = json.NewDecoder(res.Body).Decode(&resArticles)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(mockArticles) != len(resArticles) {
		t.Fatalf("Expected %v, got %v", len(mockArticles), len(resArticles))
	}
}

func TestGetArticle(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockArticles)

	mockArticle := mockArticles[1]

	res, err := http.Get(fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID))
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

	var resArticle models.Article
	err = json.NewDecoder(res.Body).Decode(&resArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(resArticle, mockArticle) {
		t.Fatalf("Expected %v, got %v", mockArticle, resArticle)
	}
}
