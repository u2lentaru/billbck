package models

//SealStatus struct
type SealStatus struct {
	Id             int    `json:"id"`
	SealStatusName string `json:"sealstatusname"`
}

//AddSealStatus struct
type AddSealStatus struct {
	SealStatusName string `json:"sealstatusname"`
}

//SealStatus_count  struct
type SealStatus_count struct {
	Values []SealStatus `json:"values"`
	Count  int          `json:"count"`
	Auth   Auth         `json:"auth"`
}
