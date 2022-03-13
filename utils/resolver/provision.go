package resolver

import (
	"github.com/srrmendez/private-api-order/model"
)

func SendToProvisionSystem(order model.Order) error {
	if order.Type == model.VPSOrder {
	}

	return nil
}
