package service

import (
	"context"

	"github.com/srrmendez/private-api-offers/conf"
	"github.com/srrmendez/private-api-offers/model"
	"github.com/srrmendez/private-api-offers/repository"
	log "github.com/srrmendez/services-interface-tools/pkg/logger"
)

type OfferService interface {
	Search(ctx context.Context, appID string, active *bool, category *model.CategoryType) ([]model.Offer, error)
	Sync(ctx context.Context, appID string, bssSyncOffer model.BssSyncOfferRequest) error
	Get(ctx context.Context, id string, appID string) (*model.Offer, error)
	GetSecondaryOffers(ctx context.Context, ids []string) ([]model.Offer, error)
}

type service struct {
	logger                  log.Log
	repository              repository.OfferRepository
	supplementaryRepository repository.OfferRepository
	confCategories          map[string]conf.Category
}
