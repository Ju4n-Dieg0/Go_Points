package pagination

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestNewParams(t *testing.T) {
	app := fiber.New()

	t.Run("default values", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		req := app.AcquireCtx(ctx)
		defer app.ReleaseCtx(req)

		params := NewParams(req)

		assert.Equal(t, DefaultPage, params.Page)
		assert.Equal(t, DefaultLimit, params.Limit)
		assert.Equal(t, "desc", params.Order)
		assert.Equal(t, "", params.Sort)
		assert.Equal(t, "", params.Search)
		assert.Empty(t, params.Filter)
	})

	t.Run("custom values", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		req := app.AcquireCtx(ctx)
		defer app.ReleaseCtx(req)

		// Simular query string: ?page=2&limit=20&sort=name&order=asc&search=test
		req.Request().URI().SetQueryString("page=2&limit=20&sort=name&order=asc&search=test")

		params := NewParams(req)

		assert.Equal(t, 2, params.Page)
		assert.Equal(t, 20, params.Limit)
		assert.Equal(t, "name", params.Sort)
		assert.Equal(t, "asc", params.Order)
		assert.Equal(t, "test", params.Search)
	})

	t.Run("with filters", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		req := app.AcquireCtx(ctx)
		defer app.ReleaseCtx(req)

		req.Request().URI().SetQueryString("filter[status]=active&filter[price__gte]=100")

		params := NewParams(req)

		assert.Equal(t, 2, len(params.Filter))
		assert.Equal(t, "active", params.Filter["status"])
		assert.Equal(t, "100", params.Filter["price__gte"])
	})

	t.Run("validates page minimum", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		req := app.AcquireCtx(ctx)
		defer app.ReleaseCtx(req)

		req.Request().URI().SetQueryString("page=0")

		params := NewParams(req)

		assert.Equal(t, DefaultPage, params.Page)
	})

	t.Run("validates limit maximum", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		req := app.AcquireCtx(ctx)
		defer app.ReleaseCtx(req)

		req.Request().URI().SetQueryString("limit=500")

		params := NewParams(req)

		assert.Equal(t, MaxLimit, params.Limit)
	})

	t.Run("normalizes order", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		req := app.AcquireCtx(ctx)
		defer app.ReleaseCtx(req)

		req.Request().URI().SetQueryString("order=invalid")

		params := NewParams(req)

		assert.Equal(t, "desc", params.Order)
	})
}

func TestParamsGetOffset(t *testing.T) {
	tests := []struct {
		name   string
		page   int
		limit  int
		offset int
	}{
		{"page 1", 1, 10, 0},
		{"page 2", 2, 10, 10},
		{"page 3", 3, 20, 40},
		{"page 5", 5, 25, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &Params{Page: tt.page, Limit: tt.limit}
			assert.Equal(t, tt.offset, params.GetOffset())
		})
	}
}

func TestParamsGetOrderBy(t *testing.T) {
	tests := []struct {
		name    string
		sort    string
		order   string
		orderBy string
	}{
		{"no sort", "", "desc", ""},
		{"sort asc", "name", "asc", "name ASC"},
		{"sort desc", "created_at", "desc", "created_at DESC"},
		{"sort with uppercase", "price", "ASC", "price ASC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &Params{Sort: tt.sort, Order: tt.order}
			assert.Equal(t, tt.orderBy, params.GetOrderBy())
		})
	}
}

func TestParamsCalculateMeta(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		limit      int
		total      int64
		totalPages int
	}{
		{"exact pages", 1, 10, 100, 10},
		{"partial page", 1, 10, 95, 10},
		{"single page", 1, 10, 5, 1},
		{"no results", 1, 10, 0, 1},
		{"large dataset", 5, 25, 1000, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &Params{Page: tt.page, Limit: tt.limit}
			meta := params.CalculateMeta(tt.total)

			assert.Equal(t, tt.page, meta.Page)
			assert.Equal(t, tt.limit, meta.Limit)
			assert.Equal(t, tt.total, meta.Total)
			assert.Equal(t, tt.totalPages, meta.TotalPages)
		})
	}
}

func TestIsValidFieldName(t *testing.T) {
	tests := []struct {
		name  string
		field string
		valid bool
	}{
		{"valid simple", "name", true},
		{"valid with underscore", "created_at", true},
		{"valid with numbers", "field123", true},
		{"invalid with dash", "created-at", false},
		{"invalid with dot", "user.name", false},
		{"invalid with space", "user name", false},
		{"invalid with special char", "name!", false},
		{"invalid with sql injection", "name; DROP TABLE", false},
		{"valid with operator", "price__gte", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, isValidFieldName(tt.field))
		})
	}
}

func TestNewResponse(t *testing.T) {
	type TestModel struct {
		ID   uint
		Name string
	}

	data := []TestModel{
		{ID: 1, Name: "Test"},
	}
	testMeta := Meta{
		Page:       1,
		Limit:      10,
		Total:      1,
		TotalPages: 1,
	}

	response := NewResponse(data, testMeta)

	assert.NotNil(t, response)
	assert.Equal(t, data, response.Data)
	assert.Equal(t, testMeta, response.Meta)
}
