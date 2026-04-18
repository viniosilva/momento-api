package listopts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"momento/pkg/listopts"
)

func TestNewPaginationInput(t *testing.T) {
	t.Run("should use provided values when valid", func(t *testing.T) {
		got := listopts.NewPaginationInput(2, 20)

		assert.Equal(t, 2, got.Page)
		assert.Equal(t, 20, got.PageSize)
	})

	t.Run("should default page to 1 when page < 1", func(t *testing.T) {
		got := listopts.NewPaginationInput(0, 10)

		assert.Equal(t, 1, got.Page)
	})

	t.Run("should default pageSize to 10 when pageSize < 1", func(t *testing.T) {
		got := listopts.NewPaginationInput(1, 0)

		assert.Equal(t, 10, got.PageSize)
	})

	t.Run("should default pageSize to 10 when pageSize > 100", func(t *testing.T) {
		got := listopts.NewPaginationInput(1, 101)

		assert.Equal(t, 10, got.PageSize)
	})
}

func TestPaginationInput_ToSkip(t *testing.T) {
	t.Run("should return 0 for first page", func(t *testing.T) {
		p := listopts.PaginationInput{Page: 1, PageSize: 10}

		assert.Equal(t, int64(0), p.ToSkip())
	})

	t.Run("should return correct skip for page 2", func(t *testing.T) {
		p := listopts.PaginationInput{Page: 2, PageSize: 10}

		assert.Equal(t, int64(10), p.ToSkip())
	})

	t.Run("should return correct skip for page 3 with custom page size", func(t *testing.T) {
		p := listopts.PaginationInput{Page: 3, PageSize: 20}

		assert.Equal(t, int64(40), p.ToSkip())
	})
}

func TestNewPaginated(t *testing.T) {
	t.Run("should create paginated result with correct pagination", func(t *testing.T) {
		data := []string{"a", "b", "c"}
		input := listopts.PaginationInput{Page: 1, PageSize: 3}

		got := listopts.NewPaginated(data, 10, input)

		assert.Equal(t, data, got.Data)
		assert.Equal(t, int64(10), got.Pagination.TotalCount)
		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 3, got.Pagination.PageSize)
		assert.Equal(t, 4, got.Pagination.TotalPages)
	})

	t.Run("should calculate total pages correctly when exact fit", func(t *testing.T) {
		data := []string{"a", "b"}
		input := listopts.PaginationInput{Page: 1, PageSize: 2}

		got := listopts.NewPaginated(data, 4, input)

		assert.Equal(t, 2, got.Pagination.TotalPages)
	})

	t.Run("should round up total pages when remainder", func(t *testing.T) {
		data := []string{"a"}
		input := listopts.PaginationInput{Page: 1, PageSize: 3}

		got := listopts.NewPaginated(data, 7, input)

		assert.Equal(t, 3, got.Pagination.TotalPages)
	})
}

func TestPaginationApplicationToResponse(t *testing.T) {
	t.Run("should convert output to response", func(t *testing.T) {
		output := listopts.PaginationOutput{
			TotalCount: 100,
			Page:       2,
			PageSize:   10,
			TotalPages: 10,
		}

		got := listopts.PaginationApplicationToResponse(output)

		assert.Equal(t, int64(100), got.TotalCount)
		assert.Equal(t, 2, got.Page)
		assert.Equal(t, 10, got.PageSize)
		assert.Equal(t, 10, got.TotalPages)
	})
}

func TestOrderTypeFromString(t *testing.T) {
	t.Run("should return asc for 'asc'", func(t *testing.T) {
		assert.Equal(t, listopts.OrderTypeAsc, listopts.OrderTypeFromString("asc"))
	})

	t.Run("should return asc for 'ASC' (case insensitive)", func(t *testing.T) {
		assert.Equal(t, listopts.OrderTypeAsc, listopts.OrderTypeFromString("ASC"))
	})

	t.Run("should return desc for 'desc'", func(t *testing.T) {
		assert.Equal(t, listopts.OrderTypeDesc, listopts.OrderTypeFromString("desc"))
	})

	t.Run("should return desc for unknown value", func(t *testing.T) {
		assert.Equal(t, listopts.OrderTypeDesc, listopts.OrderTypeFromString("invalid"))
	})

	t.Run("should return desc for empty string", func(t *testing.T) {
		assert.Equal(t, listopts.OrderTypeDesc, listopts.OrderTypeFromString(""))
	})
}

func TestNewSortInput(t *testing.T) {
	t.Run("should use provided values when valid", func(t *testing.T) {
		got := listopts.NewSortInput("updated_at", listopts.OrderTypeAsc)

		assert.Equal(t, "updated_at", got.Field)
		assert.Equal(t, listopts.OrderTypeAsc, got.Order)
	})

	t.Run("should default field to created_at when invalid", func(t *testing.T) {
		got := listopts.NewSortInput("invalid_field", listopts.OrderTypeAsc)

		assert.Equal(t, "created_at", got.Field)
	})

	t.Run("should default order to desc when invalid", func(t *testing.T) {
		got := listopts.NewSortInput("created_at", "invalid_order")

		assert.Equal(t, listopts.OrderTypeDesc, got.Order)
	})

	t.Run("should accept created_at field", func(t *testing.T) {
		got := listopts.NewSortInput("created_at", listopts.OrderTypeDesc)

		assert.Equal(t, "created_at", got.Field)
		assert.Equal(t, listopts.OrderTypeDesc, got.Order)
	})
}

func TestListParams_ToFindOptions(t *testing.T) {
	t.Run("should generate find options for asc order", func(t *testing.T) {
		params := listopts.ListParams{
			Pagination: listopts.PaginationInput{Page: 2, PageSize: 10},
			Sort:       listopts.SortInput{Field: "created_at", Order: listopts.OrderTypeAsc},
		}

		findOpts := params.ToFindOptions()

		assert.Equal(t, int64(10), *findOpts.Skip)
		assert.Equal(t, int64(10), *findOpts.Limit)
	})

	t.Run("should generate find options for desc order", func(t *testing.T) {
		params := listopts.ListParams{
			Pagination: listopts.PaginationInput{Page: 1, PageSize: 5},
			Sort:       listopts.SortInput{Field: "updated_at", Order: listopts.OrderTypeDesc},
		}

		findOpts := params.ToFindOptions()

		assert.Equal(t, int64(0), *findOpts.Skip)
		assert.Equal(t, int64(5), *findOpts.Limit)
	})
}
