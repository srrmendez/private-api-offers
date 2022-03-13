package model

import (
	"encoding/json"
	"time"
)

type OrderStatusType string
type OrderType string
type CategoryType string
type CustomTimeStamp int64
type EventType string
type ContactType string
type CustomerType string
type AccessType string
type SourceType string

const (
	TechnicianContact   ContactType = "TECHNICIAN"
	ModeratorContact    ContactType = "MODERATOR"
	NotificationContact ContactType = "NOTIFICATION"

	PendingOrderStatus  OrderStatusType = "PENDING"
	DeclinedOrderStatus OrderStatusType = "DECLINED"
	CompleteOrderStatus OrderStatusType = "COMPLETED"

	DataCenterCategory       CategoryType = "DATACENTER"
	YellowPagesCategory      CategoryType = "YELLOW_PAGES_WEB"
	YellowPagesPrintCategory CategoryType = "YELLOW_PAGES_PRINT"

	VPSOrder               OrderType = "VPS"
	WebHostingOrder        OrderType = "WEB_HOSTING"
	VirtualDataCenterOrder OrderType = "VIRTUAL_DATA_CENTER"
	HouseLeasingOrder      OrderType = "HOUSE_LEASING"
	DedicatedServerOrder   OrderType = "DEDICATED_SERVER"

	YellowPagesOrder      OrderType = "YELLOW_PAGES_WEB"
	YellowPagesPrintOrder OrderType = "YELLOW_PAGES_PRINT"

	CommercialValidationSendEvent      EventType = "COMMERCIAL_VALIDATION_SENDED"
	CommercialValidationCompletedEvent EventType = "COMMERCIAL_VALIDATION_COMPLETED"
	CommercialValidationDeclinedEvent  EventType = "COMMERCIAL_VALIDATION_DECLINED"

	ProvisionSendEvent      EventType = "PROVISION_SENDED"
	ProvisionCompletedEvent EventType = "PROVISION_COMPLETED"
	ProvisionDeclinedEvent  EventType = "PROVISION_DECLINED"

	IndividualCustomer  CustomerType = "INDIVIDUAL"
	CorporativeCustomer CustomerType = "CORPORATIVE"

	ExternalAccess AccessType = "EXTERNAL"
	InternalAccess AccessType = "INTERNAL"

	WebPortalSource        SourceType = "WEB_PORTAL"
	CommercialSystemSource SourceType = "COMMERCIAL_SYSTEM"
)

type Event struct {
	Type        EventType       `json:"type" bson:"type"`
	Date        CustomTimeStamp `json:"date" bson:"date"`
	Description string          `json:"description" bson:"description"`
}

type Contact struct {
	Firstname string      `json:"firstname"`
	Lastname  string      `json:"lastname"`
	Email     string      `json:"email"`
	Phone     string      `json:"phone"`
	Type      ContactType `json:"type"`
}

type Access struct {
	Account *string `json:"account,omitempty" bson:"account,omitempty"` // natural client

	Access        *AccessType `json:"access"`
	IPAddress     *string     `json:"ip_address,omitempty" bson:"ip_address,omitempty"`
	ServiceNumber *string     `json:"service_number,omitempty" bson:"service_number,omitempty"`
}

type VPS struct {
	ServerName *string   `json:"server_name,omitempty" bson:"server_name,omitempty"`
	Contacts   []Contact `json:"contacts,omitempty" bson:"contacts,omitempty"`
	Access     Access    `json:"access" bson:"access"`
	Forum      *bool     `json:"forum,omitempty" bson:"forum"`
	SLA        *bool     `json:"sla,omitempty" bson:"sla,omitempty"`
}

type HouseLeasing struct {
	ServerName string    `json:"server_name" bson:"server_name"`
	Contacts   []Contact `json:"contacts,omitempty" bson:"contacts,omitempty"`
	Access     Access    `json:"access" bson:"access"`
	SLA        bool      `json:"sla" bson:"sla"`
}

type VirtualDataCenter struct {
	ServerName string    `json:"server_name" bson:"server_name"`
	Contacts   []Contact `json:"contacts,omitempty" bson:"contacts,omitempty"`
	Access     Access    `json:"access" bson:"access"`
	SLA        bool      `json:"sla" bson:"sla"`
}

type DedicatedServer struct {
	ServerName string    `json:"server_name" bson:"server_name"`
	Contacts   []Contact `json:"contacts,omitempty" bson:"contacts,omitempty"`
	Access     Access    `json:"access" bson:"access"`
	SLA        bool      `json:"sla" bson:"sla"`
}

type WebHosting struct {
	Access []Access `json:"access" bson:"access"`

	Domain              string  `json:"domain" bson:"domain"`
	ProgrammingLanguage *string `json:"programming_language,omitempty" bson:"programming_language,omitempty"`
	Database            string  `json:"database" bson:"database"`
	Advertising         bool    `json:"advertising" bson:"advertising"`
	WebsiteType         *string `json:"website_type,omitempty" bson:"website_type,omitempty"`
	Framework           *string `json:"framework,omitempty" bson:"framework,omitempty"`
	CMS                 *string `json:"cms,omitempty" bson:"cms,omitempty"`
	LMS                 *string `json:"lms,omitempty" bson:"lms,omitempty"`
	CMF                 *string `json:"cmf,omitempty" bson:"cmf,omitempty"`

	Topic       *string `json:"topic,omitempty" bson:"topic,omitempty"`
	WebsiteGoal *string `json:"website_goal,omitempty" bson:"website_goal,omitempty"`
	Remark      *string `json:"remark,omitempty" bson:"remark,omitempty"`

	Alias       []string `json:"alias,omitempty" bson:"alias,omitempty"`
	FTPAccounts []string `json:"ftp_accounts,omitempty" bson:"ftp_accounts,omitempty"`

	Contacts []Contact `json:"contacts,omitempty" bson:"contacts,omitempty"`
}

type Order struct {
	ID         string          `json:"id" bson:"_id"`
	AppID      string          `json:"-" bson:"app_id"`
	ExternalID *string         `json:"external_id,omitempty" bson:"external_id,omitempty"`
	Source     SourceType      `json:"source" bson:"source"`
	CreatedAt  CustomTimeStamp `json:"created_at" bson:"created_at"` // timestamp
	UpdatedAt  CustomTimeStamp `json:"updated_at" bson:"updated_at"` // timestamp
	Status     OrderStatusType `json:"status" bson:"status"`
	Category   CategoryType    `json:"category" bson:"category"`
	Type       OrderType       `json:"type" bson:"type"`

	Events       []Event      `json:"events" bson:"events"`
	CustomerType CustomerType `json:"customer_type" bson:"customer_type"`

	VPS               *VPS               `json:"vps,omitempty" bson:"vps,omitempty"`
	WebHosting        *WebHosting        `json:"web_hosting,omitempty" bson:"web_hosting,omitempty"`
	VirtualDataCenter *VirtualDataCenter `json:"virtual_data_center,omitempty" bson:"virtual_data_center,omitempty"`
	DedicatedServer   *DedicatedServer   `json:"dedicated_server,omitempty" bson:"dedicated_server,omitempty"`
	HouseLeasing      *HouseLeasing      `json:"house_leasing,omitempty" bson:"house_leasing,omitempty"`
}

func (c CustomTimeStamp) MarshalJSON() ([]byte, error) {
	tstamp := int64(c)

	t := time.Unix(tstamp, 0).Format("2006-01-02 15:04:05")

	return json.Marshal(t)
}

type UpdateOrderRequest struct {
	Status      OrderStatusType `json:"status"`
	Description string          `json:"description"`
}
