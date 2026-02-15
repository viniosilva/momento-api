package dto

import "strings"

var (
	validSortFields = map[string]bool{
		"created_at": true,
		"updated_at": true,
	}

	validSortOrders = map[OrderType]bool{
		OrderTypeAsc:  true,
		OrderTypeDesc: true,
	}
)

type OrderType string

const (
	OrderTypeAsc  OrderType = "asc"
	OrderTypeDesc OrderType = "desc"
)

func OrderTypeFromString(orderType string) OrderType {
	if strings.ToLower(orderType) == string(OrderTypeAsc) {
		return OrderTypeAsc
	}

	return OrderTypeDesc
}

type SortInput struct {
	Field string
	Order OrderType
}

func NewSortInput(sortBy string, sortOrder OrderType) SortInput {
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	if !validSortOrders[sortOrder] {
		sortOrder = OrderTypeDesc
	}

	return SortInput{
		Field: sortBy,
		Order: sortOrder,
	}
}
