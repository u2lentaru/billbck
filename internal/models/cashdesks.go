package models

//Cashdesk struct
type Cashdesk struct {
	Id           int    `json:"id"`
	CashdeskName string `json:"cashdeskname"`
}

//AddCashdesk struct
type AddCashdesk struct {
	CashdeskName string `json:"cashdeskname"`
}

//Cashdesk_count  struct
type Cashdesk_count struct {
	Values []Cashdesk `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}
