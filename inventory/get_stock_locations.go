package inventory

import (
	"context"
	"errors"
	"net/http"

	lw_api "github.com/MMC-BK/lw-api/client"
	"github.com/MMC-BK/lw-api/inventory/models"
)

type GetStockLocationsRequestBuilder struct {
	ctx    context.Context
	client lw_api.MakeRequest
}

func (b *GetStockLocationsRequestBuilder) Do() ([]models.StockLocation, error) {
	if b == nil {
		return nil, errors.New("builder is nil")
	}

	var out []models.StockLocation
	if err := b.client.DoJSON(b.ctx, http.MethodGet, "/api/Inventory/GetStockLocations", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
