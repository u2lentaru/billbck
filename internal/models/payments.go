package models

//Payment struct
type Payment struct {
	Id           int         `json:"id"`
	PaymentDate  string      `json:"paymentdate"`
	Contract     Contract    `json:"contract"`
	Object       Object      `json:"object"`
	PaymentType  PaymentType `json:"paymenttype"`
	ChargeType   ChargeType  `json:"chargetype"`
	Cashdesk     Cashdesk    `json:"cashdesk"`
	BundleNumber int         `json:"bundlenumber"`
	Amount       float32     `json:"amount"`
}

//AddPayment struct
type AddPayment struct {
	PaymentDate  string      `json:"paymentdate"`
	Contract     Contract    `json:"contract"`
	Object       Object      `json:"object"`
	PaymentType  PaymentType `json:"paymenttype"`
	ChargeType   ChargeType  `json:"chargetype"`
	Cashdesk     Cashdesk    `json:"cashdesk"`
	BundleNumber int         `json:"bundlenumber"`
	Amount       float32     `json:"amount"`
}

//Payment_count  struct
type Payment_count struct {
	Values []Payment `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}
