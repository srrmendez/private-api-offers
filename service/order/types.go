package service

import (
	"context"

	"github.com/srrmendez/private-api-order/model"
	pkgRepository "github.com/srrmendez/private-api-order/repository"
	log "github.com/srrmendez/services-interface-tools/pkg/logger"
)

type Order interface {
	Create(ctx context.Context, order model.Order, appID string) (*model.Order, error)
	Search(ctx context.Context, appID string, status *model.OrderStatusType, category *model.CategoryType, orderType *model.OrderType) ([]model.Order, error)
	Get(ctx context.Context, id string) (*model.Order, error)
	CreateServiceOrder(ctx context.Context, order model.ServiceOrderRequest, appID string, transactionID string) (*model.ServiceOrderResponse, error)
	UpdateOrderStatus(ctx context.Context, request model.UpdateOrderRequest, id string, clientID string) (*model.Order, error)
}

type service struct {
	repository pkgRepository.OrderRepository
	logger     log.Log
}
