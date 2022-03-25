package model

type BSSClientType string
type BssPaymodeType string
type BssStatus string
type BssSaleType string

const (
	IndividualBssOfferClient  BSSClientType = "1"
	CorporativeBssOfferClient BSSClientType = "2"
	FamilyBssOfferClient      BSSClientType = "3"

	PrepaidBssPaymode  BssPaymodeType = "0"
	PostpaidBssPaymode BssPaymodeType = "1"
	AllBssPaymode      BssPaymodeType = "2"

	DraftBssStatus      BssStatus = "0"
	TestBssStatus       BssStatus = "1"
	ReleaseBssStatus    BssStatus = "2"
	SuspendBssStatus    BssStatus = "3"
	RetirementBssStatus BssStatus = "4"

	NoSaleAlone BssSaleType = "0"
	saleAlone   BssSaleType = "1"
)

type BssAttribute struct {
	Code        string `json:"attr_code"`
	Value       string `json:"attr_value"`
	Type        string `json:"attr_type"`
	Description string `json:"attr_value_desc"`
}

type BssAttributeList struct {
	Attribute []BssAttribute `json:"attr"`
}

type BssAttached struct {
	ID           string `json:"offeringID"`
	RelationType string `json:"releationType"`
}

type BssRelationshipList struct {
	Attached []BssAttached `json:"attached"`
}

type BssCatalog struct {
	ID   string `json:"catalog_id"`
	Name string `json:"catalog_name"`
}

type BssCatalogList struct {
	Catalogs []BssCatalog `json:"catalog"`
}

type BssOffer struct {
	ID             string               `json:"offer_id"`
	Code           string               `json:"offer_code"`
	PrimaryFlag    string               `json:"primary_flag"`
	BundleFlag     string               `json:"bundle_flag"`
	Name           string               `json:"offer_name"`
	Category       string               `json:"offer_cata"`
	EffectiveDate  *string              `json:"eff_date"`
	ExpirationDate *string              `json:"exp_date"`
	Status         BssStatus            `json:"status"`
	OneOfFee       float64              `json:"oneoff_fee"`
	MontlyFee      float64              `json:"monthly_fee"`
	Description    string               `json:"offer_desc"`
	Type           BSSClientType        `json:"offer_type"`
	Attributes     *BssAttributeList    `json:"attrList"`
	Relationships  *BssRelationshipList `json:"relationshipList"`
	PayMode        BssPaymodeType       `json:"pay_mode"`
	OnSale         BssSaleType          `json:"on_sale"`
	Catalogs       BssCatalogList       `json:"catalogList"`
	TopupFee       string               `json:"topupFee"`
}

type BssSyncOfferRequest struct {
	SyncOffers []BssOffer `json:"syncOffers"`
}
