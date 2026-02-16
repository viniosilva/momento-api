package listopts

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListParams struct {
	Pagination PaginationInput
	Sort       SortInput
}

func (p ListParams) ToFindOptions() *options.FindOptions {
	order := 1
	if p.Sort.Order == OrderTypeDesc {
		order = -1
	}

	return options.Find().
		SetSkip(int64((p.Pagination.Page - 1) * p.Pagination.PageSize)).
		SetLimit(int64(p.Pagination.PageSize)).
		SetSort(bson.D{{
			Key:   p.Sort.Field,
			Value: order,
		}})
}
