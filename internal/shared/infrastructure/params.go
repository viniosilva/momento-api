package infrastructure

import (
	shareddto "pinnado/internal/shared/application/dto"
	"pinnado/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListParams struct {
	Pagination pagination.PaginationInput
	Sort       shareddto.SortInput
}

func (p ListParams) ToFindOptions() *options.FindOptions {
	order := 1
	if p.Sort.Order == shareddto.OrderTypeDesc {
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
