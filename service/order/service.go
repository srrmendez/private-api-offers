package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/srrmendez/private-api-order/model"
	pkgRepository "github.com/srrmendez/private-api-order/repository"
	"github.com/srrmendez/private-api-order/utils/resolver"
	log "github.com/srrmendez/services-interface-tools/pkg/logger"
)

func NewService(repo pkgRepository.OrderRepository,
	logger log.Log) *service {

	return &service{
		repository: repo,
		logger:     logger,
	}
}

func (s *service) Create(ctx context.Context, order model.Order, appID string) (*model.Order, error) {
	order.Status = model.PendingOrderStatus
	order.AppID = appID
	order.ExternalID = nil
	order.Source = model.WebPortalSource

	now := time.Now().Add(1 * time.Millisecond).Unix()

	order.Events = []model.Event{
		{
			Type:        model.CommercialValidationSendEvent,
			Date:        model.CustomTimeStamp(now),
			Description: "send order to commercial validation",
		},
	}

	ord, err := s.repository.Upsert(ctx, order)
	if err != nil {
		return nil, err
	}

	// BSS Process order
	go func() {

	}()

	return ord, nil
}

func (s *service) Search(ctx context.Context, appID string, status *model.OrderStatusType, category *model.CategoryType, orderType *model.OrderType) ([]model.Order, error) {
	if status == nil && category == nil && orderType == nil {
		return s.repository.All(ctx, appID)
	}

	return s.repository.Search(ctx, appID, category, status, orderType)
}

func (s *service) Get(ctx context.Context, id string) (*model.Order, error) {
	return s.repository.Get(ctx, id)
}

func (s *service) CreateServiceOrder(ctx context.Context, serviceOrder model.ServiceOrderRequest, appID string,
	transactionID string) (*model.ServiceOrderResponse, error) {
	order, err := s.repository.GetByExternalID(ctx, serviceOrder.OrderInfo.OrderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		order = s.mapServiceOrderToOrder(serviceOrder, appID)
	}

	now := time.Now().Add(1 * time.Millisecond).Unix()

	events := []model.Event{
		{
			Type:        model.CommercialValidationCompletedEvent,
			Date:        model.CustomTimeStamp(now),
			Description: "order validated successfully",
		},
		{
			Type:        model.ProvisionSendEvent,
			Date:        model.CustomTimeStamp(now),
			Description: "send to provision system",
		},
	}

	order.Events = append(order.Events, events...)

	order, err = s.repository.Upsert(ctx, *order)
	if err != nil {
		return nil, err
	}

	go s.provisionOrder(*order)

	return &model.ServiceOrderResponse{
		BSSTransactionID: transactionID,
		OrderID:          order.ID,
	}, nil
}

func (s *service) UpdateOrderStatus(ctx context.Context, request model.UpdateOrderRequest, id string, appID string) (*model.Order, error) {
	order, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, nil
	}

	order.Status = request.Status

	event := model.Event{
		Type:        model.ProvisionCompletedEvent,
		Description: request.Description,
	}

	if order.Status == model.DeclinedOrderStatus {
		event.Type = model.ProvisionDeclinedEvent
	}

	order.Events = append(order.Events, event)

	return s.repository.Upsert(ctx, *order)
}

func (s *service) mapServiceOrderToOrder(o model.ServiceOrderRequest, appID string) *model.Order {
	order := model.Order{
		AppID:        appID,
		ExternalID:   &o.OrderInfo.OrderID,
		Source:       model.CommercialSystemSource,
		Status:       model.PendingOrderStatus,
		Events:       []model.Event{},
		Category:     model.DataCenterCategory,
		CustomerType: model.IndividualCustomer,
	}

	if o.SubscriberInfo.NetworkType == model.YellowPagesWebNetwork {
		order.Category = model.YellowPagesCategory
	}

	if o.SubscriberInfo.NetworkType == model.YellowPagesPrintNetwork {
		order.Category = model.YellowPagesPrintCategory
	}

	if o.OrderInfo.CustomerType == model.ServiceOrderCorporateCustomerType {
		order.CustomerType = model.CorporativeCustomer
	}

	switch o.SubscriberInfo.ProvisionServiceType {
	case model.DedicatedServerServiceOrderProvision:
		order.Type = model.DedicatedServerOrder
		order.DedicatedServer = s.mapServiceOrderDedicatedServerToDedicatedServer(o.SubscriberInfo.DedicatedServer)
	case model.HouseLeasingServiceOrderProvision:
		order.Type = model.HouseLeasingOrder
		order.HouseLeasing = s.mapServiceOrderHouseLeasingToHouseLeasing(o.SubscriberInfo.HousingLeasing)
	case model.VirtualDataCenterServiceOrderProvision:
		order.Type = model.VirtualDataCenterOrder
		order.VirtualDataCenter = s.mapServiceOrderVirtualDataCenterToVirtualDC(o.SubscriberInfo.VirtualDataCenter)
	case model.WebHostingServiceOrderProvision:
		order.Type = model.WebHostingOrder
		order.WebHosting = s.mapServiceOrderWHToWH(o.SubscriberInfo.WebHosting)
	default:
		order.Type = model.VPSOrder
		order.VPS = s.mapServiceOrderVPSToVPS(o.SubscriberInfo.VPS)
	}

	return &order
}

func (s *service) mapResponsibleToContacts(responsibles []model.ServiceOrderResponsible) []model.Contact {
	contacts := make([]model.Contact, 0, len(responsibles))

	for _, responsible := range responsibles {
		contact := model.Contact{
			Email: responsible.Email,
			Phone: responsible.Phone,
			Type:  model.TechnicianContact,
		}

		if responsible.ResponsibleType == model.ModeratorServiceOrderResponsible {
			contact.Type = model.ModeratorContact
		}

		if responsible.ResponsibleType == model.NotificationServiceOrderResponsible {
			contact.Type = model.NotificationContact
		}

		splittedFullName := strings.Split(responsible.FullName, " ")

		if len(splittedFullName) == 0 {
			contact.Firstname = responsible.FullName
			contacts = append(contacts, contact)

			continue
		}

		for i := range splittedFullName {
			if i == 0 {
				contact.Firstname = splittedFullName[i]
				continue
			}

			if contact.Lastname == "" {
				contact.Lastname = splittedFullName[i]
				continue
			}

			contact.Lastname = fmt.Sprintf("%s %s", contact.Lastname, splittedFullName[i])
		}

		contacts = append(contacts, contact)
	}

	return contacts
}

func (s *service) mapServiceOrderAccessInfoToAccess(accessInfo model.ServiceOrderAccessInfo) model.Access {
	access := model.Access{
		Account:       accessInfo.AccessAcount,
		IPAddress:     accessInfo.IPAddress,
		ServiceNumber: accessInfo.ServiceNumber,
	}

	if accessInfo.AccessType != model.InternalServiceOrderAccess {
		link := model.InternalAccess
		access.Access = &link
	}

	if accessInfo.AccessType != model.ExternalServiceOrderAccess {
		link := model.ExternalAccess
		access.Access = &link
	}

	return access
}

func (s *service) mapServiceOrderDedicatedServerToDedicatedServer(d *model.ServiceOrderDedicatedServer) *model.DedicatedServer {
	if d == nil {
		return nil
	}

	return &model.DedicatedServer{
		ServerName: d.ServerName,
		SLA:        d.SLAFlag,
		Contacts:   s.mapResponsibleToContacts(d.Responsible),
		Access:     s.mapServiceOrderAccessInfoToAccess(d.AccessInfo),
	}
}

func (s *service) mapServiceOrderHouseLeasingToHouseLeasing(d *model.ServiceOrderHousingLeasing) *model.HouseLeasing {
	if d == nil {
		return nil
	}

	return &model.HouseLeasing{
		ServerName: d.ServerName,
		SLA:        d.SLAFlag,
		Contacts:   s.mapResponsibleToContacts(d.Responsible),
		Access:     s.mapServiceOrderAccessInfoToAccess(d.AccessInfo),
	}
}

func (s *service) mapServiceOrderVirtualDataCenterToVirtualDC(d *model.ServiceOrderVirtualDataCenter) *model.VirtualDataCenter {
	if d == nil {
		return nil
	}

	return &model.VirtualDataCenter{
		ServerName: d.ServerName,
		SLA:        d.SLAFlag,
		Contacts:   s.mapResponsibleToContacts(d.Responsible),
		Access:     s.mapServiceOrderAccessInfoToAccess(d.AccessInfo),
	}
}

func (s *service) mapServiceOrderVPSToVPS(d *model.ServiceOrderVPS) *model.VPS {
	if d == nil {
		return nil
	}

	return &model.VPS{
		ServerName: d.ServerName,
		Contacts:   s.mapResponsibleToContacts(d.Responsible),
		Access:     s.mapServiceOrderAccessInfoToAccess(d.AccessInfo),
		Forum:      d.ForumFlag,
		SLA:        d.SLAFlag,
	}
}

func (s *service) mapServiceOrderWHToWH(d *model.ServiceOrderWebHosting) *model.WebHosting {
	if d == nil {
		return nil
	}

	wh := model.WebHosting{
		Contacts:            s.mapResponsibleToContacts(d.Responsible),
		Access:              []model.Access{},
		Domain:              d.MainDomain,
		Database:            d.DataBase,
		ProgrammingLanguage: d.ProgrammingLanguage,
		Advertising:         false,
		WebsiteType:         d.WebSiteType,
		WebsiteGoal:         d.WebsiteGoal,
		Framework:           d.FrameWorkName,
		CMS:                 d.CMSName,
		LMS:                 d.LMSVersion,
		CMF:                 d.CMFVersion,
		Topic:               d.TopicToDiscuss,
		Remark:              d.Remark,
		Alias:               []string{},
		FTPAccounts:         []string{},
	}

	for _, accessInfo := range d.AccessInfo {
		wh.Access = append(wh.Access, s.mapServiceOrderAccessInfoToAccess(accessInfo))
	}

	if d.AdvertisemtFlag == model.YesAdvertisemtFlag {
		wh.Advertising = true
	}

	for _, alias := range d.Alias {
		wh.Alias = append(wh.Alias, alias.AliasName)
	}

	for _, ftp := range d.FTPInfo {
		wh.FTPAccounts = append(wh.FTPAccounts, ftp.FTPAccount)
	}

	return &wh
}

func (s *service) provisionOrder(order model.Order) {
	ctx := context.Background()

	now := time.Now().Add(1 * time.Millisecond).Unix()
	err := resolver.SendToProvisionSystem(order)

	if err != nil {
		order.Events = append(order.Events, model.Event{
			Type:        model.ProvisionDeclinedEvent,
			Date:        model.CustomTimeStamp(now),
			Description: err.Error(),
		})

		order.Status = model.DeclinedOrderStatus

		if _, err = s.repository.Upsert(ctx, order); err != nil {
			nErr := fmt.Errorf("updating order [%s] after provision [%s]", order.ID, err)
			s.logger.Error(nErr.Error())
		}

		return
	}

	order.Status = model.CompleteOrderStatus

	order.Events = append(order.Events, model.Event{
		Type:        model.ProvisionCompletedEvent,
		Date:        model.CustomTimeStamp(now),
		Description: "provision completed successfully",
	})

	if _, err = s.repository.Upsert(ctx, order); err != nil {
		nErr := fmt.Errorf("updating order [%s] after provision [%s]", order.ID, err)
		s.logger.Error(nErr.Error())
	}

	// reply order to commercial system
}
