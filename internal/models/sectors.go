package models

//Sector struct
type Sector struct {
	Id         int    `json:"id"`
	SectorName string `json:"sectorname"`
}

//AddSector struct
type AddSector struct {
	SectorName string `json:"sectorname"`
}

//Sector_count  struct
type Sector_count struct {
	Values []Sector `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}
