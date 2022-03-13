package repository

import (
	"context"

	"github.com/srrmendez/private-api-order/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	All(ctx context.Context, appID string) ([]model.Order, error)
	Upsert(ctx context.Context, order model.Order) (*model.Order, error)
	Search(ctx context.Context, appID string, category *model.CategoryType, status *model.OrderStatusType, orderType *model.OrderType) ([]model.Order, error)
	Get(ctx context.Context, id string) (*model.Order, error)
	GetByExternalID(ctx context.Context, externalID string) (*model.Order, error)
}

type repository struct {
	collection *mongo.Collection
}
