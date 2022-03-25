package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (s *service) Sync(ctx context.Context, appID string, bssSyncOffer model.BssSyncOfferRequest) {
	go s.sync(context.Background(), appID, bssSyncOffer)
}

func (s *service) sync(ctx context.Context, appID string, bssSyncOffer model.BssSyncOfferRequest) {
	offers, err := s.repository.All(ctx)
	if err != nil {
		msg := fmt.Sprintf("syncing getting all offers error: [%s]", err)
		s.logger.Error(msg)
		return
	}

	for _, bssOffer := range bssSyncOffer.SyncOffers {
		if bssOffer.PrimaryFlag == "1" {
			err = s.syncPrimaryOffer(ctx, bssOffer)
			if err != nil {
				msg := fmt.Sprintf("syncing offer [%s] [%s]", bssOffer.Name, err)
				s.logger.Error(msg)
			}

			continue
		}

		err = s.syncSupplementaryOffer(ctx, bssOffer, offers)
		if err != nil {
			msg := fmt.Sprintf("syncing offer [%s] [%s]", bssOffer.Name, err)
			s.logger.Error(msg)
		}
	}
}

func (s *service) syncPrimaryOffer(ctx context.Context, bssOffer model.BssOffer) error {
	if bssOffer.Status == model.SuspendBssStatus || bssOffer.Status == model.RetirementBssStatus {
		return s.repository.RemoveByExternalID(ctx, bssOffer.ID)
	}

	offer, err := s.repository.GetByExternalID(ctx, bssOffer.ID)
	if err != nil {
		return err
	}

	nOffer := s.mapBssOfferToOffer(bssOffer)
	if offer != nil {
		nOffer.ID = offer.ID
		nOffer.CreatedAt = offer.CreatedAt
		nOffer.UpdatedAt = offer.UpdatedAt
	}

	if bssOffer.Relationships != nil && len(bssOffer.Relationships.Attached) > 0 {
		nOffer.Childrens = make([]model.Offer, 0, len(bssOffer.Relationships.Attached))

		for i := range bssOffer.Relationships.Attached {
			now := time.Now().Unix()

			nOffer.Childrens = append(nOffer.Childrens, model.Offer{
				ExternalID: &bssOffer.Relationships.Attached[i].ID,
				CreatedAt:  model.CustomTimeStamp(now),
				UpdatedAt:  model.CustomTimeStamp(now),
				ID:         uuid.NewString(),
			})
		}
	}

	_, err = s.repository.Upsert(ctx, nOffer)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) syncSupplementaryOffer(ctx context.Context, bssOffer model.BssOffer, offers []model.Offer) error {
	for _, offer := range offers {
		if len(offer.Childrens) == 0 {
			continue
		}

		var index *int

		for i := range offer.Childrens {
			if offer.Childrens[i].ExternalID == nil || *offer.Childrens[i].ExternalID != bssOffer.ID {
				continue
			}

			if bssOffer.Status == model.SuspendBssStatus || bssOffer.Status == model.RetirementBssStatus {
				index = &i
				break
			}

			nOffer := s.mapBssOfferToOffer(bssOffer)
			nOffer.ID = offer.Childrens[i].ID
			nOffer.CreatedAt = offer.Childrens[i].CreatedAt

			now := time.Now().Unix()
			nOffer.UpdatedAt = model.CustomTimeStamp(now)

			offer.Childrens[i] = nOffer

			_, err := s.repository.Upsert(ctx, offer)
			if err != nil {
				return err
			}

			return nil
		}

		if index != nil {
			nChildren := make([]model.Offer, 0)
			for i := range offer.Childrens {
				if i == *index {
					continue
				}

				nChildren = append(nChildren, offer.Childrens[i])
			}

			offer.Childrens = nChildren

			_, err := s.repository.Upsert(ctx, offer)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func (s *service) mapBssOfferToOffer(bssOffer model.BssOffer) model.Offer {
	offer := model.Offer{
		ExternalID: &bssOffer.ID,
		Name:       bssOffer.Name,
		Code:       &bssOffer.Code,
		Category:   bssOffer.Category,
		ClientType: model.IndividualClienType,
		Paymode:    model.PostpaidPayMode,
		StandAlone: true,
		OneOfFee:   bssOffer.OneOfFee,
		MonthlyFee: bssOffer.MontlyFee,
	}

	if bssOffer.OnSale == model.NoSaleAlone {
		offer.StandAlone = false
	}

	if bssOffer.PayMode == model.AllBssPaymode {
		offer.Paymode = model.AllPayMode
	}

	if bssOffer.PayMode == model.PostpaidBssPaymode {
		offer.Paymode = model.PostpaidPayMode
	}

	if bssOffer.Type == model.CorporativeBssOfferClient {
		offer.ClientType = model.CorporativeClienType
	}

	if bssOffer.Attributes != nil && len((*bssOffer.Attributes).Attribute) > 0 {
		offer.Metadata = map[string]string{}

		for _, attributte := range (*bssOffer.Attributes).Attribute {
			offer.Metadata[attributte.Code] = attributte.Value
		}
	}

	if bssOffer.Description != "" {
		offer.Description = &bssOffer.Description
	}

	if bssOffer.EffectiveDate != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", *bssOffer.EffectiveDate)
		ts := model.CustomTimeStamp(t.Unix())

		offer.EffectiveDate = &ts
	}

	if bssOffer.ExpirationDate != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", *bssOffer.ExpirationDate)

		ts := model.CustomTimeStamp(t.Unix())

		offer.ExpirationDate = &ts
	}

	return offer
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
