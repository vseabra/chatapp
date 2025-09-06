package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpPkg "chatapp/internal/http"

	"github.com/gin-gonic/gin"
)

// mockService implements Service for handler tests
type mockService struct {
	registerFn func(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	loginFn    func(ctx context.Context, req LoginRequest, jwtSecret string, expiresIn string) (*LoginResponse, error)
}

func (m *mockService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	return m.registerFn(ctx, req)
}

func (m *mockService) Login(ctx context.Context, req LoginRequest, jwtSecret string, expiresIn string) (*LoginResponse, error) {
	return m.loginFn(ctx, req, jwtSecret, expiresIn)
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := httpPkg.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func TestHandleRegister_Success(t *testing.T) {
	svc := &mockService{
		registerFn: func(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
			return &RegisterResponse{ID: "id123", Name: req.Name, Email: req.Email}, nil
		},
	}
	h := NewHandler(svc)
	r := setupRouter(h)

	body, _ := json.Marshal(RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "secret123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestHandleRegister_BadRequest(t *testing.T) {
	svc := &mockService{registerFn: func(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) { return nil, nil }}
	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleLogin_Success(t *testing.T) {
	svc := &mockService{
		loginFn: func(ctx context.Context, req LoginRequest, jwtSecret string, expiresIn string) (*LoginResponse, error) {
			return &LoginResponse{AccessToken: "token"}, nil
		},
	}
	h := NewHandler(svc)
	r := setupRouter(h)

	body, _ := json.Marshal(LoginRequest{Email: "alice@example.com", Password: "secret"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandleLogin_BadRequest(t *testing.T) {
	svc := &mockService{loginFn: func(ctx context.Context, req LoginRequest, jwtSecret string, expiresIn string) (*LoginResponse, error) {
		return nil, nil
	}}
	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
