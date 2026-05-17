package listopts

import "fmt"

type ListParams struct {
	Pagination PaginationInput
	Sort       SortInput
}

func (p ListParams) ToSQLOrder() string {
	order := "ASC"
	if p.Sort.Order == OrderTypeDesc {
		order = "DESC"
	}

	return fmt.Sprintf("%s %s", p.Sort.Field, order)
}
