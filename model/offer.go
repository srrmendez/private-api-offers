package model

import (
	"encoding/json"
	"time"
)

type CustomTimeStamp int64
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
	ID         string          `json:"id" bson:"_id"`
	ExternalID *string         `json:"external_id,omitempty" bson:"external_id,omitempty"`
	CreatedAt  CustomTimeStamp `json:"created_at" bson:"created_at"`
	UpdatedAt  CustomTimeStamp `json:"updated_at" bson:"updated_at"`
	Name       string          `json:"name,omitempty" bson:"name,omitempty"`
	Code       *string         `json:"code,omitempty" bson:"code,omitempty"`
	ClientType ClientType      `json:"client_type,omitempty" bson:"client_type,omitempty"`
	Paymode    PayModeType     `json:"pay_mode,omitempty" bson:"pay_mode,omitempty"`
	StandAlone bool            `json:"standalone,omitempty" bson:"standalone,omitempty"`

	Category CategoryType `json:"category,omitempty" bson:"category,omitempty"`
	Type     OfferType    `json:"type,omitempty" bson:"type,omitempty"`

	EffectiveDate  *CustomTimeStamp `json:"-" bson:"effective_date,omitempty"`
	ExpirationDate *CustomTimeStamp `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"`

	Description *string           `json:"description,omitempty" bson:"description,omitempty"`
	MonthlyFee  float64           `json:"monthly_fee,omitempty" bson:"monthly_fee,omitempty"`
	OneOfFee    float64           `json:"one_of_fee,omitempty" bson:"one_of_fee,omitempty"`
	Childrens   []Offer           `json:"childrens,omitempty" bson:"childrens,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

func (c CustomTimeStamp) MarshalJSON() ([]byte, error) {
	tstamp := int64(c)

	t := time.Unix(tstamp, 0).Format("2006-01-02 15:04:05")

	return json.Marshal(t)
}
