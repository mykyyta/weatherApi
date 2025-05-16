package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"weatherApi/internal/model"
)

func setupTestRouterWithDB(t *testing.T) *gin.Engine {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}

	err = db.AutoMigrate(&model.Subscription{})
	if err != nil {
		t.Fatalf("failed to migrate test DB: %v", err)
	}

	SetDB(db)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterRoutes(r)
	return r
}

func TestSubscribe_Success(t *testing.T) {
	router := setupTestRouterWithDB(t)

	form := url.Values{}
	form.Add("email", "test@example.com")
	form.Add("city", "Kyiv")
	form.Add("frequency", "daily")

	req := httptest.NewRequest(http.MethodPost, "/api/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"message":"Subscription successful. Confirmation email sent."}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestSubscribe_MissingEmail(t *testing.T) {
	router := setupTestRouterWithDB(t)

	form := url.Values{}
	form.Add("city", "Kyiv")
	form.Add("frequency", "daily")

	req := httptest.NewRequest(http.MethodPost, "/api/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input") // залежить від твоєї валідації
}

func TestSubscribe_InvalidFrequency(t *testing.T) {
	router := setupTestRouterWithDB(t)

	form := url.Values{}
	form.Add("email", "test@example.com")
	form.Add("city", "Kyiv")
	form.Add("frequency", "weekly")

	req := httptest.NewRequest(http.MethodPost, "/api/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input") // залежить від твоєї помилки
}

func TestSubscribe_DuplicateEmail(t *testing.T) {
	router := setupTestRouterWithDB(t)

	err := DB.Create(&model.Subscription{
		ID:             uuid.New().String(),
		Email:          "duplicate@example.com",
		City:           "Kyiv",
		Frequency:      "daily",
		IsConfirmed:    true,
		IsUnsubscribed: false,
		Token:          "some-token",
		CreatedAt:      time.Now(),
	}).Error
	require.NoError(t, err)

	form := url.Values{}
	form.Add("email", "duplicate@example.com")
	form.Add("city", "Kyiv")
	form.Add("frequency", "daily")

	req := httptest.NewRequest(http.MethodPost, "/api/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Email already subscribed")
}
