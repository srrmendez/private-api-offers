package model

import (
	"time"
)

type ServiceOrderBusinnessType string
type ServiceOrderDateTime time.Time
type ServiceOrderCustomerType string
type ServiceOrderNetworkType string
type ServiceResponsibleType string
type ServiceOrderActionType string
type ServiceOrderAccessType string
type AdvertisemtFlagType string
type ServiceOrderProvisionType string
type ServiceOrderOfferingType string

const (
	NewSubscriberBusinness               ServiceOrderBusinnessType = "CO064"
	SuspendBusiness                      ServiceOrderBusinnessType = "CO066"
	ResumeBusiness                       ServiceOrderBusinnessType = "CO067"
	ChangePrimaryOfferingBusinness       ServiceOrderBusinnessType = "CO024"
	ChangeSupplementaryOfferingBusinness ServiceOrderBusinnessType = "CO025"
	ChangeSubscriberDataBusinness        ServiceOrderBusinnessType = "CO019"
	DeactivateSubscriberDataBusinness    ServiceOrderBusinnessType = "CO022"
	TransferOwnershipBusinness           ServiceOrderBusinnessType = "CO004"

	ServiceOrderIndividualCustomer    ServiceOrderCustomerType = "1"
	ServiceOrderCorporateCustomerType ServiceOrderCustomerType = "2"

	DataCenterNetwork       ServiceOrderNetworkType = "1"
	YellowPagesWebNetwork   ServiceOrderNetworkType = "2"
	YellowPagesPrintNetwork ServiceOrderNetworkType = "3"

	TechinicianServiceOrderResponsible  ServiceResponsibleType = "1"
	ModeratorServiceOrderResponsible    ServiceResponsibleType = "2"
	NotificationServiceOrderResponsible ServiceResponsibleType = "3"

	AddServiceOrderAction    ServiceOrderActionType = "1"
	DeleteServiceOrderAction ServiceOrderActionType = "2"
	ModifyServiceOrderAction ServiceOrderActionType = "3"
	KeepServiceOrderAction   ServiceOrderActionType = "4"

	InternalServiceOrderAccess ServiceOrderAccessType = "1"
	ExternalServiceOrderAccess ServiceOrderAccessType = "2"
	AAAServiceOrderAccess      ServiceOrderAccessType = "3"

	NoAdvertisemtFlag  AdvertisemtFlagType = "0"
	YesAdvertisemtFlag AdvertisemtFlagType = "1"

	VPSServiceOrderProvision               ServiceOrderProvisionType = "1"
	WebHostingServiceOrderProvision        ServiceOrderProvisionType = "2"
	VirtualDataCenterServiceOrderProvision ServiceOrderProvisionType = "3"
	DedicatedServerServiceOrderProvision   ServiceOrderProvisionType = "4"
	HouseLeasingServiceOrderProvision      ServiceOrderProvisionType = "5"

	PrimaryServiceOrderOffering       ServiceOrderOfferingType = "1"
	SupplementaryServiceOrderOffering ServiceOrderOfferingType = "2"
	OtherServiceOrderOffering         ServiceOrderOfferingType = "9"
)

type AdditionalProperty struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type OrderInfo struct {
	OrderID              string                    `json:"orderId"`
	OrderBusinessType    ServiceOrderBusinnessType `json:"orderBusiType"`
	CreateDateTime       *ServiceOrderDateTime     `json:"createDate"`
	CustomerType         ServiceOrderCustomerType  `json:"custType"`
	AdditionalProperties []AdditionalProperty      `json:"additionalProperty"`
}

type ServiceOrderName struct {
	Firstname  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
}

type ServiceOrderResponsible struct {
	ResponsibleID      string                 `json:"responsibleId"`
	ResponsibleType    ServiceResponsibleType `json:"responsibleType"`
	FullName           string                 `json:"fullName"`
	Email              string                 `json:"email"`
	Phone              string                 `json:"phone"`
	ActionType         ServiceOrderActionType `json:"actionType"`
	AdditionalProperty []AdditionalProperty   `json:"additionalProperty"`
}

type ServiceOrderAccessInfo struct {
	AccessID           string                 `json:"accessId"`
	AccessType         ServiceOrderAccessType `json:"accessType"`
	AccessAcount       *string                `json:"accessAcount"`
	IPAddress          *string                `json:"ipAddress"`
	ServiceNumber      *string                `json:"serviceNumber"`
	ActionType         ServiceOrderActionType `json:"actionType"`
	AdditionalProperty []AdditionalProperty   `json:"additionalProperty"`
}

type ServiceOrderVPS struct {
	ServerName         *string                   `json:"serverName"`
	Responsible        []ServiceOrderResponsible `json:"responsible"`
	AccessInfo         ServiceOrderAccessInfo    `json:"accessInfo"`
	SLAFlag            *bool                     `json:"slaFlag"`
	ForumFlag          *bool                     `json:"forumFlag"`
	AdditionalProperty []AdditionalProperty      `json:"additionalProperty"`
}

type ServiceOrderAliasInfo struct {
	AliasID            string                 `json:"aliasId"`
	AliasName          string                 `json:"aliasName"`
	ActionType         ServiceOrderActionType `json:"actionType"`
	AdditionalProperty []AdditionalProperty   `json:"additionalProperty"`
}

type ServiceOrderFTPInfo struct {
	FTPID              string                 `json:"ftpId"`
	FTPAccount         string                 `json:"ftpAccount"`
	ActionType         ServiceOrderActionType `json:"actionType"`
	AdditionalProperty []AdditionalProperty   `json:"additionalProperty"`
}

type ServiceOrderWebHosting struct {
	MainDomain          string                    `json:"mainDomain"`
	ProgrammingLanguage *string                   `json:"programmingLanguage"`
	DataBase            string                    `json:"dataBase"`
	AdvertisemtFlag     AdvertisemtFlagType       `json:"advertisemtFlag"`
	WebSiteType         *string                   `json:"webSiteType"`
	FrameWorkName       *string                   `json:"frameWorkName"`
	CMSName             *string                   `json:"cmsName"`
	LMSVersion          *string                   `json:"lmsVersion"`
	CMFVersion          *string                   `json:"cmfVersion"`
	Responsible         []ServiceOrderResponsible `json:"responsible"`
	TopicToDiscuss      *string                   `json:"topicsToDiscuss"`
	WebsiteGoal         *string                   `json:"websiteGoal"`
	Remark              *string                   `json:"remark"`
	Alias               []ServiceOrderAliasInfo   `json:"alias"`
	FTPInfo             []ServiceOrderFTPInfo     `json:"ftpInfo"`
	AccessInfo          []ServiceOrderAccessInfo  `json:"accessInfo"`
	AdditionalProperty  []AdditionalProperty      `json:"additionalProperty"`
}

type ServiceOrderVirtualDataCenter struct {
	ServerName         string                    `json:"serverName"`
	AccessInfo         ServiceOrderAccessInfo    `json:"accessInfo"`
	Responsible        []ServiceOrderResponsible `json:"responsible"`
	SLAFlag            bool                      `json:"slaFlag"`
	AdditionalProperty []AdditionalProperty      `json:"additionalProperty"`
}

type ServiceOrderDedicatedServer struct {
	ServerName         string                    `json:"serverName"`
	AccessInfo         ServiceOrderAccessInfo    `json:"accessInfo"`
	Responsible        []ServiceOrderResponsible `json:"responsible"`
	SLAFlag            bool                      `json:"slaFlag"`
	AdditionalProperty []AdditionalProperty      `json:"additionalProperty"`
}

type ServiceOrderHousingLeasing struct {
	ServerName         string                    `json:"serverName"`
	AccessInfo         ServiceOrderAccessInfo    `json:"accessInfo"`
	Responsible        []ServiceOrderResponsible `json:"responsible"`
	SLAFlag            bool                      `json:"slaFlag"`
	AdditionalProperty []AdditionalProperty      `json:"additionalProperty"`
}

type SubscriberInfo struct {
	SubscriberID         string                         `json:"subId"`
	SubName              ServiceOrderName               `json:"subName"`
	ServiceNumber        string                         `json:"serviceNumber"`
	NetworkType          ServiceOrderNetworkType        `json:"networkType"`
	InitialDate          string                         `json:"initialDate"`
	VPS                  *ServiceOrderVPS               `json:"vps"`
	WebHosting           *ServiceOrderWebHosting        `json:"webHosting"`
	VirtualDataCenter    *ServiceOrderVirtualDataCenter `json:"virtualDataCenter"`
	DedicatedServer      *ServiceOrderDedicatedServer   `json:"dedicatedServer"`
	HousingLeasing       *ServiceOrderHousingLeasing    `json:"housingLeasing"`
	ProvisionServiceType ServiceOrderProvisionType      `json:"provisionServiceType"`
	AdditionalProperty   []AdditionalProperty           `json:"additionalProperty"`
}

type ServiceOrderOfferInst struct {
	ActionType   ServiceOrderActionType `json:"actionType"`
	AttrCode     string                 `json:"attrCode"`
	AttrNewValue *string                `json:"attrNewValue"`
	AttrOldValue *string                `json:"attrOldValue"`
}

type ServiceOrderOffering struct {
	OfferingID         string                   `json:"offeringID"`
	OfferingName       string                   `json:"offeringName"`
	OfferingType       ServiceOrderOfferingType `json:"offeringType"`
	ActionType         ServiceOrderActionType   `json:"actionType"`
	EffectiveTime      ServiceOrderDateTime     `json:"effectiveTime"`
	ExpirationTime     ServiceOrderDateTime     `json:"expirationTime"`
	OfferInstAttr      []ServiceOrderOfferInst  `json:"offerInstAttr"`
	AdditionalProperty []AdditionalProperty     `json:"additionalProperty"`
}

type ServiceOrderRequest struct {
	OrderInfo          OrderInfo              `json:"orderInfo"`
	SubscriberInfo     SubscriberInfo         `json:"subInfo"`
	OfferingInfo       []ServiceOrderOffering `json:"offeringInfo"`
	AdditionalProperty []AdditionalProperty   `json:"additionalProperty"`
}

type ServiceOrderResponse struct {
	BSSTransactionID string `json:"transaction"`
	OrderID          string `json:"id"`
}

func (t *ServiceOrderDateTime) UnmarshalJSON(data []byte) error {
	dt, err := time.Parse("20060102150405", string(data))
	if err != nil {
		return err
	}

	*t = ServiceOrderDateTime(dt)

	return nil
}
