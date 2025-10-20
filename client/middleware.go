package client

import (
	"bytes"
	"fmt"
	"github.com/MMC-BK/lw-api/auth"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"time"
)

// ========== 4) RT middleware (декораторы транспорта) ==========
type RTMiddleware func(http.RoundTripper) http.RoundTripper

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// Auth: вставляет заголовок
func AuthMiddleware(a auth.Authenticator) RTMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if err := a.SetAuth(req); err != nil {
				return nil, err
			}
			return next.RoundTrip(req)
		})
	}
}

// Auto-refresh: если получили 401 — один раз форсим refresh и повторяем запрос
func AutoRefreshOn401(ts auth.TokenSource) RTMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			// тело запроса может потребоваться повторно прочитать
			var bodyCopy []byte
			if req.Body != nil {
				bodyCopy, _ = io.ReadAll(req.Body)
				req.Body.Close()
				req.Body = io.NopCloser(bytes.NewReader(bodyCopy))
			}
			resp, err := next.RoundTrip(req)
			if err != nil || (resp != nil && resp.StatusCode != http.StatusUnauthorized) {
				return resp, err
			}
			if resp != nil {
				resp.Body.Close()
			}

			// Пробуем форс-рефреш и повтор
			if err := ts.ForceRefresh(req.Context()); err != nil {
				return nil, fmt.Errorf("401 and refresh failed: %w", err)
			}
			// Восстанавливаем тело и повторяем
			if bodyCopy != nil {
				req.Body = io.NopCloser(bytes.NewReader(bodyCopy))
			}
			return next.RoundTrip(req)
		})
	}
}

// Прочие полезные прослойки
func RateLimitMiddleware(rps int, per time.Duration) RTMiddleware {
	lim := rate.NewLimiter(rate.Every(per/time.Duration(rps)), rps)
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if err := lim.Wait(req.Context()); err != nil {
				return nil, err
			}
			return next.RoundTrip(req)
		})
	}
}

func RetryMiddleware(max int, backoff time.Duration) RTMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			var lastErr error
			var bodyBytes []byte
			if req.Body != nil {
				bodyBytes, _ = io.ReadAll(req.Body)
				req.Body.Close()
			}
			for attempt := 0; attempt <= max; attempt++ {
				if bodyBytes != nil {
					req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				}
				resp, err := next.RoundTrip(req)
				// 5xx/429 — кандидаты на повтор, остальное — сразу выходим
				if err == nil && resp.StatusCode < 500 && resp.StatusCode != 429 {
					return resp, nil
				}
				if resp != nil {
					resp.Body.Close()
				}
				lastErr = err
				select {
				case <-req.Context().Done():
					return nil, req.Context().Err()
				case <-time.After(time.Duration(attempt+1) * backoff):
				}
			}
			if lastErr != nil {
				return nil, lastErr
			}
			return nil, fmt.Errorf("retry attempts exceeded")
		})
	}
}

// Утилита: собрать цепочку вокруг исходного транспорта
func chain(rt http.RoundTripper, mws ...RTMiddleware) http.RoundTripper {
	out := rt
	for i := len(mws) - 1; i >= 0; i-- {
		out = mws[i](out)
	}
	return out
}
