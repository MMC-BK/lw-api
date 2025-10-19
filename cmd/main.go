package main

import (
	"context"
	"encoding/json"
	"fmt"
	lw_api "github.com/bogudini/lw-api"
	"github.com/bogudini/lw-api/processedorders/models"
	"time"
)

func main() {

	linnworks, err := lw_api.NewLinnworksAPIBuilder().
		AppID("057abf7c-1684-485e-8cb6-00ab6533313a").
		AppSecret("8d8a9e2d-c6f3-475b-9cc3-ff2f720c43cc").
		Token("1bf62f655b558617aa199166746b2acc").
		BaseURL("https://us-ext.linnworks.net").
		Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	orderDetails, err := linnworks.Orders.GetOrderDetailsByNumOrderId(context.Background()).OrderID(334059).Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(orderDetails.CustomerInfo.Address.Address1)
	marshal, err := json.Marshal(orderDetails)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonString := string(marshal)
	fmt.Println(jsonString)

	toDate := time.Now()
	fromDate := time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC)

	searchResult, err := linnworks.Orders.SearchProcessedOrders(context.Background()).
		PageNumber(1).
		ResultsPerPage(50).
		SearchSorting(models.SearchSortingSortDirectionASC, "dReceivedDate").
		FromDate(fromDate).
		ToDate(toDate).
		DateField(models.SearchProcessedOrdersRequestDateFieldReceived).
		Do()
	if err != nil {
		fmt.Println(err)
		return
	}

	marshal, err = json.Marshal(searchResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonString = string(marshal)
	fmt.Println(jsonString)

}
