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
	ActivationFate  float64  `json:"activation_fare,omitempty" bson:"activation_fare,omitempty"`
	Supplementaries []string `json:"supplementaries,omitempty" bson:"supplementaries,omitempty"`
}

type DataCenterResourceAttributtes struct {
	RAM                 *RAM      `json:"ram,omitempty" bson:"ram,omitempty"`
	HDD                 *HDD      `json:"hdd,omitempty" bson:"hdd,omitempty"`
	CPUQty              *int      `json:"cpu_quantity,omitempty" bson:"cpu_quantity,omitempty"`
	Database            *Database `json:"database,omitempty" bson:"database,omitempty"`
	FTPQty              *int      `json:"ftp_quantity,omitempty" bson:"ftp_quantity,omitempty"`
	AliasQty            *int      `json:"alias_quantity,omitempty" bson:"alias_quantity,omitempty"`
	NetworkInterfaceQty *int      `json:"network_interface_qty,omitempty" bson:"network_interface_qty,omitempty"`
	PublicIPAddress     *string   `json:"public_ip_address,omitempty" bson:"public_ip_address,omitempty"`
	LANIPAddress        *string   `json:"lan_ip_address,omitempty" bson:"lan_ip_address,omitempty"`
	WANIPAddress        *string   `json:"wan_ip_address,omitempty" bson:"wan_ip_address,omitempty"`
	VPN                 *VPN      `json:"vpn,omitempty" bson:"vpn,omitempty"`
	DNS                 *DNS      `json:"dns,omitempty" bson:"dns,omitempty"`
	Bandwidth           *BandWith `json:"bandwidth,omitempty" bson:"bandwidth,omitempty"`
	SaveVM              *bool     `json:"save_vm,omitempty" bson:"save_vm,omitempty"`
	AccessType          *string   `json:"access_type,omitempty" bson:"access_type,omitempty"`
}

type RAM struct {
	Amount float64 `json:"amount" bson:"amount"`
	Unit   string  `json:"unit" bson:"unit"`
}

type HDD struct {
	Amount float64 `json:"amount" bson:"amount"`
	Unit   string  `json:"unit" bson:"unit"`
}

type Database struct {
	Quantity int     `json:"quantity" bson:"quantity"`
	Amount   float64 `json:"amount" bson:"amount"`
	Unit     string  `json:"unit" bson:"unit"`
}

type BandWith struct {
	Amount float64 `json:"amount" bson:"amount"`
	Unit   string  `json:"unit" bson:"unit"`
}

type VPN struct {
	IPAddress string `json:"ip_address" bson:"ip_address"`
	Name      string `json:"name" bson:"name"`
}

type DNS struct {
	Name string `json:"name" bson:"name"`
	DNS  string `json:"dns" bson:"dns"`
}
