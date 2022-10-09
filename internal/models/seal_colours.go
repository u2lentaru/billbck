package models

//SealColour struct
type SealColour struct {
	Id             int    `json:"id"`
	SealColourName string `json:"sealcolourname"`
}

//AddSealColour struct
type AddSealColour struct {
	SealColourName string `json:"sealcolourname"`
}

//SealColour_count  struct
type SealColour_count struct {
	Values []SealColour `json:"values"`
	Count  int          `json:"count"`
	Auth   Auth         `json:"auth"`
}
