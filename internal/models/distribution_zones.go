package models

//DistributionZone struct
type DistributionZone struct {
	Id                   int    `json:"id"`
	DistributionZoneName string `json:"distributionzonename"`
}

//AddDistributionZone struct
type AddDistributionZone struct {
	DistributionZoneName string `json:"distributionzonename"`
}

//DistributionZone_count  struct
type DistributionZone_count struct {
	Values []DistributionZone `json:"values"`
	Count  int                `json:"count"`
	Auth   Auth               `json:"auth"`
}
