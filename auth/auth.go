package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bogudini/lw-api/auth/models"
	"github.com/go-openapi/strfmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Authenticator interface {
	SetAuth(req *http.Request) error
}

type SimpleAuth struct{ ts TokenSource }

func NewSimpleAuth(source TokenSource) *SimpleAuth {
	return &SimpleAuth{ts: source}
}

func (sa SimpleAuth) SetAuth(req *http.Request) error {
	token, err := sa.ts.Token(req.Context())
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token.String())
	return nil
}

type App struct {
	ApplicationId     string `json:"applicationId"`
	ApplicationSecret string `json:"applicationSecret"`
	Token             string `json:"token"`
}

type Credentials struct {
	App          App
	RefreshToken string
}

// ----------------------------------------------------------------------------
// ========== 2) TokenSource, который ходит за новым access-token ==========
type TokenSource interface {
	Token(ctx context.Context) (strfmt.UUID, error) // вернуть актуальный access-token (с кешированием)
	ForceRefresh(ctx context.Context) error         // принудительно обновить (на 401, например)
}

// Реализация через refresh-эндпоинт
type RefreshTokenSource struct {
	creds   Credentials
	baseURL *url.URL
	hc      *http.Client

	mu    sync.Mutex
	token strfmt.UUID
	exp   time.Time
}

// Конструктор
func NewRefreshTokenSource(creds Credentials) *RefreshTokenSource {
	hc := &http.Client{Timeout: 30 * time.Second}
	parsedAuthURL, _ := url.Parse("https://api.linnworks.net")
	return &RefreshTokenSource{hc: hc, baseURL: parsedAuthURL, creds: creds}
}

func (r *RefreshTokenSource) Token(ctx context.Context) (strfmt.UUID, error) {
	r.mu.Lock()
	// если токен ещё валиден с запасом — отдать из кеша
	if r.token != "" && time.Until(r.exp) > 30*time.Second {
		t := r.token
		r.mu.Unlock()
		return t, nil
	}
	r.mu.Unlock()
	// иначе — обновить
	if err := r.refresh(ctx); err != nil {
		return "", err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.token, nil
}

func (r *RefreshTokenSource) ForceRefresh(ctx context.Context) error {
	return r.refresh(ctx)
}

func (r *RefreshTokenSource) refresh(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Собираем URL: /auth/refresh — поменяй на свой путь
	u := r.baseURL.ResolveReference(&url.URL{Path: "/api/Auth/AuthorizeByApplication"})
	body := map[string]string{
		"applicationId":     r.creds.App.ApplicationId,
		"applicationSecret": r.creds.App.ApplicationSecret,
		"token":             r.creds.App.Token,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("refresh failed: %d %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	var rr models.BaseSession
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return err
	}
	if rr.Token == "" {
		return errors.New("refresh returned empty access token")
	}
	r.token = rr.Token
	// Если API отдаёт ExpiresIn, рассчитываем exp, иначе — ставим разумный TTL (например, 10 мин)
	if rr.TTL > 0 {
		r.exp = time.Now().Add(time.Duration(rr.TTL) * time.Second)
	} else {
		r.exp = time.Now().Add(10 * time.Minute)
	}
	return nil
}
