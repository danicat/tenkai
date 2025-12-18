package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rep-1/internal/api"
	"rep-1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of service.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUsers() []service.User {
	args := m.Called()
	return args.Get(0).([]service.User)
}

func (m *MockUserService) CreateUser(user service.User) (service.User, error) {
	args := m.Called(user)
	return args.Get(0).(service.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id string) (service.User, bool) {
	args := m.Called(id)
	return args.Get(0).(service.User), args.Bool(1)
}

func TestUserHandler_GetUsers(t *testing.T) {
	mockService := new(MockUserService)
	expectedUsers := []service.User{
		{ID: "1", Name: "John Doe", Email: "john@example.com"},
		{ID: "2", Name: "Jane Smith", Email: "jane@example.com"},
	}
	mockService.On("GetAllUsers").Return(expectedUsers)

	handler := api.NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualUsers []service.User
	err := json.NewDecoder(rr.Body).Decode(&actualUsers)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, actualUsers)
	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockService := new(MockUserService)
	newUser := service.User{Name: "John Doe", Email: "john@example.com"}
	createdUser := service.User{ID: "uuid-123", Name: "John Doe", Email: "john@example.com"}
	mockService.On("CreateUser", newUser).Return(createdUser, nil)

	handler := api.NewUserHandler(mockService)

	requestBody, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var actualUser service.User
	err := json.NewDecoder(rr.Body).Decode(&actualUser)
	assert.NoError(t, err)
	assert.Equal(t, createdUser, actualUser)
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser(t *testing.T) {
	mockService := new(MockUserService)
	expectedUser := service.User{ID: "1", Name: "John Doe", Email: "john@example.com"}
	mockService.On("GetUserByID", "1").Return(expectedUser, true)

	handler := api.NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req, "1")

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualUser service.User
	err := json.NewDecoder(rr.Body).Decode(&actualUser)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	mockService.On("GetUserByID", "non-existent").Return(service.User{}, false)

	handler := api.NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/non-existent", nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req, "non-existent")

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertExpectations(t)
}
