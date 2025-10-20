package lw_api

import (
	"errors"
	"github.com/MMC-BK/lw-api/auth"
	client2 "github.com/MMC-BK/lw-api/client"
	"github.com/MMC-BK/lw-api/orders"
	"github.com/MMC-BK/lw-api/processedorders"
	"net/url"
	"time"
)

type LinnworksAPI struct {
	//Branches of API
	*orders.Orders
	*processedorders.ProcessedOrders
}

type LinnworksAPIBuilder struct {
	baseURL      string
	appID        string
	appSecret    string
	token        string
	refreshToken string
	err          []error
}

func NewLinnworksAPIBuilder() *LinnworksAPIBuilder {
	return &LinnworksAPIBuilder{err: make([]error, 0)}
}

func (b *LinnworksAPIBuilder) BaseURL(base string) *LinnworksAPIBuilder {
	if b == nil {
		return nil
	}
	if base == "" {
		b.err = append(b.err, errors.New("baseURL is required"))
		return b
	}
	if _, err := url.Parse(base); err != nil {
		b.err = append(b.err, err)
		return b
	}
	b.baseURL = base
	return b
}

func (b *LinnworksAPIBuilder) Token(token string) *LinnworksAPIBuilder {
	if b == nil {
		return nil
	}
	if token == "" {
		b.err = append(b.err, errors.New("token is required"))
		return b
	}
	b.token = token
	return b
}

func (b *LinnworksAPIBuilder) AppID(appID string) *LinnworksAPIBuilder {
	if b == nil {
		return nil
	}
	if appID == "" {
		b.err = append(b.err, errors.New("applicationID is required"))
		return b
	}
	b.appID = appID
	return b
}

func (b *LinnworksAPIBuilder) AppSecret(appSecret string) *LinnworksAPIBuilder {
	if b == nil {
		return nil
	}
	if appSecret == "" {
		b.err = append(b.err, errors.New("applicationSecret is required"))
		return b
	}
	b.appSecret = appSecret
	return b
}

func (b *LinnworksAPIBuilder) RefreshToken(refresh string) *LinnworksAPIBuilder {
	if b == nil {
		return nil
	}
	b.refreshToken = refresh
	return b
}

func (b *LinnworksAPIBuilder) Build() (*LinnworksAPI, error) {
	if b == nil {
		return nil, errors.New("builder is nil")
	}

	b.requiredFields()
	if len(b.err) > 0 {
		return nil, errors.Join(b.err...)
	}

	creds := auth.Credentials{
		App: auth.App{
			ApplicationId:     b.appID,
			ApplicationSecret: b.appSecret,
			Token:             b.token,
		},
		RefreshToken: b.refreshToken,
	}

	ts := auth.NewRefreshTokenSource(creds)
	authAPI := auth.NewSimpleAuth(ts)

	c, err := client2.NewClient(
		client2.WithBaseURL(b.baseURL),
		client2.WithTransportChain(
			client2.RateLimitMiddleware(10, time.Second),
			client2.RetryMiddleware(2, 200*time.Millisecond),
			client2.AuthMiddleware(authAPI),
			client2.AutoRefreshOn401(ts),
		),
	)
	if err != nil {
		return nil, err
	}

	ordersAPI := orders.NewOrders(c)
	processedOrdersAPI := processedorders.NewProcessedOrders(c)

	return &LinnworksAPI{
		Orders:          ordersAPI,
		ProcessedOrders: processedOrdersAPI,
	}, nil
}

func (b *LinnworksAPIBuilder) requiredFields() {
	if b.baseURL == "" {
		b.err = append(b.err, errors.New("baseURL is required"))
	}
	if b.token == "" {
		b.err = append(b.err, errors.New("token is required"))
	}
	if b.appID == "" {
		b.err = append(b.err, errors.New("applicationID is required"))
	}
	if b.appSecret == "" {
		b.err = append(b.err, errors.New("applicationSecret is required"))
	}
}
