package models

//ClaimType struct
type ClaimType struct {
	Id            int    `json:"id"`
	ClaimTypeName string `json:"claimtypename"`
}

//AddClaimType struct
type AddClaimType struct {
	ClaimTypeName string `json:"claimtypename"`
}

//ClaimType_count struct
type ClaimType_count struct {
	Values []ClaimType `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
