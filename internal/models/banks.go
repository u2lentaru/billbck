package models

//Bank struct
type Bank struct {
	Id        int    `json:"id"`
	BankName  string `json:"bankname"`
	BankDescr string `json:"bankdescr"`
	Mfo       string `json:"mfo"`
}

//AddBank struct
type AddBank struct {
	BankName  string `json:"bankname"`
	BankDescr string `json:"bankdescr"`
	Mfo       string `json:"mfo"`
}

//Bank_count  struct
type Bank_count struct {
	Values []Bank `json:"values"`
	Count  int    `json:"count"`
	Auth   Auth   `json:"auth"`
}
