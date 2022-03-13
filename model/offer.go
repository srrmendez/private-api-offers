package model

type CustomTimeStamp int64
type ClientType string
type PayModeType string

const (
	IndividualClienType  string = "INDIVIDUAL"
	CorporativeClienType string = "CORPORATIVE"

	PrepaidPayMode  PayModeType = "PREPAID"
	PostpaidPayMode PayModeType = "POSTPAID"
	AllPayMode      PayModeType = "ALL"
)

type Offer struct {
	ID         string          `json:"id" bson:"_id"`
	ExternalID *string         `json:"external_id,omitempty" bson:"external_id,omitempty"`
	CreatedAt  CustomTimeStamp `json:"created_at" bson:"created_at"`
	UpdatedAt  CustomTimeStamp `json:"updated_at" bson:"updated_at"`
	Name       string          `json:"name" bson:"name"`
	Code       *string         `json:"code,omitempty" bson:"code,omitempty"`
	ClientType ClientType      `json:"client_type" bson:"client_type"`
	Paymode    PayModeType     `json:"pay_mode" bson:"pay_mode"`
	StandAlone bool            `json:"standalone" bson:"standalone"`

	Category       string          `json:"category" bson:"category"`
	EffectiveDate  CustomTimeStamp `json:"-" bson:"effective_date"`
	ExpirationDate CustomTimeStamp `json:"expiration_date" bson:"expiration_date"`

	Description *string           `json:"description,omitempty" bson:"description,omitempty"`
	MonthlyFee  float64           `json:"monthly_fee" bson:"monthly_fee"`
	OneOfFee    float64           `json:"one_of_fee" bson:"one_of_fee"`
	Childrens   []Offer           `json:"childrens,omitempty" bson:"childrens,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" bson:"metadata,omitempty"`
}
