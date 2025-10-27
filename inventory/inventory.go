package inventory

import (
	"context"
	lw_api "github.com/MMC-BK/lw-api/client"
)

type Inventory struct{ c lw_api.MakeRequest }

func NewInventory(c lw_api.MakeRequest) *Inventory {
	return &Inventory{c: c}
}

func (i Inventory) GetStockLocations(ctx context.Context) *GetStockLocationsRequestBuilder {
	return &GetStockLocationsRequestBuilder{
		ctx:    ctx,
		client: i.c,
	}
}
