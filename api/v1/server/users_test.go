package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/JonathanGzzBen/ingenialists/api/v1/repository/mocks"
	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
)

var mockUsers = []models.User{
	{ID: 123, Name: "First User", Role: models.RoleReader},
	{ID: 456, Name: "Second User", Role: models.RoleReader},
	{ID: 789, Name: "Third User", Role: models.RoleReader},
}

func TestGetAllUsers(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, u := range mockUsers {
		s.UsersRepo.CreateUser(&u)
	}

	mockUsersRepo := &mocks.UsersRepository{}
	mockUsersRepo.On("GetAllUsers").Return(mockUsers, nil)
	s.UsersRepo = mockUsersRepo

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
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	uToGet := mockUsers[1]

	mockUsersRepo := &mocks.UsersRepository{}
	mockUsersRepo.On("GetUser", uToGet.ID).Return(&uToGet, nil)
	s.UsersRepo = mockUsersRepo

	res, err := http.Get(fmt.Sprintf("%s/v1/users/%d", ts.URL, uToGet.ID))
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
	if uToGet != resUser {
		t.Fatalf("Expected %v, got %v", uToGet, resUser)
	}
}

func TestUpdateUserChangeNameAsAdministratorReturnOkDontMakeChanges(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, u := range mockUsers {
		s.UsersRepo.CreateUser(&u)
	}
	mockUser := mockUsers[1]
	mockUser.Name = "User Updated"

	muJSONBytes, err := json.Marshal(mockUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/users/%d", ts.URL, mockUser.ID), bytes.NewBuffer(muJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
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

	uInDb, err := s.UsersRepo.GetUser(mockUser.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if uInDb.Name == mockUser.Name {
		t.Fatalf("Expected %v, got %v", uInDb.Name, mockUser.Name)
	}
}

func TestUpdateUserChangeRoleAsAdministratorReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, u := range mockUsers {
		s.UsersRepo.CreateUser(&u)
	}
	mockUser := mockUsers[1]
	mockUser.Role = models.RoleAdministrator

	muJSONBytes, err := json.Marshal(mockUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/users/%d", ts.URL, mockUser.ID), bytes.NewBuffer(muJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
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

	uInDb, err := s.UsersRepo.GetUser(mockUser.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(uInDb.Role) != string(mockUser.Role) {
		t.Fatalf("Expected %v, got %v", mockUser.Role, uInDb.Role)
	}
}

// TestUpdateUserChangeNameAsDifferentUserReturnForbidden tests a request
// in which a user with a different role from Administrator tries to update
// a user with a different ID than his own.
//
// In testing mode, authenticated user will have ID = 1.
func TestUpdateUserChangeNameAsDifferentUserReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, u := range mockUsers {
		s.UsersRepo.CreateUser(&u)
	}
	// mockUser has ID different from 1
	mockUser := mockUsers[1]
	mockUser.Name = "Updated name"

	muJSONBytes, err := json.Marshal(mockUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/users/%d", ts.URL, mockUser.ID), bytes.NewBuffer(muJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %d, got %v", http.StatusForbidden, res.StatusCode)
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

	uInDB, err := s.UsersRepo.GetUser(mockUser.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// If user was updated in db
	if uInDB.Name != mockUsers[1].Name {
		t.Fatalf("Expected %v, got %v", mockUsers[1].Name, uInDB.Name)
	}
}

// TestUpdateUserChangeNameAsSameUserReturnOk tests a request
// in which a user with a different role from Administrator tries to update
// a user with a different ID than his own.
//
// In testing mode, authenticated user will have ID = 1.
func TestUpdateUserChangeNameAsSameUserReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	mockUsers[1].ID = 1
	for _, u := range mockUsers {
		s.UsersRepo.CreateUser(&u)
	}
	mockUser := mockUsers[1]
	mockUser.Name = "Updated name"

	muJSONBytes, err := json.Marshal(mockUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/users/%d", ts.URL, mockUser.ID), bytes.NewBuffer(muJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "AccessToken")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %v", http.StatusOK, res.StatusCode)
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

	uInDB, err := s.UsersRepo.GetUser(mockUser.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if uInDB.Name != mockUser.Name {
		t.Fatalf("Expected %v, got %v", mockUser.Name, uInDB.Name)
	}
}
