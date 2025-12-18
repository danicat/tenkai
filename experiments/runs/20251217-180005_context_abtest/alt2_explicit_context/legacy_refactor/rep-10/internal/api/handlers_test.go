package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"experiment/internal/api"
	"experiment/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(name, email string) (service.User, error) {
	args := m.Called(name, email)
	return args.Get(0).(service.User), args.Error(1)
}

func (m *MockUserService) GetUsers() []service.User {
	args := m.Called()
	return args.Get(0).([]service.User)
}

func (m *MockUserService) GetUser(id string) (service.User, error) {
	args := m.Called(id)
	return args.Get(0).(service.User), args.Error(1)
}

func TestUserHandlers_GetUsers(t *testing.T) {
	mockService := new(MockUserService)
	expectedUsers := []service.User{
		{ID: "1", Name: "John", Email: "john@example.com"},
		{ID: "2", Name: "Jane", Email: "jane@example.com"},
	}
	mockService.On("GetUsers").Return(expectedUsers)

	handlers := &api.UserHandlers{UserService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handlers.GetUsers(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var actualUsers []service.User
	json.NewDecoder(rr.Body).Decode(&actualUsers)
	assert.Equal(t, expectedUsers, actualUsers)
	mockService.AssertExpectations(t)
}

func TestUserHandlers_CreateUser(t *testing.T) {
	mockService := new(MockUserService)
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	expectedUser := service.User{ID: "abc-123", Name: "Alice", Email: "alice@example.com"}
	mockService.On("CreateUser", newUser.Name, newUser.Email).Return(expectedUser, nil)

	handlers := &api.UserHandlers{UserService: mockService}

	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handlers.CreateUser(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var actualUser service.User
	json.NewDecoder(rr.Body).Decode(&actualUser)
	assert.Equal(t, expectedUser, actualUser)
	mockService.AssertExpectations(t)
}

func TestUserHandlers_CreateUser_BadRequest(t *testing.T) {
	mockService := new(MockUserService)
	handlers := &api.UserHandlers{UserService: mockService}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid json")))
	rr := httptest.NewRecorder()

	handlers.CreateUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "CreateUser")
}

func TestUserHandlers_CreateUser_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	mockService.On("CreateUser", newUser.Name, newUser.Email).Return(service.User{}, errors.New("db error"))

	handlers := &api.UserHandlers{UserService: mockService}

	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handlers.CreateUser(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandlers_GetUserByID(t *testing.T) {
	mockService := new(MockUserService)
	expectedUser := service.User{ID: "abc-123", Name: "Alice", Email: "alice@example.com"}
	mockService.On("GetUser", "abc-123").Return(expectedUser, nil)

	handlers := &api.UserHandlers{UserService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/users/abc-123", nil)
	rr := httptest.NewRecorder()

	handlers.GetUserByID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var actualUser service.User
	json.NewDecoder(rr.Body).Decode(&actualUser)
	assert.Equal(t, expectedUser, actualUser)
	mockService.AssertExpectations(t)
}

func TestUserHandlers_GetUserByID_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	mockService.On("GetUser", "non-existent").Return(service.User{}, service.ErrUserNotFound)

	handlers := &api.UserHandlers{UserService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/users/non-existent", nil)
	rr := httptest.NewRecorder()

	handlers.GetUserByID(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandlers_GetUserByID_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	mockService.On("GetUser", "some-id").Return(service.User{}, errors.New("db error"))

	handlers := &api.UserHandlers{UserService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/users/some-id", nil)
	rr := httptest.NewRecorder()

	handlers.GetUserByID(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandlers_GetUserByID_NoID(t *testing.T) {
	mockService := new(MockUserService)
	handlers := &api.UserHandlers{UserService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	rr := httptest.NewRecorder()

	handlers.GetUserByID(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertNotCalled(t, "GetUser")
}
