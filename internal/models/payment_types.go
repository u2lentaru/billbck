package models

//PaymentType struct
type PaymentType struct {
	Id              int    `json:"id"`
	PaymentTypeName string `json:"paymenttypename"`
}

//AddPaymentType struct
type AddPaymentType struct {
	PaymentTypeName string `json:"paymenttypename"`
}

//PaymentType_count  struct
type PaymentType_count struct {
	Values []PaymentType `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
