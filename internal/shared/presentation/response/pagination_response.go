package response

import "pinnado/internal/shared/application/dto"

type PaginatedResponse[T any] struct {
	Data       []T                `json:"data" example:"[item1, item2, item3]"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	TotalCount int64 `json:"total_count" example:"100"`
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalPages int   `json:"total_pages" example:"5"`
}

func PaginationApplicationToResponse(pagination dto.PaginationOutput) PaginationResponse {
	return PaginationResponse(pagination)
}
