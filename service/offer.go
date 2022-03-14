package service

import (
	"context"
	"fmt"

	"github.com/srrmendez/private-api-offers/model"
	"github.com/srrmendez/private-api-offers/repository"
	log "github.com/srrmendez/services-interface-tools/pkg/logger"
)

func NewService(repository repository.OfferRepository, logger log.Log,
) *service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) Search(ctx context.Context, appID string, active *bool) ([]model.Offer, error) {
	if active == nil {
		offers, err := s.repository.All(ctx)
		if err != nil {
			msg := fmt.Sprintf("[%s] searching offers error [%s]", appID, err)

			s.logger.Error(msg)

			return nil, err
		}

		return offers, nil
	}

	offers, err := s.repository.Search(ctx, *active)
	if err != nil {
		msg := fmt.Sprintf("[%s] searching offers error [%s]", appID, err)

		s.logger.Error(msg)

		return nil, err
	}

	return offers, nil
}

func (s *service) Create(ctx context.Context, appID string, bssSyncOffer model.BssSyncOfferRequest) error {
	/*allOffers, err := s.repository.All(ctx)
	if err != nil {
		msg := fmt.Sprintf("[%s] searching offers error [%s]", appID, err)

		s.logger.Error(msg)

		return err
	}*/

	return nil
}

func (s *service) Get(ctx context.Context, id string, appID string) (*model.Offer, error) {
	offer, err := s.repository.Get(ctx, id)
	if err != nil {
		msg := fmt.Sprintf("[%s] getting offer [%s] error [%s]", appID, id, err)

		s.logger.Error(msg)

		return nil, err
	}

	return offer, nil
}
