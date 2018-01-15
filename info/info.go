package info

type DeliveryDetails struct {
	Fee   float64 `yaml:"fee"`
	Promo string  `yaml:"promo"`
}

type PageInfo struct {
	SkuID        int
	ItemID       int
	SellerID     int
	DeliveryInfo map[string]DeliveryDetails
}
