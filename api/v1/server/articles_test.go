package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
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

func TestCreateArticleAsUnauthenticatedUserReturnForbidden(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockCategories)
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockArticles)
	mockArticle := models.Article{
		UserID:     mockUsers[1].ID,
		CategoryID: mockCategories[1].ID,
		Title:      "New Article",
	}

	maJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/articles", ts.URL), bytes.NewBuffer(maJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
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

func TestCreateArticleAsReaderReturnForbidden(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockCategories)
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockArticles)
	mockArticle := models.Article{
		UserID:     mockUsers[1].ID,
		CategoryID: mockCategories[1].ID,
		Title:      "New Article",
	}

	maJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/articles", ts.URL), bytes.NewBuffer(maJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Reader")
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

func TestCreateArticleAsWriterReturnOk(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockCategories)
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockArticles)
	mockArticle := models.Article{
		UserID:     mockUsers[1].ID,
		CategoryID: mockCategories[1].ID,
		Title:      "New Article",
	}

	maJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/articles", ts.URL), bytes.NewBuffer(maJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, res.StatusCode)
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
	if resArticle.Title != mockArticle.Title {
		t.Fatalf("Expected %v, got %v", mockArticle.Title, resArticle.Title)
	}
}

func TestCreateArticleAsAdministratorReturnOk(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockCategories)
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockArticles)
	mockArticle := models.Article{
		UserID:     mockUsers[1].ID,
		CategoryID: mockCategories[1].ID,
		Title:      "New Article",
	}

	maJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/articles", ts.URL), bytes.NewBuffer(maJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, res.StatusCode)
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
	if resArticle.Title != mockArticle.Title {
		t.Fatalf("Expected %v, got %v", mockArticle.Title, resArticle.Title)
	}
}

func TestUpdateArticleAsUnauthenticatedUserReturnForbidden(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]
	mockArticle.Title = "Article Updated"

	mcJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
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

func TestUpdateArticleAsReaderReturnForbidden(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]
	mockArticle.Title = "Article Updated"

	mcJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Access Token")
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

func TestUpdateArticleAsWriterThatOwnsArticleReturnOk(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	mockArticles[1].UserID = 1
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]
	mockArticle.Title = "Article Updated"

	mcJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}

func TestUpdateArticleAsWriterThatDoesNotOwnArticleReturnOk(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	mockArticles[1].UserID = 2
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]
	mockArticle.Title = "Article Updated"

	mcJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
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

func TestUpdateArticleAsAdministratorThatOwnsArticleReturnOk(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	mockArticles[1].UserID = 1
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]
	mockArticle.Title = "Article Updated"

	mcJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}

func TestUpdateArticleAsAdministratorThatDoesNotOwnArticleReturnOk(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	mockArticles[1].UserID = 2
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]
	mockArticle.Title = "Article Updated"

	mcJSONBytes, err := json.Marshal(mockArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
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

func TestDeleteArticleAsUnauthenticatedUserReturnForbidden(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
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

	// Verify that mockCategory is still in database
	var aInDB *models.Article
	e.DB.Find(&aInDB, mockArticle.ID)
	if aInDB.Title != mockArticle.Title {
		t.Fatalf("Expected %v, got %v", mockArticle.Title, aInDB.Title)
	}
}

func TestDeleteArticleAsReaderReturnForbidden(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Reader")
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

	// Verify that mockCategory is still in database
	var aInDB *models.Article
	e.DB.Find(&aInDB, mockArticle.ID)
	if aInDB.Title != mockArticle.Title {
		t.Fatalf("Expected %v, got %v", mockArticle.Title, aInDB.Title)
	}
}

func TestDeleteArticleAsWriterThatDoesNotOwnArticleReturnForbidden(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	mockArticles[1].UserID = 2
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
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

	// Verify that mockCategory is still in database
	var aInDB *models.Article
	e.DB.Find(&aInDB, mockArticle.ID)
	if aInDB.Title != mockArticle.Title {
		t.Fatalf("Expected %v, got %v", mockArticle.Title, aInDB.Title)
	}
}

func TestDeleteArticleAsWriterThatOwnsArticleReturnNoContent(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	mockArticles[1].UserID = 1
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %v, got %v", http.StatusNoContent, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "text/plain; charset=utf-8" {
		t.Fatalf("Expected \"text/plain; charset=utf-8\", got %s", val[0])
	}

	// Verify that mockCategory is not in database
	var aInDB *models.Article
	tx := e.DB.Find(&aInDB, mockArticle.ID)
	if tx.RowsAffected != 0 {
		t.Fatalf("Expected %v, got %v", 0, tx.RowsAffected)
	}
}

func TestDeleteArticleAsAdministratorThatDoesNotOwnArticleReturnNoContent(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	mockArticles[1].UserID = 2
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %v, got %v", http.StatusNoContent, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "text/plain; charset=utf-8" {
		t.Fatalf("Expected \"text/plain; charset=utf-8\", got %s", val[0])
	}

	// Verify that mockCategory is not in database
	var aInDB *models.Article
	tx := e.DB.Find(&aInDB, mockArticle.ID)
	if tx.RowsAffected != 0 {
		t.Fatalf("Expected %v, got %v", 0, tx.RowsAffected)
	}
}

func TestDeleteArticleAsAdministratorThatOwnsArticleReturnNoContent(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	mockArticles[1].UserID = 1
	e.DB.Create(&mockUsers)
	e.DB.Create(&mockCategories)
	e.DB.Create(&mockArticles)
	mockArticle := mockArticles[1]

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, mockArticle.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %v, got %v", http.StatusNoContent, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "text/plain; charset=utf-8" {
		t.Fatalf("Expected \"text/plain; charset=utf-8\", got %s", val[0])
	}

	// Verify that mockCategory is not in database
	var aInDB *models.Article
	tx := e.DB.Find(&aInDB, mockArticle.ID)
	if tx.RowsAffected != 0 {
		t.Fatalf("Expected %v, got %v", 0, tx.RowsAffected)
	}
}
