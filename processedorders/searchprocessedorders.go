package processedorders

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"

	lw_api "github.com/MMC-BK/lw-api/client"
	processedmodels "github.com/MMC-BK/lw-api/processedorders/models"
)

type ProcessedOrders struct{ c lw_api.MakeRequest }

func NewProcessedOrders(c lw_api.MakeRequest) *ProcessedOrders {
	return &ProcessedOrders{c: c}
}

type SearchProcessedOrdersRequestBuilder struct {
	ctx     context.Context
	client  lw_api.MakeRequest
	payload *processedmodels.ProcessedOrdersSearchProcessedOrdersRequest
	request *processedmodels.SearchProcessedOrdersRequest
	err     []error
}

func (o ProcessedOrders) SearchProcessedOrders(ctx context.Context) *SearchProcessedOrdersRequestBuilder {
	req := &processedmodels.SearchProcessedOrdersRequest{}
	return &SearchProcessedOrdersRequestBuilder{
		ctx:    ctx,
		client: o.c,
		payload: &processedmodels.ProcessedOrdersSearchProcessedOrdersRequest{
			Request: req,
		},
		request: req,
		err:     make([]error, 0),
	}
}

func (b *SearchProcessedOrdersRequestBuilder) DateField(value string) *SearchProcessedOrdersRequestBuilder {

	b.request.DateField = value
	return b
}

func (b *SearchProcessedOrdersRequestBuilder) FromDate(value time.Time) *SearchProcessedOrdersRequestBuilder {

	b.request.FromDate = strfmt.DateTime(value)
	return b
}

func (b *SearchProcessedOrdersRequestBuilder) ToDate(value time.Time) *SearchProcessedOrdersRequestBuilder {

	b.request.ToDate = strfmt.DateTime(value)
	return b
}

func (b *SearchProcessedOrdersRequestBuilder) PageNumber(page int32) *SearchProcessedOrdersRequestBuilder {
	if page <= 0 {
		b.err = append(b.err, errors.New("page number must be greater than 0"))
	}
	b.request.PageNumber = page
	return b
}

func (b *SearchProcessedOrdersRequestBuilder) ResultsPerPage(value int32) *SearchProcessedOrdersRequestBuilder {

	b.request.ResultsPerPage = value
	return b
}

func (b *SearchProcessedOrdersRequestBuilder) SearchTerm(value string) *SearchProcessedOrdersRequestBuilder {

	b.request.SearchTerm = value
	return b
}

func (b *SearchProcessedOrdersRequestBuilder) SearchFilters(filters ...*processedmodels.SearchFilters) *SearchProcessedOrdersRequestBuilder {
	b.request.SearchFilters = filters
	return b
}

func (b *SearchProcessedOrdersRequestBuilder) AddSearchFilter(filter *processedmodels.SearchFilters) *SearchProcessedOrdersRequestBuilder {
	if filter == nil {
		b.err = append(b.err, errors.New("search filter cannot be nil"))
		return b
	}
	b.request.SearchFilters = append(b.request.SearchFilters, filter)
	return b
}
func (b *SearchProcessedOrdersRequestBuilder) SearchSortingASC(sortField string) *SearchProcessedOrdersRequestBuilder {
	if b.request.SearchSorting != nil {
		b.err = append(b.err, errors.New("search sorting might be used once"))
	}
	return b.searchSorting(processedmodels.SearchSortingSortDirectionASC, sortField)
}

func (b *SearchProcessedOrdersRequestBuilder) SearchSortingDESC(sortField string) *SearchProcessedOrdersRequestBuilder {
	if b.request.SearchSorting != nil {
		b.err = append(b.err, errors.New("search sorting might be used once"))
	}
	return b.searchSorting(processedmodels.SearchSortingSortDirectionDESC, sortField)
}

func (b *SearchProcessedOrdersRequestBuilder) searchSorting(sortDirection, sortField string) *SearchProcessedOrdersRequestBuilder {
	sorting := processedmodels.SearchSorting{
		SortDirection: sortDirection,
		SortField:     sortField,
	}
	formats := strfmt.NewFormats()
	err := sorting.Validate(formats)
	if err != nil {
		b.err = append(b.err, err)
		return nil
	}
	b.request.SearchSorting = &sorting
	return b
}

func (b *SearchProcessedOrdersRequestBuilder) build() (*processedmodels.ProcessedOrdersSearchProcessedOrdersRequest, error) {
	if b == nil {
		return nil, errors.New("builder is nil")
	}
	errs := make([]error, len(b.err))
	copy(errs, b.err)
	if b.payload == nil {
		errs = append(errs, errors.New("payload is required"))
	} else if b.payload.Request == nil {
		errs = append(errs, errors.New("request body is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return b.payload, nil
}

func (b *SearchProcessedOrdersRequestBuilder) Do() (*processedmodels.GenericPagedResultProcessedOrderWeb, error) {
	if b == nil {
		return nil, errors.New("builder is nil")
	}
	req, err := b.build()
	if err != nil {
		return nil, err
	}
	var out processedmodels.SearchProcessedOrdersResponse
	if err := b.client.DoJSON(b.ctx, http.MethodPost, "/api/ProcessedOrders/SearchProcessedOrders", nil, req, &out); err != nil {
		return nil, err
	}
	return out.ProcessedOrders, nil
}
