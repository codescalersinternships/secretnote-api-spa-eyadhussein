package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/api/middlewares"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	validRegister = models.NewRegisterUserRequest(
		"testuser",
		"test@example.com",
		"password12345",
		"password12345",
	)

	inValidRegister = models.NewRegisterUserRequest(
		"testuser",
		"test@example.com",
		"password12345",
		"password1234",
	)

	validLogin       = models.NewLoginUserRequest("testuser", "password12345")
	inValidLogin     = models.NewLoginUserRequest("", "password1234")
	userNotFound     = models.NewLoginUserRequest("testus", "password1234")
	passwordMismatch = models.NewLoginUserRequest("testuser", "password1234")

	validNote1 = models.NewCreateNoteRequest("test note", "this is a test note", 2, time.Now().Add(time.Hour*1))
	validNote2 = models.NewCreateNoteRequest("test note 2", "this is a test note 2", 3, time.Now().Add(time.Hour*2))

	inValidNote          = models.NewCreateNoteRequest("", "this is a test note", 2, time.Now())
	maxViewsLessThanOne  = models.NewCreateNoteRequest("test note", "this is a test note", 0, validNote1.ExpiresAt.Add(time.Hour*2))
	expiresDateNotFuture = models.NewCreateNoteRequest("test note", "this is a test note", 1, validNote1.ExpiresAt.Add(-time.Hour*2))
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

func mockVerifyToken(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, util.NewResponseError(err, http.StatusBadRequest))
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

func TestServer_handleRegisterUser(t *testing.T) {
	store := storage.NewMemory()
	s := NewMockServer(store)
	s.router.POST("/api/register", s.handleRegisterUser)

	registerRequestTests := []struct {
		name     string
		request  *models.RegisterUserRequest
		expected int
	}{
		{"register user successfully", validRegister, http.StatusCreated},
		{"invalid register request", inValidRegister, http.StatusBadRequest},
	}
	for _, tt := range registerRequestTests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := createRequest("POST", "/api/register", tt.request.String())
			assertNoError(t, err)

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			assertEqual(t, tt.expected, w.Code)
		})
	}

	t.Run("user already exists", func(t *testing.T) {
		req, err := createRequest("POST", "/api/register", validRegister.String())
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusBadRequest, w.Code)
	})

	t.Run("cookies set successfully", func(t *testing.T) {
		_ = s.store.Clear()
		req, err := createRequest("POST", "/api/register", validRegister.String())
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, "token", w.Result().Cookies()[0].Value)
		assertEqual(t, "testuser", w.Result().Cookies()[1].Value)
	})
}

func TestServer_handleLoginUser(t *testing.T) {
	store := storage.NewMemory()
	_ = store.CreateUser(models.NewUser(validRegister.Username, validRegister.Email, validRegister.Password))
	s := NewMockServer(store)
	s.router.POST("/api/login", s.handleLoginUser)

	loginRequestTests := []struct {
		name     string
		request  *models.LoginUserRequest
		expected int
	}{
		{"login user successfully", validLogin, http.StatusOK},
		{"invalid login request", inValidLogin, http.StatusBadRequest},
		{"user not found", userNotFound, http.StatusNotFound},
		{"password mismatch", passwordMismatch, http.StatusUnauthorized},
	}

	for _, tt := range loginRequestTests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createRequest("POST", "/api/login", tt.request.String())
			assertNoError(t, err)

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			assertEqual(t, tt.expected, w.Code)
		})
	}
}

func TestServer_handleLogoutUser(t *testing.T) {
	s := NewServer("", storage.NewMemory())
	s.router.POST("/api/logout", s.handleLogoutUser)

	t.Run("logout user successfully", func(t *testing.T) {
		req, err := createRequest("POST", "/api/logout", "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusOK, w.Code)
	})

	t.Run("cookies cleared successfully", func(t *testing.T) {

		req, err := createRequest("POST", "/api/logout", "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, "", w.Result().Cookies()[0].Value)
		assertEqual(t, "", w.Result().Cookies()[1].Value)
	})

	t.Run("user is not authorized", func(t *testing.T) {
		s := NewServer("", storage.NewMemory())

		s.router.POST("/api/logout", mockVerifyToken, s.handleLogoutUser)

		req, err := createRequest("POST", "/api/logout", "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("user is authorized", func(t *testing.T) {
		s := NewServer("", storage.NewMemory())
		_ = s.store.CreateUser(models.NewUser(validRegister.Username, validRegister.Email, validRegister.Password))

		s.router.POST("/api/logout", mockVerifyToken, s.handleLogoutUser)

		req, err := createRequest("POST", "/api/logout", "")
		assertNoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "token",
		})
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "testuser",
		})

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusOK, w.Code)
	})
}

func TestServer_handleCreateNote(t *testing.T) {
	store := storage.NewMemory()
	s := NewServer("", store)
	_ = s.store.CreateUser(models.NewUser(validRegister.Username, validRegister.Email, validRegister.Password))
	s.router.POST("/api/notes", mockVerifyToken, middlewares.JwtAuthMiddleware(s.store), s.handleCreateNote)

	noteTests := []struct {
		name           string
		request        *models.CreateNoteRequest
		expectedStatus int
		expectedErr    error
	}{
		{"create note successfully", validNote1, http.StatusCreated, nil},
		{"invalid create note request", inValidNote, http.StatusBadRequest, util.ErrBadRequest},
		{"max views less than one", maxViewsLessThanOne, http.StatusBadRequest, util.ErrMaxViewsLessThanOne},
		{"note expired", expiresDateNotFuture, http.StatusBadRequest, util.ErrExpiresAtBeforeNow},
	}

	for _, tt := range noteTests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createRequest("POST", "/api/notes", tt.request.String())
			assertNoError(t, err)

			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: "token",
			})
			req.AddCookie(&http.Cookie{
				Name:  "user",
				Value: "testuser",
			})

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)
			assertEqual(t, tt.expectedStatus, w.Code)

			if tt.expectedErr != nil {
				var response *util.ResponseError
				err = json.NewDecoder(w.Body).Decode(&response)
				assertNoError(t, err)
				assertEqual(t, tt.expectedErr.Error(), response.Message)
			}
		})
	}
}

func TestServer_handleGetNoteByID(t *testing.T) {
	store := storage.NewMemory()
	s := NewServer("", store)
	_ = s.store.CreateNote(models.NewNote(validNote1.Title, validNote1.Content, validNote1.MaxViews, validNote1.ExpiresAt))

	s.router.GET("/api/notes/:id", s.handleGetNoteByID)

	validID := uuid.UUID([]byte(strconv.Itoa(int(store.NotesIDCounter - 1))))
	inValidID := uuid.UUID([]byte(strconv.Itoa(int(store.NotesIDCounter + 1))))
	t.Run("get note by id successfully", func(t *testing.T) {
		req, err := createRequest("GET", "/api/notes/"+validID.String(), "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusOK, w.Code)
	})

	t.Run("note not found", func(t *testing.T) {
		req, err := createRequest("GET", "/api/notes/"+inValidID.String(), "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusNotFound, w.Code)
	})

	t.Run("note has expired", func(t *testing.T) {
		_ = s.store.Clear()
		err := s.store.CreateNote(models.NewNote(validNote1.Title, validNote1.Content, validNote1.MaxViews, validNote1.ExpiresAt))
		assertNoError(t, err)

		note, err := s.store.GetNoteByID(validID.String())
		assertNoError(t, err)
		note.ExpiresAt = time.Now().Add(-time.Hour * 2)

		req, err := createRequest("GET", "/api/notes/"+validID.String(), "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusNotFound, w.Code)

		_, err = s.store.GetNoteByID(validID.String())
		assertEqual(t, util.ErrNotFound, err)
	})

	t.Run("note reached max views", func(t *testing.T) {
		_ = s.store.Clear()
		err := s.store.CreateNote(models.NewNote(validNote1.Title, validNote1.Content, validNote1.MaxViews, validNote1.ExpiresAt))
		assertNoError(t, err)

		note, err := s.store.GetNoteByID(validID.String())
		assertNoError(t, err)
		note.CurrentViews = 2

		req, err := createRequest("GET", "/api/notes/"+validID.String(), "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusNotFound, w.Code)

		_, err = s.store.GetNoteByID(validID.String())
		assertEqual(t, util.ErrNotFound, err)
	})

	t.Run("note views incremented successfully", func(t *testing.T) {
		_ = s.store.Clear()
		err := s.store.CreateNote(models.NewNote(validNote1.Title, validNote1.Content, validNote1.MaxViews, validNote1.ExpiresAt))
		assertNoError(t, err)

		note, err := s.store.GetNoteByID(validID.String())
		assertNoError(t, err)
		note.CurrentViews = 0

		req, err := createRequest("GET", "/api/notes/"+validID.String(), "")
		assertNoError(t, err)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusOK, w.Code)

		assertEqual(t, 1, note.CurrentViews)
	})
}

func TestServer_handleGetNotesByUserID(t *testing.T) {
	store := storage.NewMemory()
	s := NewServer("", store)
	s.router.GET("/api/users/notes", mockVerifyToken, middlewares.JwtAuthMiddleware(s.store), s.handleGetNotesByUserID)

	t.Run("get notes by user id successfully", func(t *testing.T) {
		_ = s.store.CreateUser(models.NewUser(validRegister.Username, validRegister.Email, validRegister.Password))
		_ = s.store.CreateNote(models.NewNote(validNote1.Title, validNote1.Content, validNote1.MaxViews, validNote1.ExpiresAt))
		_ = s.store.CreateNote(models.NewNote(validNote2.Title, validNote2.Content, validNote2.MaxViews, validNote2.ExpiresAt))
		validID1 := uuid.UUID([]byte(strconv.Itoa(int(store.NotesIDCounter - 2))))
		validID2 := uuid.UUID([]byte(strconv.Itoa(int(store.NotesIDCounter - 1))))

		note1, err := s.store.GetNoteByID(validID1.String())
		assertNoError(t, err)
		note1.UserID = store.UsersIDCounter - 1
		note2, err := s.store.GetNoteByID(validID2.String())
		assertNoError(t, err)
		note2.UserID = store.UsersIDCounter - 1
		req, err := createRequest("GET", "/api/users/notes", "")
		assertNoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "token",
		})
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "testuser",
		})

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusOK, w.Code)

		var notes *[]util.APINote
		err = json.NewDecoder(w.Body).Decode(&notes)
		assertNoError(t, err)

		assertEqual(t, 2, len(*notes))
	})

	t.Run("get notes by user id unsuccessfully", func(t *testing.T) {
		req, err := createRequest("GET", "/api/users/notes", "")
		assertNoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "token",
		})
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "testus",
		})

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assertEqual(t, http.StatusNotFound, w.Code)
	})

	t.Run("no token unauthorized", func(t *testing.T) {
		req, err := createRequest("GET", "/api/users/notes", "")
		assertNoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "testus",
		})

		w := httptest.NewRecorder()
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
