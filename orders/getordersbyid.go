package orders

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-openapi/strfmt"

	lw_api "github.com/bogudini/lw-api/client"
	"github.com/bogudini/lw-api/orders/models"
)

type GetOrdersByIdRequestBuilder struct {
	ctx    context.Context
	client lw_api.MakeRequest
	data   *models.OrdersGetOrdersByIDRequest
	err    []error
}

func (o Orders) GetOrdersById(ctx context.Context) *GetOrdersByIdRequestBuilder {
	return &GetOrdersByIdRequestBuilder{
		ctx:    ctx,
		client: o.c,
		data:   &models.OrdersGetOrdersByIDRequest{},
		err:    make([]error, 0),
	}
}

func (b *GetOrdersByIdRequestBuilder) PkOrderIds(ids []strfmt.UUID) *GetOrdersByIdRequestBuilder {
	if b == nil {
		return nil
	}
	b.ensureData()
	if len(ids) == 0 {
		b.err = append(b.err, errors.New("pkOrderIds must contain at least one value"))
		return b
	}
	b.data.PkOrderIds = append([]strfmt.UUID(nil), ids...)
	return b
}

func (b *GetOrdersByIdRequestBuilder) build() (*models.OrdersGetOrdersByIDRequest, error) {
	if b == nil {
		return nil, errors.New("builder is nil")
	}
	b.ensureData()
	errs := make([]error, len(b.err))
	copy(errs, b.err)
	if len(b.data.PkOrderIds) == 0 {
		errs = append(errs, errors.New("pkOrderIds must contain at least one value"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return b.data, nil
}

func (b *GetOrdersByIdRequestBuilder) Do() ([]models.OrderDetails, error) {
	if b == nil {
		return nil, errors.New("builder is nil")
	}
	req, err := b.build()
	if err != nil {
		return nil, err
	}
	var out []models.OrderDetails
	if err := b.client.DoJSON(b.ctx, http.MethodPost, "/api/Orders/GetOrdersById", nil, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (b *GetOrdersByIdRequestBuilder) ensureData() {
	if b.data == nil {
		b.data = &models.OrdersGetOrdersByIDRequest{}
	}
}
