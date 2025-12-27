package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/mock"

	"github.com/luisfucros/go-events-api-tutorial/internal/configs"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
	"github.com/luisfucros/go-events-api-tutorial/internal/store/cache"
	"github.com/luisfucros/go-events-api-tutorial/internal/store/mocks"
)

// Using mocks from internal/store/mocks package for tests - created in /internal/store/mocks/mocks.go

// --- helpers ---

func newTestApp() *application {
	logger := zap.NewNop().Sugar()
	cfg := configs.Config{}
	cfg.JWT.Secret = "test-secret"
	cfg.Redis.Enabled = false

	app := &application{
		config: cfg,
		logger: logger,
	}
	// default empty stores to avoid nil pointer
	app.store = store.Storage{
		Events:    &mocks.MockEvents{},
		Users:     &mocks.MockUsers{},
		Attendees: &mocks.MockAttendees{},
	}
	app.cacheStorage = cache.Storage{
		Events:    &mocks.MockCacheEvents{},
		Users:     &mocks.MockCacheUsers{},
		Attendees: &mocks.MockCacheAttendees{},
	}
	gin.SetMode(gin.TestMode)
	return app
}

func performRequest(app *application, method, path string, body []byte, token string) *httptest.ResponseRecorder {
	r := app.routes()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	return rw
}

func createToken(secret string, id int64) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userId": float64(id), "expr": time.Now().Add(time.Hour).Unix()})
	return t.SignedString([]byte(secret))
}

func futureDate() string {
	return time.Now().AddDate(1, 0, 0).Format("2006-01-02")
}

// --- tests ---

func TestGetAllEvents_OK(t *testing.T) {
	app := newTestApp()
	me := app.store.Events.(*mocks.MockEvents)
	me.On("GetAll", mock.Anything).Return([]*store.Event{{Id: 1, Name: "e1"}}, nil)

	rw := performRequest(app, "GET", "/api/v1/events", nil, "")
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rw.Code, rw.Body.String())
	}
	var resp []*store.Event
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(resp) != 1 || resp[0].Id != 1 {
		t.Fatalf("unexpected body: %v", resp)
	}
	me.AssertExpectations(t)
}

func TestGetEvent_OK_NotFound_BadID(t *testing.T) {
	app := newTestApp()
	me := app.store.Events.(*mocks.MockEvents)
	me.On("Get", mock.Anything, int64(1)).Return(&store.Event{Id: 1, Name: "ok"}, nil)
	me.On("Get", mock.Anything, int64(2)).Return(nil, nil)

	// good
	rw := performRequest(app, "GET", "/api/v1/events/1", nil, "")
	if rw.Code != http.StatusOK { t.Fatalf("expected 200 got %d", rw.Code) }

	// not found
	rw = performRequest(app, "GET", "/api/v1/events/2", nil, "")
	if rw.Code != http.StatusNotFound { t.Fatalf("expected 404 got %d", rw.Code) }

	// bad id
	rw = performRequest(app, "GET", "/api/v1/events/xxx", nil, "")
	if rw.Code != http.StatusBadRequest { t.Fatalf("expected 400 got %d", rw.Code) }

	me.AssertExpectations(t)
}

func TestCreateEvent_Auth_Success(t *testing.T) {
	app := newTestApp()
	mu := app.store.Users.(*mocks.MockUsers)
	mu.On("Get", mock.Anything, int64(1)).Return(&store.User{Id: 1}, nil)

	me := app.store.Events.(*mocks.MockEvents)
	me.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args.Get(1).(*store.Event).Id = 10
	}).Return(nil)

	token, err := createToken(app.config.JWT.Secret, 1)
	if err != nil { t.Fatalf("token: %v", err) }

	evt := store.Event{
		Name:        "New event",
		Description: "this is a long enough description",
		Date:        futureDate(),
		Location:    "Testville",
	}
	b, _ := json.Marshal(evt)
	rw := performRequest(app, "POST", "/api/v1/events", b, token)
	if rw.Code != http.StatusCreated { t.Fatalf("expected 201 got %d body=%s", rw.Code, rw.Body.String()) }
	var resp store.Event
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil { t.Fatalf("invalid json: %v", err) }
	if resp.OwnerId != 1 || resp.Id != 10 { t.Fatalf("unexpected resp: %v", resp) }
	me.AssertExpectations(t)
	mu.AssertExpectations(t)
}

func TestAddAttendeeToEvent_Success(t *testing.T) {
	app := newTestApp()
	mu := app.store.Users.(*mocks.MockUsers)
	mu.On("Get", mock.Anything, int64(1)).Return(&store.User{Id: 1}, nil)
	mu.On("Get", mock.Anything, int64(2)).Return(&store.User{Id: 2}, nil)

	me := app.store.Events.(*mocks.MockEvents)
	me.On("Get", mock.Anything, int64(1)).Return(&store.Event{Id: 1, OwnerId: 1}, nil)

	ma := app.store.Attendees.(*mocks.MockAttendees)
	ma.On("GetByEventAndAttendee", mock.Anything, int64(1), int64(2)).Return(nil, nil)
	ma.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) { args.Get(1).(*store.Attendee).Id = 33 }).Return(&store.Attendee{Id: 33}, nil)

	token, _ := createToken(app.config.JWT.Secret, 1)
	path := "/api/v1/events/" + strconv.FormatInt(1, 10) + "/attendees/" + strconv.FormatInt(2, 10)
	rw := performRequest(app, "POST", path, nil, token)
	if rw.Code != http.StatusCreated { t.Fatalf("expected 201 got %d body=%s", rw.Code, rw.Body.String()) }
	var resp store.Attendee
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil { t.Fatalf("invalid json: %v", err) }
	if resp.Id != 33 { t.Fatalf("unexpected attendee: %v", resp) }
	ma.AssertExpectations(t)
	me.AssertExpectations(t)
	mu.AssertExpectations(t)
}

func TestGetAttendeesForEvent_OK(t *testing.T) {
	app := newTestApp()
	ma := app.store.Attendees.(*mocks.MockAttendees)
	ma.On("GetAttendeesByEvent", mock.Anything, int64(1)).Return([]*store.User{{Id: 1, Name: "u1"}}, nil)

	rw := performRequest(app, "GET", "/api/v1/events/1/attendees/", nil, "")
	if rw.Code != http.StatusOK { t.Fatalf("expected 200 got %d", rw.Code) }
	var resp []*store.User
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil { t.Fatalf("invalid json: %v", err) }
	if len(resp) != 1 || resp[0].Id != 1 { t.Fatalf("unexpected: %v", resp) }
	ma.AssertExpectations(t)
}

func TestRegisterAndLogin(t *testing.T) {
	app := newTestApp()

	// Register: no existing user
	mu := app.store.Users.(*mocks.MockUsers)
	mu.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	mu.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) { args.Get(1).(*store.User).Id = 5 }).Return(nil)

	reg := map[string]string{"email": "test@example.com", "password": "password123", "name": "Test User"}
	b, _ := json.Marshal(reg)
	rw := performRequest(app, "POST", "/api/v1/auth/register", b, "")
	if rw.Code != http.StatusCreated { t.Fatalf("register expected 201 got %d body=%s", rw.Code, rw.Body.String()) }
	var created store.User
	if err := json.Unmarshal(rw.Body.Bytes(), &created); err != nil { t.Fatalf("invalid json: %v", err) }
	if created.Id != 5 { t.Fatalf("unexpected created user: %v", created) }
	mu.AssertExpectations(t)

	// Login: existing user - reset the users mock to avoid previous expectation interference
	newMu := &mocks.MockUsers{}
	app.store.Users = newMu

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	newMu.On("GetByEmail", mock.Anything, "test@example.com").Return(&store.User{Id: 5, Email: "test@example.com", Password: string(hash)}, nil)

	login := map[string]string{"email": "test@example.com", "password": "password123"}
	b, _ = json.Marshal(login)
	rw = performRequest(app, "POST", "/api/v1/auth/login", b, "")
	if rw.Code != http.StatusOK { t.Fatalf("login expected 200 got %d body=%s", rw.Code, rw.Body.String()) }
	var lr map[string]string
	if err := json.Unmarshal(rw.Body.Bytes(), &lr); err != nil { t.Fatalf("invalid json: %v", err) }
	if lr["token"] == "" { t.Fatalf("token not returned") }
	newMu.AssertExpectations(t)
}

func TestAuthMiddleware_Unauthorized(t *testing.T) {
	app := newTestApp()
	// no token
	rw := performRequest(app, "POST", "/api/v1/events", nil, "")
	if rw.Code != http.StatusUnauthorized { t.Fatalf("expected 401 got %d body=%s", rw.Code, rw.Body.String()) }
}

func TestUpdateEvent_SuccessAndForbidden(t *testing.T) {
	app := newTestApp()
	mu := app.store.Users.(*mocks.MockUsers)
	// authenticated user lookup used by the middleware
	mu.On("Get", mock.Anything, int64(1)).Return(&store.User{Id: 1}, nil)

	me := app.store.Events.(*mocks.MockEvents)
	// existing owned event
	me.On("Get", mock.Anything, int64(1)).Return(&store.Event{Id: 1, OwnerId: 1, Name: "old", Description: "descdesc", Date: futureDate(), Location: "L"}, nil)
	me.On("Update", mock.Anything, mock.Anything).Return(nil)

	// success (owner)
	token, _ := createToken(app.config.JWT.Secret, 1)
	dateStr := futureDate()
	b, _ := json.Marshal(map[string]string{"name": "updated", "description": "longer description", "date": dateStr, "location": "here"})
	rw := performRequest(app, "PUT", "/api/v1/events/1", b, token)
	if rw.Code != http.StatusOK { t.Fatalf("expected 200 got %d body=%s", rw.Code, rw.Body.String()) }

	// forbidden (not owner)
	me.On("Get", mock.Anything, int64(2)).Return(&store.Event{Id: 2, OwnerId: 99}, nil)
	token2, _ := createToken(app.config.JWT.Secret, 1)
	rw = performRequest(app, "PUT", "/api/v1/events/2", b, token2)
	if rw.Code != http.StatusForbidden { t.Fatalf("expected 403 got %d body=%s", rw.Code, rw.Body.String()) }

	me.AssertExpectations(t)
	mu.AssertExpectations(t)
}

func TestDeleteEvent_SuccessAndForbidden(t *testing.T) {
	app := newTestApp()
	mu := app.store.Users.(*mocks.MockUsers)
	mu.On("Get", mock.Anything, int64(1)).Return(&store.User{Id: 1}, nil)

	me := app.store.Events.(*mocks.MockEvents)
	me.On("Get", mock.Anything, int64(1)).Return(&store.Event{Id: 1, OwnerId: 1}, nil)
	me.On("Delete", mock.Anything, int64(1)).Return(nil)

	token, _ := createToken(app.config.JWT.Secret, 1)
	rw := performRequest(app, "DELETE", "/api/v1/events/1", nil, token)
	if rw.Code != http.StatusNoContent { t.Fatalf("expected 204 got %d body=%s", rw.Code, rw.Body.String()) }

	// forbidden
	me.On("Get", mock.Anything, int64(2)).Return(&store.Event{Id: 2, OwnerId: 99}, nil)
	rw = performRequest(app, "DELETE", "/api/v1/events/2", nil, token)
	if rw.Code != http.StatusForbidden { t.Fatalf("expected 403 got %d body=%s", rw.Code, rw.Body.String()) }

	me.AssertExpectations(t)
	mu.AssertExpectations(t)
}

func TestAddAttendee_Conflict(t *testing.T) {
	app := newTestApp()
	me := app.store.Events.(*mocks.MockEvents)
	me.On("Get", mock.Anything, int64(1)).Return(&store.Event{Id: 1, OwnerId: 1}, nil)
	mu := app.store.Users.(*mocks.MockUsers)
	// authenticated user lookup (token user id 1)
	mu.On("Get", mock.Anything, int64(1)).Return(&store.User{Id: 1}, nil)
	mu.On("Get", mock.Anything, int64(2)).Return(&store.User{Id: 2}, nil)
	ma := app.store.Attendees.(*mocks.MockAttendees)
	ma.On("GetByEventAndAttendee", mock.Anything, int64(1), int64(2)).Return(&store.Attendee{Id: 5}, nil)

	token, _ := createToken(app.config.JWT.Secret, 1)
	path := "/api/v1/events/" + strconv.FormatInt(1, 10) + "/attendees/" + strconv.FormatInt(2, 10)
	rw := performRequest(app, "POST", path, nil, token)
	if rw.Code != http.StatusConflict { t.Fatalf("expected 409 got %d body=%s", rw.Code, rw.Body.String()) }

	me.AssertExpectations(t)
	mu.AssertExpectations(t)
	ma.AssertExpectations(t)
}

func TestGetEventsByAttendee_CacheHitAndMiss(t *testing.T) {
	app := newTestApp()
	app.config.Redis.Enabled = true
	cacheAtt := app.cacheStorage.Attendees.(*mocks.MockCacheAttendees)
	cacheAtt.On("GetEventsByAttendee", mock.Anything, int64(1)).Return([]store.Event{{Id: 5}}, nil)

	rw := performRequest(app, "GET", "/api/v1/attendees/1/events", nil, "")
	if rw.Code != http.StatusOK { t.Fatalf("expected 200 got %d body=%s", rw.Code, rw.Body.String()) }
	var resp []*store.Event
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil { t.Fatalf("invalid json: %v", err) }
	if len(resp) != 1 || resp[0].Id != 5 { t.Fatalf("unexpected: %v", resp) }
	cacheAtt.AssertExpectations(t)

	// cache miss -> backend hit on a fresh app instance
	app2 := newTestApp()
	app2.config.Redis.Enabled = true
	cacheAtt2 := app2.cacheStorage.Attendees.(*mocks.MockCacheAttendees)
	cacheAtt2.On("GetEventsByAttendee", mock.Anything, int64(2)).Return(nil, nil)
	// cache will be set after backend fetch
	cacheAtt2.On("SetEventsByAttendee", mock.Anything, int64(2), mock.Anything).Return(nil)
	ma2 := app2.store.Attendees.(*mocks.MockAttendees)
	ma2.On("GetEventsByAttendee", mock.Anything, int64(2)).Return([]*store.Event{{Id: 6}}, nil)

	rw = performRequest(app2, "GET", "/api/v1/attendees/2/events", nil, "")
	if rw.Code != http.StatusOK { t.Fatalf("expected 200 got %d body=%s", rw.Code, rw.Body.String()) }
	var resp2 []*store.Event
	if err := json.Unmarshal(rw.Body.Bytes(), &resp2); err != nil { t.Fatalf("invalid json: %v", err) }
	if len(resp2) != 1 || resp2[0].Id != 6 { t.Fatalf("unexpected: %v", resp2) }
	ma2.AssertExpectations(t)
	cacheAtt2.AssertExpectations(t)
}

func TestGetEvent_CacheHit(t *testing.T) {
	app := newTestApp()
	app.config.Redis.Enabled = true
	ce := app.cacheStorage.Events.(*mocks.MockCacheEvents)
	ce.On("Get", mock.Anything, int64(1)).Return(&store.Event{Id: 1, Name: "cached"}, nil)

	rw := performRequest(app, "GET", "/api/v1/events/1", nil, "")
	if rw.Code != http.StatusOK { t.Fatalf("expected 200 got %d body=%s", rw.Code, rw.Body.String()) }
	var resp store.Event
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil { t.Fatalf("invalid json: %v", err) }
	if resp.Id != 1 { t.Fatalf("unexpected: %v", resp) }
	ce.AssertExpectations(t)
}

func TestDeleteAttendeeFromEvent_SuccessAndForbidden(t *testing.T) {
	app := newTestApp()
	mu := app.store.Users.(*mocks.MockUsers)
	mu.On("Get", mock.Anything, int64(1)).Return(&store.User{Id: 1}, nil)

	me := app.store.Events.(*mocks.MockEvents)
	me.On("Get", mock.Anything, int64(1)).Return(&store.Event{Id: 1, OwnerId: 1}, nil)

	ma := app.store.Attendees.(*mocks.MockAttendees)
	ma.On("Delete", mock.Anything, int64(2), int64(1)).Return(nil)

	token, _ := createToken(app.config.JWT.Secret, 1)
	path := "/api/v1/events/1/attendees/2"
	rw := performRequest(app, "DELETE", path, nil, token)
	if rw.Code != http.StatusNoContent { t.Fatalf("expected 204 got %d body=%s", rw.Code, rw.Body.String()) }

	// forbidden
	me.On("Get", mock.Anything, int64(2)).Return(&store.Event{Id: 2, OwnerId: 99}, nil)
	rw = performRequest(app, "DELETE", "/api/v1/events/2/attendees/2", nil, token)
	if rw.Code != http.StatusForbidden { t.Fatalf("expected 403 got %d body=%s", rw.Code, rw.Body.String()) }

	ma.AssertExpectations(t)
	me.AssertExpectations(t)
	mu.AssertExpectations(t)
}

func TestGetAllEvents_CacheHit(t *testing.T) {
	app := newTestApp()
	app.config.Redis.Enabled = true
	ce := app.cacheStorage.Events.(*mocks.MockCacheEvents)
	ce.On("GetAll", mock.Anything).Return([]store.Event{{Id: 3}}, nil)

	rw := performRequest(app, "GET", "/api/v1/events", nil, "")
	if rw.Code != http.StatusOK { t.Fatalf("expected 200 got %d body=%s", rw.Code, rw.Body.String()) }
	var resp []*store.Event
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil { t.Fatalf("invalid json: %v", err) }
	if len(resp) != 1 || resp[0].Id != 3 { t.Fatalf("unexpected: %v", resp) }
	ce.AssertExpectations(t)
}
