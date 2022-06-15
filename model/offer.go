package model

type ClientType string
type PayModeType string
type CategoryType string
type OfferType string

const (
	IndividualClienType  ClientType = "INDIVIDUAL"
	CorporativeClienType ClientType = "CORPORATIVE"

	PrepaidPayMode  PayModeType = "PREPAID"
	PostpaidPayMode PayModeType = "POSTPAID"
	AllPayMode      PayModeType = "ALL"

	CategoryTypeDataCenter  CategoryType = "DATACENTER"
	CategoryTypeYellowPages CategoryType = "YELLOW_PAGES"

	OfferTypeVPS               OfferType = "VPS"
	OfferTypeWebHosting        OfferType = "WEB_HOSTING"
	OfferTypeVirtualDataCenter OfferType = "VIRTUAL_DATA_CENTER"
	OfferTypeHouseLeasing      OfferType = "HOUSE_LEASING"
	OfferTypeDedicatedServer   OfferType = "DEDICATED_SERVER"
	OfferTypeYellowPages       OfferType = "YELLOW_PAGES"
)

type Offer struct {
	ID          string      `json:"id" bson:"_id"`
	ExternalID  *string     `json:"external_id,omitempty" bson:"external_id,omitempty"`
	CreatedAt   string      `json:"created_at" bson:"created_at"`
	UpdatedAt   string      `json:"updated_at" bson:"updated_at"`
	Name        string      `json:"name,omitempty" bson:"name,omitempty"`
	ClientType  ClientType  `json:"client_type,omitempty" bson:"client_type,omitempty"`
	Paymentmode PayModeType `json:"payment_mode,omitempty" bson:"payment_mode,omitempty"`

	Category CategoryType `json:"category,omitempty" bson:"category,omitempty"`
	Type     OfferType    `json:"type,omitempty" bson:"type,omitempty"`

	DataCenterResourceAttributtes *DataCenterResourceAttributtes `json:"data_center_resource_attributes,omitempty" bson:"data_center_resource_attributes,omitempty"`

	EffectiveDate  string `json:"-" bson:"effective_date,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"`

	Fare            float64  `json:"fare,omitempty" bson:"fare,omitempty"`
	Supplementaries []string `json:"supplementaries,omitempty" bson:"supplementaries,omitempty"`
}

type DataCenterResourceAttributtes struct {
	Ram              ResourceValue `json:"ram,omitempty" bson:"ram,omitempty"`
	HDD              ResourceValue `json:"hdd,omitempty" bson:"hdd,omitempty"`
	CPU              ResourceValue `json:"cpu,omitempty" bson:"cpu,omitempty"`
	Database         ResourceValue `json:"database,omitempty" bson:"database,omitempty"`
	FTP              ResourceValue `json:"ftp,omitempty" bson:"ftp,omitempty"`
	Alias            ResourceValue `json:"alias,omitempty" bson:"alias,omitempty"`
	NetworkInterface ResourceValue `json:"network_interface,omitempty" bson:"network_interface,omitempty"`
	IPAddress        ResourceValue `json:"ip_address,omitempty" bson:"ip_address,omitempty"`
	SaveVM           *bool         `json:"save_vm,omitempty" bson:"save_vm,omitempty"`
}

type ResourceValue struct {
	Quantity *int     `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Amount   *float64 `json:"amount,omitempty" bson:"amount,omitempty"`
	Unit     *string  `json:"unit,omitempty" bson:"unit,omitempty"`
}
