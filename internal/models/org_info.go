package models

//OrgInfo struct
type OrgInfo struct {
	Id          int    `json:"id"`
	OIName      string `json:"oiname"`
	OIBin       string `json:"oibin"`
	OIAddr      string `json:"oiaddr"`
	OIBank      Bank   `json:"oibank"`
	OIAccNumber string `json:"oiaccnumber"`
	OIFName     string `json:"oifname"`
}

//AddOrgInfo struct
type AddOrgInfo struct {
	OIName      string `json:"oiname"`
	OIBin       string `json:"oibin"`
	OIAddr      string `json:"oiaddr"`
	OIBank      Bank   `json:"oibank"`
	OIAccNumber string `json:"oiaccnumber"`
	OIFName     string `json:"oifname"`
}

//OrgInfo_count  struct
type OrgInfo_count struct {
	Values []OrgInfo `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}
