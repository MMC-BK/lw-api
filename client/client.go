package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type MakeRequest interface {
	DoJSON(ctx context.Context, method, path string, query url.Values, in any, out any) error
}

// ========== 5) Client и общий doJSON ==========
type Client struct {
	hc      *http.Client
	baseURL *url.URL
}

type Option func(*Client) error

func NewClient(opts ...Option) (*Client, error) {
	c := &Client{hc: &http.Client{Timeout: 30 * time.Second}}
	if c.hc.Transport == nil {
		c.hc.Transport = http.DefaultTransport
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func WithBaseURL(raw string) Option {
	return func(c *Client) error {
		u, err := url.Parse(raw)
		if err != nil {
			return err
		}
		c.baseURL = u
		return nil
	}
}

func WithTransportChain(mws ...RTMiddleware) Option {
	return func(c *Client) error {
		base := c.hc.Transport
		c.hc.Transport = chain(base, mws...)
		return nil
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error { c.hc = hc; return nil }
}

func (c *Client) DoJSON(ctx context.Context, method, path string, query url.Values, in any, out any) error {
	u := c.baseURL.ResolveReference(&url.URL{Path: path})
	if query != nil {
		u.RawQuery = query.Encode()
	}
	fmt.Println(u.String())
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// Тут можно поставить общий User-Agent
	// req.Header.Set("User-Agent", "linnworks-sdk-go/1.0")

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// единый разбор ошибок API
		var apiErr struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		if apiErr.Message == "" {
			apiErr.Message = resp.Status
		}
		return fmt.Errorf("api error: %d %s (%s)", resp.StatusCode, apiErr.Message, apiErr.Code)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e APIError) Error() string { return fmt.Sprintf("%d %s (%s)", e.Status, e.Message, e.Code) }
