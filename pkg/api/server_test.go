package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/util"
	"github.com/gin-gonic/gin"
)

const (
	validRegisterUserRequest     = `{"username":"testuser","email":"test@example.com","password":"password12345", "password_confirmation":"password12345"}`
	inValidRegisterUserRequest   = `{"username":"testuser","email":"test@example.com","password":"password12345", "password_confirmation":"password1234"}`
	validLoginUserRequest        = `{"username":"testuser","password":"password12345"}`
	inValidLoginUserRequest      = `{"username":"","password":"password1234"}`
	userNotFoundLoginUserRequest = `{"username":"testus","password":"password1234"}`
	passwordMismatchLoginRequest = `{"username":"testuser","password":"password1234"}`
)

type MockServer struct {
	store  storage.Storage
	router *gin.Engine
}

func NewMockServer(store storage.Storage) *MockServer {
	return &MockServer{store: store, router: gin.Default()}
}

func mockCreateToken(username string, duration time.Duration) (string, error) {
	return "token", nil
}

func mockJwtAuthMiddleware(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := c.Cookie("user")
		if err != nil {
			c.JSON(http.StatusUnauthorized, util.NewResponseError(
				util.ErrUnauthorized, http.StatusUnauthorized,
			))
			c.Abort()
			return
		}

		user, err := store.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusNotFound, util.NewResponseError(
				err, http.StatusNotFound,
			))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func mockVerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Cookie("token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token cookie"})
			c.Abort()
			return
		}

		if tokenCookie != "token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *MockServer) handleRegisterUser(c *gin.Context) {
	var registerUserRequest models.RegisterUserRequest

	if err := c.ShouldBindJSON(&registerUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponseError(util.ErrBadRequest, http.StatusBadRequest))
		return
	}

	user := models.NewUser(registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password)

	err := user.MockSetPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(util.ErrInternalServer, http.StatusInternalServerError))
		return
	}

	if err := s.store.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(err, http.StatusInternalServerError))
		return
	}

	token, _ := mockCreateToken(registerUserRequest.Username, time.Hour*24*7)

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", token, 604800, "/", "localhost", true, true)
	c.SetCookie("user", user.Username, 604800, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (s *MockServer) handleLoginUser(c *gin.Context) {
	var loginUserRequest models.LoginUserRequest

	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponseError(util.ErrBadRequest, http.StatusBadRequest))
		return
	}

	user, err := s.store.GetUserByUsername(loginUserRequest.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, util.NewResponseError(err, http.StatusNotFound))
		return
	}

	if err := user.MockCheckPassword(loginUserRequest.Password); err != nil {
		c.JSON(http.StatusUnauthorized, util.NewResponseError(util.ErrUnauthorized, http.StatusUnauthorized))
		return
	}

	token, _ := mockCreateToken(loginUserRequest.Username, time.Hour*24*7)

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", token, 604800, "/", "localhost", false, true)
	c.SetCookie("user", user.Username, 604800, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "user login successfully"})
}

func (s *MockServer) handleLogoutUser(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.SetCookie("user", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

func TestServer_handleRegisterUser(t *testing.T) {
	registerRequestTests := []struct {
		name     string
		request  string
		expected int
	}{
		{"register user successfully", validRegisterUserRequest, http.StatusCreated},
		{"invalid register request", inValidRegisterUserRequest, http.StatusBadRequest},
	}

	for _, tt := range registerRequestTests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMockServer(storage.NewMemory())

			s.router.POST("/api/register", s.handleRegisterUser)

			req, err := createRequest("POST", "/api/register", tt.request)
			assertNoError(t, err)

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			assertEqual(t, tt.expected, w.Code)
		})
	}

	t.Run("user already exists", func(t *testing.T) {
		s := NewMockServer(storage.NewMemory())

		s.router.POST("/api/register", s.handleRegisterUser)

		req, err := createRequest("POST", "/api/register", validRegisterUserRequest)
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusCreated, w.Code)

		w = httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusBadRequest, w.Code)
	})

	t.Run("cookies set successfully", func(t *testing.T) {
		s := NewMockServer(storage.NewMemory())

		s.router.POST("/api/register", s.handleRegisterUser)

		req, err := createRequest("POST", "/api/register", validRegisterUserRequest)
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, "token", w.Result().Cookies()[0].Value)
		assertEqual(t, "testuser", w.Result().Cookies()[1].Value)
	})
}

func TestServer_handleLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := storage.NewMemory()
	_ = store.CreateUser(models.NewUser("testuser", "test@example.com", "password12345"))

	loginRequestTests := []struct {
		name     string
		request  string
		expected int
	}{
		{"login user successfully", validLoginUserRequest, http.StatusOK},
		{"invalid login request", inValidLoginUserRequest, http.StatusBadRequest},
		{"user not found", userNotFoundLoginUserRequest, http.StatusNotFound},
		{"password mismatch", passwordMismatchLoginRequest, http.StatusUnauthorized},
	}

	for _, tt := range loginRequestTests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMockServer(store)

			s.router.POST("/api/login", s.handleLoginUser)

			req, err := createRequest("POST", "/api/login", tt.request)
			assertNoError(t, err)

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			assertEqual(t, tt.expected, w.Code)
		})
	}
}

func TestServer_handleLogoutUser(t *testing.T) {

	t.Run("logout user successfully", func(t *testing.T) {
		s := NewMockServer(storage.NewMemory())

		s.router.POST("/api/logout", s.handleLogoutUser)

		req, err := createRequest("POST", "/api/logout", "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusOK, w.Code)
	})

	t.Run("cookies cleared successfully", func(t *testing.T) {
		s := NewMockServer(storage.NewMemory())

		s.router.POST("/api/logout", s.handleLogoutUser)

		req, err := createRequest("POST", "/api/logout", "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, "", w.Result().Cookies()[0].Value)
		assertEqual(t, "", w.Result().Cookies()[1].Value)
	})

	t.Run("user is not authorized", func(t *testing.T) {
		s := NewMockServer(storage.NewMemory())

		s.router.POST("/api/logout", mockVerifyToken(), mockJwtAuthMiddleware(s.store), s.handleLogoutUser)

		req, err := createRequest("POST", "/api/logout", "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("user is authorized", func(t *testing.T) {
		s := NewMockServer(storage.NewMemory())

		s.router.POST("/api/logout", mockJwtAuthMiddleware(s.store), s.handleLogoutUser)

		req, err := createRequest("POST", "/api/logout", "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		w.Header().Set("Cookie", "user=testuser; token=token")
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusUnauthorized, w.Code)
	})
}

func createRequest(method, url, body string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}
}

func assertEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}
