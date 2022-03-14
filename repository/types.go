package repository

import (
	"context"

	"github.com/srrmendez/private-api-offers/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type OfferRepository interface {
	All(ctx context.Context) ([]model.Offer, error)
	Upsert(ctx context.Context, offer model.Offer) (*model.Offer, error)
	BatchUpsert(ctx context.Context, offers []model.Offer) error
	Get(ctx context.Context, id string) (*model.Offer, error)
	Search(ctx context.Context, active bool) ([]model.Offer, error)
	Remove(ctx context.Context, ids []string) error
}

type repository struct {
	collection *mongo.Collection
}
