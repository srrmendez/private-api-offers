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

func (s *service) Sync(ctx context.Context, appID string, bssSyncOffer model.BssSyncOfferRequest) error {
	return s.sync(context.Background(), appID, bssSyncOffer)
}

func (s *service) sync(ctx context.Context, appID string, bssSyncOffer model.BssSyncOfferRequest) error {
	for i := range bssSyncOffer.SyncOffers {
		if bssSyncOffer.SyncOffers[i].Offer.PrimaryFlag == "1" {
			if err := s.syncPrimaryOffer(ctx, bssSyncOffer.SyncOffers[i].Offer); err != nil {
				msg := fmt.Sprintf("syncing offer [%s] [%s]", bssSyncOffer.SyncOffers[i].Offer.Name, err)
				s.logger.Error(msg)

				return err
			}

			continue
		}

		if err := s.syncSupplementaryOffer(ctx, bssSyncOffer.SyncOffers[i].Offer); err != nil {
			msg := fmt.Sprintf("syncing offer [%s] [%s]", bssSyncOffer.SyncOffers[i].Offer.Name, err)
			s.logger.Error(msg)

			return err
		}
	}

	return nil
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
			supOffer, err := s.supplementaryRepository.GetByExternalID(ctx, bssOffer.Relationships.Attached[i].ID)
			if err != nil {
				return nil
			}

			if supOffer != nil {
				nOffer.Supplementaries = append(nOffer.Supplementaries, supOffer.ID)
				continue
			}

			sOffer, err := s.supplementaryRepository.Upsert(ctx, model.Offer{
				ExternalID: &bssOffer.Relationships.Attached[i].ID,
			})
			if err != nil {
				return err
			}

			nOffer.Supplementaries = append(nOffer.Supplementaries, sOffer.ID)
		}
	}

	if _, err = s.repository.Upsert(ctx, *nOffer); err != nil {
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
		ActivationFate:  bssOffer.OneOfFee,
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
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				d, err := strconv.Atoi(attributte.Value)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.AliasQty = &d
			case "C_BD_NUM":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.Database = s.checkDatabaseNil(offer.DataCenterResourceAttributtes.Database)

				d, err := strconv.Atoi(attributte.Value)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.Database.Quantity = d
			case "CN_BD_SPACE":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.Database = s.checkDatabaseNil(offer.DataCenterResourceAttributtes.Database)

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.Database.Amount = d
			case "C_BD_SPACE_UNIT":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.Database = s.checkDatabaseNil(offer.DataCenterResourceAttributtes.Database)

				offer.DataCenterResourceAttributtes.Database.Unit = attributte.Value
			case "CN_CPU_NUM":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				d, err := strconv.Atoi(attributte.Value)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.CPUQty = &d
			case "CN_FTP_NUM":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				d, err := strconv.Atoi(attributte.Value)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.FTPQty = &d
			case "CN_PORT_NUM":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				d, err := strconv.Atoi(attributte.Value)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.NetworkInterfaceQty = &d
			case "CN_IP_NUM":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.PublicIPAddress = &attributte.Value
			case "CN_RAM_SPACE":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.RAM = s.checkRAMNil(offer.DataCenterResourceAttributtes.RAM)

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.RAM.Amount = d
			case "C_RAM_SPACE_UNIT":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.RAM = s.checkRAMNil(offer.DataCenterResourceAttributtes.RAM)

				offer.DataCenterResourceAttributtes.RAM.Unit = attributte.Value
			case "C_DISK_SPACE":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.HDD = s.checkHDDNil(offer.DataCenterResourceAttributtes.HDD)

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.HDD.Amount = d
			case "C_DISK_SPACE_UNIT":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.HDD = s.checkHDDNil(offer.DataCenterResourceAttributtes.HDD)

				offer.DataCenterResourceAttributtes.HDD.Unit = attributte.Value

			case "C_RATE_NUM":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.Bandwidth = s.checkBandwithNil(offer.DataCenterResourceAttributtes.Bandwidth)

				d, err := strconv.ParseFloat(attributte.Value, 64)
				if err != nil {
					return nil, err
				}

				offer.DataCenterResourceAttributtes.Bandwidth.Amount = d

			case "C_RATE_UNIT":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.Bandwidth = s.checkBandwithNil(offer.DataCenterResourceAttributtes.Bandwidth)

				offer.DataCenterResourceAttributtes.Bandwidth.Unit = attributte.Value

			case "CN_VPN_LANIP":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.VPN = s.checkVPNNil(offer.DataCenterResourceAttributtes.VPN)

				offer.DataCenterResourceAttributtes.VPN.IPAddress = attributte.Value

			case "CN_VPN_NAME":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.VPN = s.checkVPNNil(offer.DataCenterResourceAttributtes.VPN)

				offer.DataCenterResourceAttributtes.VPN.Name = attributte.Value

			case "CN_DNS":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.DNS = s.checkDNSNil(offer.DataCenterResourceAttributtes.DNS)

				offer.DataCenterResourceAttributtes.DNS.DNS = attributte.Value

			case "CN_DNS_CNAME":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.DNS = s.checkDNSNil(offer.DataCenterResourceAttributtes.DNS)

				offer.DataCenterResourceAttributtes.DNS.Name = attributte.Value

			case "C_DATAC_ACCESS_TYPE":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.AccessType = &attributte.Value

			case "CN_VPS_LANIP":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.LANIPAddress = &attributte.Value

			case "CN_VPS_WANIP":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				offer.DataCenterResourceAttributtes.WANIPAddress = &attributte.Value

			case "C_SAVEVM_FALG":
				offer.DataCenterResourceAttributtes = s.checkDataCenterAttributesNil(offer.DataCenterResourceAttributtes)

				value := false

				if attributte.Value == "1" {
					value = true
				}

				offer.DataCenterResourceAttributtes.SaveVM = &value
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

func (s *service) checkDataCenterAttributesNil(v *model.DataCenterResourceAttributtes) *model.DataCenterResourceAttributtes {
	if v == nil {
		return &model.DataCenterResourceAttributtes{}
	}

	return v
}

func (s *service) checkDatabaseNil(v *model.Database) *model.Database {
	if v == nil {
		return &model.Database{}
	}

	return v
}

func (s *service) checkVPNNil(v *model.VPN) *model.VPN {
	if v == nil {
		return &model.VPN{}
	}

	return v
}

func (s *service) checkDNSNil(v *model.DNS) *model.DNS {
	if v == nil {
		return &model.DNS{}
	}

	return v
}

func (s *service) checkRAMNil(v *model.RAM) *model.RAM {
	if v == nil {
		return &model.RAM{}
	}

	return v
}

func (s *service) checkHDDNil(v *model.HDD) *model.HDD {
	if v == nil {
		return &model.HDD{}
	}

	return v
}

func (s *service) checkBandwithNil(v *model.BandWith) *model.BandWith {
	if v == nil {
		return &model.BandWith{}
	}

	return v
}
