package orders

import (
	lw_api "github.com/MMC-BK/lw-api/client"
)

type Orders struct{ c lw_api.MakeRequest }

func NewOrders(c lw_api.MakeRequest) *Orders {
	return &Orders{c: c}
}
