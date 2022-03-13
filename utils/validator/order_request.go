package validator

import (
	"errors"
	"fmt"

	"github.com/srrmendez/private-api-order/model"
)

func NewOrderRequestValidator() *OrderRequestValidator {
	dataCenterOrderTypes := []model.OrderType{model.DedicatedServerOrder, model.WebHostingOrder, model.HouseLeasingOrder, model.VPSOrder,
		model.VirtualDataCenterOrder}

	return &OrderRequestValidator{
		categories: map[model.CategoryType][]model.OrderType{
			model.DataCenterCategory:       dataCenterOrderTypes,
			model.YellowPagesCategory:      {model.YellowPagesOrder},
			model.YellowPagesPrintCategory: {model.YellowPagesPrintOrder},
		},
	}
}

func (v *OrderRequestValidator) ValidateStatus(status model.OrderStatusType) error {
	sts := []model.OrderStatusType{model.CompleteOrderStatus, model.DeclinedOrderStatus}

	isValid := false

	for _, st := range sts {
		if status == st {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("incorrect status value posible values are %s, %s", model.CompleteOrderStatus, model.DeclinedOrderStatus)
	}

	return nil
}

func (v *OrderRequestValidator) Validate(order model.Order) error {
	if _, ok := v.categories[order.Category]; !ok {
		return fmt.Errorf("incorrect category type posible values are %s, %s", model.DataCenterCategory, model.YellowPagesCategory)
	}

	isValid := false

	for _, orderType := range v.categories[order.Category] {
		if orderType == order.Type {
			isValid = true
			break
		}
	}

	if !isValid {
		oErr := "incorrect category type posible values are"

		for i := range v.categories[order.Category] {
			if i == len(v.categories[order.Category])-1 {
				oErr = fmt.Sprintf("%s %s", oErr, v.categories[order.Category][i])
				continue
			}

			oErr = fmt.Sprintf("%s %s,", oErr, v.categories[order.Category][i])
		}

		return errors.New(oErr)
	}

	if order.CustomerType != model.IndividualCustomer && order.CustomerType != model.CorporativeCustomer {
		return fmt.Errorf("incorrect client type posible values are %s, %s", model.IndividualCustomer, model.CorporativeCustomer)
	}

	switch order.Type {
	case model.DedicatedServerOrder:
		return v.validateDedicatedServerOrderType(order.DedicatedServer, order.CustomerType)
	case model.HouseLeasingOrder:
		return v.validateHouseLeasingOrderType(order.HouseLeasing, order.CustomerType)
	case model.VirtualDataCenterOrder:
		return v.validateVirtualDataCenter(order.VirtualDataCenter, order.CustomerType)
	case model.WebHostingOrder:
		return v.validateWebHosting(order.WebHosting, order.CustomerType)
	case model.VPSOrder:
		return v.validateVPSRequest(order.VPS, order.CustomerType)
	default:
		return nil
	}
}

func (v *OrderRequestValidator) validateAccess(access model.Access, CustomerType model.CustomerType) error {
	if CustomerType == model.IndividualCustomer {
		if access.IPAddress != nil {
			return errors.New("ip address cannot be combined with natural client")
		}

		if access.ServiceNumber != nil {
			return errors.New("service number cannot be combined with natural client")
		}

		if access.Access != nil {
			return errors.New("access cannot be combined with natural client")
		}

		if access.Account == nil || len(*access.Account) == 0 {
			return errors.New("access account must be provided")
		}

		return nil
	}

	if access.Account != nil {
		return errors.New("access account cannot be combined with corporate client")
	}

	if access.Access == nil || (*access.Access != model.InternalAccess && *access.Access != model.ExternalAccess) {
		return fmt.Errorf("incorrect access type posible values are %s, %s", model.InternalAccess, model.ExternalAccess)
	}

	return nil
}

func (v *OrderRequestValidator) validateContacts(contacts []model.Contact) error {
	for _, contact := range contacts {
		if contact.Email == "" {
			return errors.New("contact email must be provided")
		}

		if contact.Firstname == "" {
			return errors.New("contact first name must be provided")
		}

		if contact.Lastname == "" {
			return errors.New("contact last name must be provided")
		}

		if contact.Phone == "" {
			return errors.New("contact phone must be provided")
		}

		isValid := false
		for _, cTypes := range []model.ContactType{model.ModeratorContact, model.NotificationContact, model.TechnicianContact} {
			if cTypes == contact.Type {
				isValid = true
				break
			}
		}

		if !isValid {
			return fmt.Errorf("incorrect contact type posible values are %s, %s, %s", model.ModeratorContact, model.NotificationContact, model.TechnicianContact)
		}
	}

	return nil
}

func (v *OrderRequestValidator) validateDedicatedServerOrderType(ds *model.DedicatedServer, CustomerType model.CustomerType) error {
	if ds == nil {
		return errors.New("dedicated server object must be provided")
	}

	err := v.validateAccess(ds.Access, CustomerType)
	if err != nil {
		return err
	}

	if len(ds.ServerName) == 0 {
		return errors.New("server name must be provided")
	}

	return v.validateContacts(ds.Contacts)
}

func (v *OrderRequestValidator) validateHouseLeasingOrderType(ds *model.HouseLeasing, CustomerType model.CustomerType) error {
	if ds == nil {
		return errors.New("dedicated server object must be provided")
	}

	err := v.validateAccess(ds.Access, CustomerType)
	if err != nil {
		return err
	}

	if len(ds.ServerName) == 0 {
		return errors.New("server name must be provided")
	}

	return v.validateContacts(ds.Contacts)
}

func (v *OrderRequestValidator) validateVirtualDataCenter(ds *model.VirtualDataCenter, CustomerType model.CustomerType) error {
	if ds == nil {
		return errors.New("dedicated server object must be provided")
	}

	err := v.validateAccess(ds.Access, CustomerType)
	if err != nil {
		return err
	}

	if len(ds.ServerName) == 0 {
		return errors.New("server name must be provided")
	}

	return v.validateContacts(ds.Contacts)
}

func (v *OrderRequestValidator) validateWebHosting(ds *model.WebHosting, CustomerType model.CustomerType) error {
	if ds == nil {
		return errors.New("dedicated server object must be provided")
	}

	for _, access := range ds.Access {
		err := v.validateAccess(access, CustomerType)
		if err != nil {
			return err
		}
	}

	if ds.Domain == "" {
		return errors.New("domain must be provided")
	}

	if len(ds.Alias) == 0 {
		return errors.New("aliases must be provided")
	}

	for i := range ds.Alias {
		if ds.Alias[i] == "" {
			return errors.New("alias must be provided")
		}
	}

	if len(ds.FTPAccounts) == 0 {
		return errors.New("ftp accounts must be provided")
	}

	for i := range ds.FTPAccounts {
		if ds.FTPAccounts[i] == "" {
			return errors.New("ftp account must be provided")
		}
	}

	if ds.Database == "" {
		return errors.New("database must be provided")
	}

	if ds.LMS != nil && *ds.LMS == "" {
		return errors.New("lms must be provided")
	}

	if ds.CMS != nil && *ds.CMS == "" {
		return errors.New("cms must be provided")
	}

	if ds.CMF != nil && *ds.CMF == "" {
		return errors.New("cmf must be provided")
	}

	if ds.Framework != nil && *ds.Framework == "" {
		return errors.New("framework must be provided")
	}

	if ds.ProgrammingLanguage != nil && *ds.ProgrammingLanguage == "" {
		return errors.New("programming language must be provided")
	}

	return v.validateContacts(ds.Contacts)
}

func (v *OrderRequestValidator) validateVPSRequest(vps *model.VPS, CustomerType model.CustomerType) error {
	if vps == nil {
		return errors.New("vps object must be provided")
	}

	err := v.validateAccess(vps.Access, CustomerType)
	if err != nil {
		return err
	}

	if CustomerType == model.IndividualCustomer {
		if vps.Forum == nil {
			return errors.New("forum must be provided")
		}

		return nil
	}

	if vps.ServerName == nil || len(*vps.ServerName) == 0 {
		return errors.New("server name must be provided")
	}

	if vps.SLA == nil {
		return errors.New("sla must be provided")
	}

	return v.validateContacts(vps.Contacts)
}
