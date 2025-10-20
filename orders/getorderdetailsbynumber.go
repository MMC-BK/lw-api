package orders

import (
	"context"
	"errors"
	lw_api "github.com/MMC-BK/lw-api/client"
	"github.com/MMC-BK/lw-api/orders/models"
	"net/http"
	"net/url"
	"strconv"
)

type GetOrderDetailsByNumOrderIdRequest struct {
	OrderID *int32 `json:"orderId,omitempty"`
}

type GetOrderDetailsByNumOrderIdRequestBuilder struct {
	ctx    context.Context
	client lw_api.MakeRequest
	data   *GetOrderDetailsByNumOrderIdRequest
	err    []error
}

func (o Orders) GetOrderDetailsByNumOrderId(ctx context.Context) *GetOrderDetailsByNumOrderIdRequestBuilder {
	return &GetOrderDetailsByNumOrderIdRequestBuilder{
		ctx:    ctx,
		client: o.c,
		data:   &GetOrderDetailsByNumOrderIdRequest{},
		err:    make([]error, 0),
	}
}

func (b *GetOrderDetailsByNumOrderIdRequestBuilder) OrderID(orderID int32) *GetOrderDetailsByNumOrderIdRequestBuilder {
	if b == nil {
		return nil
	}
	if orderID <= 0 {
		b.err = append(b.err, errors.New("orderId must be greater than 0"))
		return b
	}
	value := orderID
	b.data.OrderID = &value
	return b
}

func (b *GetOrderDetailsByNumOrderIdRequestBuilder) build() (*GetOrderDetailsByNumOrderIdRequest, error) {
	if b == nil {
		return nil, errors.New("builder is nil")
	}
	errs := make([]error, len(b.err))
	copy(errs, b.err)
	if b.data.OrderID == nil {
		errs = append(errs, errors.New("orderId is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return b.data, nil
}

func (b *GetOrderDetailsByNumOrderIdRequestBuilder) Do() (*models.OrderDetails, error) {
	if b == nil {
		return nil, errors.New("builder is nil")
	}
	req, err := b.build()
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Set("OrderId", strconv.FormatInt(int64(*req.OrderID), 10))
	path := "/api/Orders/GetOrderDetailsByNumOrderId"

	var out models.OrderDetails
	if err := b.client.DoJSON(b.ctx, http.MethodGet, path, query, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
