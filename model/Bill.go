package model

import "time"

type Bill struct {
	ID       string    `json:"id"`
	Datetime time.Time `json:"datetime"`
	BillDate time.Time `json:"bill_date"`
	BillID   string    `json:"bill_id"`
	BillType struct {
		Status int    `json:"status"`
		Label  string `json:"label"`
	} `json:"bill_type"`
	Table  string `json:"table"`
	Status struct {
		Datetime    time.Time `json:"datetime"`
		Status      int       `json:"status"`
		Label       string    `json:"label"`
		Description string    `json:"description"`
		By          string    `json:"by"`
	} `json:"status"`
	StatusHistory []struct {
		Datetime    time.Time `json:"datetime"`
		Status      int       `json:"status"`
		Label       string    `json:"label"`
		Description string    `json:"description"`
		By          string    `json:"by"`
	} `json:"status_history"`
	Customer struct {
		Name         string `json:"name"`
		MemberID     string `json:"member_id"`
		Aka          string `json:"aka"`
		Rank         int    `json:"rank"`
		CustomerType string `json:"customer_type"`
	} `json:"customer"`
	Order []struct {
		ID            string    `json:"id"`
		Datetime      time.Time `json:"datetime"`
		Name          string    `json:"name"`
		Image         string    `json:"image"`
		Price         float64   `json:"price"`
		Category      int       `json:"category"`
		Subcategory   int       `json:"subcategory"`
		Size          string    `json:"size"`
		Quantity      int       `json:"quantity"`
		OrderTimeline []struct {
			Datetime    time.Time `json:"datetime"`
			Status      int       `json:"status"`
			Label       string    `json:"label"`
			Description string    `json:"description"`
			By          string    `json:"by"`
		} `json:"order_timeline"`
		Status struct {
			Datetime    time.Time `json:"datetime"`
			Status      int       `json:"status"`
			Label       string    `json:"label"`
			Description string    `json:"description"`
			By          string    `json:"by"`
		} `json:"status"`
		Seller         string `json:"seller"`
		SellerID       string `json:"seller_id"`
		SellerName     string `json:"seller_name"`
		SellerUsername string `json:"seller_username"`
		SellerAka      string `json:"seller_aka"`
		Package        bool   `json:"package"`
		Points         int    `json:"points"`
		Commission     bool   `json:"commission"`
	} `json:"order"`
	Operator             string        `json:"operator"`
	DiscountList         []interface{} `json:"discount_list"`
	Discount             float64       `json:"discount"`
	Price                float64       `json:"price"`
	Points               float64       `json:"points"`
	PriceWithDiscount    float64       `json:"price_with_discount"`
	PaymentProof         string        `json:"payment_proof"`
	ServiceCharge        float64       `json:"service_charge"`
	Vat                  float64       `json:"vat"`
	Rounding             float64       `json:"rounding"`
	Total                float64       `json:"total"`
	Payment              string        `json:"payment"`
	PercentServiceCharge int8          `json:"percent_service_charge"`
	IsServiceCharge      bool          `json:"is_service_charge"`
	IsVat                bool          `json:"is_vat"`
}
