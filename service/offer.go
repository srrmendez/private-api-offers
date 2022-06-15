package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/srrmendez/private-api-offers/conf"
	"github.com/srrmendez/private-api-offers/model"
	"github.com/srrmendez/private-api-offers/repository"
	log "github.com/srrmendez/services-interface-tools/pkg/logger"
)

func NewService(repository repository.OfferRepository, supplementary repository.OfferRepository, logger log.Log,
	confCategories map[string]conf.Category,
) *service {
	return &service{
		repository:              repository,
		supplementaryRepository: supplementary,
		logger:                  logger,
		confCategories:          confCategories,
	}
}

func (s *service) Search(ctx context.Context, appID string, active *bool, category *model.CategoryType) ([]model.Offer, error) {
	if active == nil && category == nil {
		offers, err := s.repository.All(ctx)
		if err != nil {
			msg := fmt.Sprintf("[%s] searching offers error [%s]", appID, err)

			s.logger.Error(msg)

			return nil, err
		}

		return offers, nil
	}

	offers, err := s.repository.Search(ctx, active, category)
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
	for i := range bssSyncOffer.SyncOffers {
		if bssSyncOffer.SyncOffers[i].Offer.PrimaryFlag == "1" {
			if err := s.syncPrimaryOffer(ctx, bssSyncOffer.SyncOffers[i].Offer); err != nil {
				msg := fmt.Sprintf("syncing offer [%s] [%s]", bssSyncOffer.SyncOffers[i].Offer.Name, err)
				s.logger.Error(msg)
			}

			continue
		}

		if err := s.syncSupplementaryOffer(ctx, bssSyncOffer.SyncOffers[i].Offer); err != nil {
			msg := fmt.Sprintf("syncing offer [%s] [%s]", bssSyncOffer.SyncOffers[i].Offer.Name, err)
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

	nOffer, err := s.mapBssOfferToOffer(bssOffer)
	if err != nil {
		return err
	}

	if offer != nil {
		nOffer.ID = offer.ID
		nOffer.CreatedAt = offer.CreatedAt
		nOffer.UpdatedAt = offer.UpdatedAt
	}

	if bssOffer.Relationships != nil && len(bssOffer.Relationships.Attached) > 0 {
		nOffer.Supplementaries = make([]string, 0, len(bssOffer.Relationships.Attached))

		for i := range bssOffer.Relationships.Attached {
			sOffer, err := s.supplementaryRepository.Upsert(ctx, model.Offer{
				ExternalID: &bssOffer.Relationships.Attached[i].ID,
			})
			if err != nil {
				return err
			}

			nOffer.Supplementaries = append(nOffer.Supplementaries, sOffer.ID)
		}
	}

	_, err = s.repository.Upsert(ctx, *nOffer)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) syncSupplementaryOffer(ctx context.Context, bssOffer model.BssOffer) error {
	if bssOffer.Status == model.SuspendBssStatus || bssOffer.Status == model.RetirementBssStatus {
		return s.supplementaryRepository.RemoveByExternalID(ctx, bssOffer.ID)
	}

	offer, err := s.repository.GetByExternalID(ctx, bssOffer.ID)
	if err != nil {
		return err
	}

	nOffer, err := s.mapBssOfferToOffer(bssOffer)
	if err != nil {
		return err
	}

	if offer != nil {
		nOffer.ID = offer.ID
		nOffer.CreatedAt = offer.CreatedAt
		nOffer.UpdatedAt = offer.UpdatedAt
	}

	if _, err = s.supplementaryRepository.Upsert(ctx, *nOffer); err != nil {
		return err
	}

	return nil
}

func (s *service) mapBssOfferToOffer(bssOffer model.BssOffer) (*model.Offer, error) {
	offer := model.Offer{
		ExternalID:      &bssOffer.ID,
		Name:            bssOffer.Name,
		ClientType:      model.IndividualClienType,
		Paymentmode:     model.PostpaidPayMode,
		Fare:            bssOffer.MontlyFee,
		Supplementaries: []string{},
	}

	if bssOffer.PayMode == model.AllBssPaymode {
		offer.Paymentmode = model.AllPayMode
	}

	if bssOffer.PayMode == model.PostpaidBssPaymode {
		offer.Paymentmode = model.PostpaidPayMode
	}

	if bssOffer.Type == model.CorporativeBssOfferClient {
		offer.ClientType = model.CorporativeClienType
	}

	if bssOffer.Attributes != nil && len((*bssOffer.Attributes).Attribute) > 0 {
		for _, attributte := range (*bssOffer.Attributes).Attribute {
			if v, ok := s.confCategories[attributte.Value]; ok {
				offer.Category = v.Category
				offer.Type = v.Type
				continue
			}

			switch attributte.Code {
			case "CN_ALIAS_NUM":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.Alias.Amount = &d
			case "C_BD_NUM":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.Atoi(attributte.Value)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.Database.Quantity = &d
			case "CN_BD_SPACE":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.Database.Amount = &d
			case "C_BD_SPACE_UNIT":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				offer.DataCenterResourceAttributtes.Database.Unit = &attributte.Value
			case "CN_CPU_NUM":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.CPU.Amount = &d
			case "CN_FTP_NUM":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.FTP.Amount = &d
			case "CN_PORT_NUM":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.NetworkInterface.Amount = &d
			case "CN_IP_NUM":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.IPAddress.Amount = &d
			case "CN_RAM_SPACE":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.Ram.Amount = &d
			case "C_RAM_SPACE_UNIT":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				offer.DataCenterResourceAttributtes.Ram.Unit = &attributte.Value
			case "C_DISK_SPACE":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.HDD.Amount = &d
			case "C_DISK_SPACE_UNIT":
				if offer.DataCenterResourceAttributtes == nil {
					offer.DataCenterResourceAttributtes = &model.DataCenterResourceAttributtes{}
				}

				offer.DataCenterResourceAttributtes.HDD.Unit = &attributte.Value
			}
		}
	}

	if bssOffer.EffectiveDate != nil {
		offer.EffectiveDate = *bssOffer.EffectiveDate
	}

	if bssOffer.ExpirationDate != nil {
		offer.ExpirationDate = *bssOffer.ExpirationDate
	}

	return &offer, nil
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

func (s *service) GetSecondaryOffers(ctx context.Context, ids []string) ([]model.Offer, error) {
	return s.supplementaryRepository.GetByIDList(ctx, ids)
}
