package models

//Street struct
type Street struct {
	Id         int     `json:"id"`
	StreetName string  `json:"streetname"`
	Created    string  `json:"created"`
	Closed     *string `json:"closed"`
	City       City    `json:"city"`
}

//AddStreet struct
type AddStreet struct {
	StreetName string `json:"streetname"`
	Created    string `json:"created"`
	City       City   `json:"city"`
}

//Street_count  struct
type Street_count struct {
	Values []Street `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}

//StreetClose struct
type StreetClose struct {
	Id    int    `json:"id"`
	Close string `json:"close"`
}

//UpdStreet struct
type UpdStreet struct {
	Id         int    `json:"id"`
	StreetName string `json:"streetname"`
	Created    string `json:"created"`
}
