package pagination

import "math"

type PaginationInput struct {
	Page     int
	PageSize int
}

type PaginationOutput struct {
	TotalCount int64
	Page       int
	PageSize   int
	TotalPages int
}

func NewPaginationInput(page, pageSize int) PaginationInput {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return PaginationInput{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p PaginationInput) ToSkip() int64 {
	return int64((p.Page - 1) * p.PageSize)
}

type Paginated[T any] struct {
	Data       []T
	Pagination PaginationOutput
}

func NewPaginated[T any](data []T, totalCount int64, pagination PaginationInput) Paginated[T] {
	return Paginated[T]{
		Data: data,
		Pagination: PaginationOutput{
			TotalCount: totalCount,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: int(math.Ceil(float64(totalCount) / float64(pagination.PageSize))),
		},
	}
}
