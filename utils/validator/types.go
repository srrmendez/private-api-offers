package validator

import "github.com/srrmendez/private-api-order/model"

type OrderRequestValidator struct {
	categories map[model.CategoryType][]model.OrderType
}
