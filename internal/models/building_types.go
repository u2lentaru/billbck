package models

//BuildingType struct
type BuildingType struct {
	Id               int    `json:"id"`
	BuildingTypeName string `json:"buildingtypename"`
}

//AddBuildingType struct
type AddBuildingType struct {
	BuildingTypeName string `json:"buildingtypename"`
}

//BuildingType_count  struct
type BuildingType_count struct {
	Values []BuildingType `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}
