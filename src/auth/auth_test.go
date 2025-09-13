package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

type mockProducer struct {
	ProducedMessages []*kafka.Message
	ProduceError     error
}

func (m *mockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	m.ProducedMessages = append(m.ProducedMessages, msg)
	return m.ProduceError
}
func (m *mockProducer) Close() {}

func setupHandler(t *testing.T) (*AuthHandler, sqlmock.Sqlmock, *mockProducer) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	publicKey := &privateKey.PublicKey

	producer := &mockProducer{}

	handler := &AuthHandler{
		db:         db,
		jwtPrivate: privateKey,
		jwtPublic:  publicKey,
		producer:   producer,
	}

	return handler, mock, producer
}

func TestSignup_Success(t *testing.T) {
	handler, mock, _ := setupHandler(t)
	defer handler.db.Close()

	signupInfo := SignupInfo{
		Id:       "user123",
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	bodyBytes, err := json.Marshal(signupInfo)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	req.ContentLength = int64(len(bodyBytes))

	mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM "users" WHERE id = \$1\)`).
		WithArgs("user123").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(signupInfo.Id, signupInfo.Name, signupInfo.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO "passwords"`).
		WithArgs(signupInfo.Id, signupInfo.hash()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	handler.signup(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "Signup successful")
	assert.NotNil(t, resp.Cookies())
}

func TestLogin_Success(t *testing.T) {
	handler, mock, _ := setupHandler(t)
	defer handler.db.Close()

	loginInfo := LoginInfo{
		Id:       "user123",
		Password: "password123",
	}

	bodyBytes, err := json.Marshal(loginInfo)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(bodyBytes))
	req.ContentLength = int64(len(bodyBytes))

	mock.ExpectQuery(`SELECT password FROM "passwords" WHERE user_id = \$1`).
		WithArgs(loginInfo.Id).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(loginInfo.hash()))

	w := httptest.NewRecorder()
	handler.login(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "Login successful")
	assert.NotNil(t, resp.Cookies())
}

func TestLogin_InvalidPassword(t *testing.T) {
	handler, mock, _ := setupHandler(t)
	defer handler.db.Close()

	loginInfo := LoginInfo{
		Id:       "user123",
		Password: "wrongpassword",
	}

	bodyBytes, err := json.Marshal(loginInfo)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(bodyBytes))
	req.ContentLength = int64(len(bodyBytes))

	mock.ExpectQuery(`SELECT password FROM "passwords" WHERE user_id = \$1`).
		WithArgs(loginInfo.Id).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("differenthash"))

	w := httptest.NewRecorder()
	handler.login(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Contains(t, string(body), "Invalid username or password")
}

func TestWhoami_Success(t *testing.T) {
	handler, mock, _ := setupHandler(t)
	defer handler.db.Close()

	tokenString := handler.genToken("user123")

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: tokenString})

	mock.ExpectQuery(`SELECT name FROM "users" WHERE id = \$1`).
		WithArgs("user123").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Test User"))

	w := httptest.NewRecorder()
	handler.whoami(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "Hello, Test User")
}

func TestWhoami_MissingCookie(t *testing.T) {
	handler, _, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	w := httptest.NewRecorder()

	handler.whoami(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Contains(t, string(body), "Cookie is missing")
}

func TestWhoami_InvalidToken(t *testing.T) {
	handler, _, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: "invalidtoken"})

	w := httptest.NewRecorder()
	handler.whoami(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(body), "Invalid token")
}

func TestSignup_KafkaProduceError(t *testing.T) {
	handler, mock, producer := setupHandler(t)
	defer handler.db.Close()

	producer.ProduceError = errors.New("kafka error")

	signupInfo := SignupInfo{
		Id:       "user123",
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	bodyBytes, err := json.Marshal(signupInfo)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	req.ContentLength = int64(len(bodyBytes))

	mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM "users" WHERE id = \$1\)`).
		WithArgs("user123").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(signupInfo.Id, signupInfo.Name, signupInfo.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO "passwords"`).
		WithArgs(signupInfo.Id, signupInfo.hash()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	handler.signup(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, string(body), "Error sending with kafka")
	assert.Len(t, producer.ProducedMessages, 1)
}
