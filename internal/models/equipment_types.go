package models

//EquipmentType struct
type EquipmentType struct {
	Id                 int     `json:"id"`
	EquipmentTypeName  string  `json:"equipmenttypename"`
	EquipmentTypePower float32 `json:"equipmenttypepower"`
}

//AddEquipmentType struct
type AddEquipmentType struct {
	EquipmentTypeName  string  `json:"equipmenttypename"`
	EquipmentTypePower float32 `json:"equipmenttypepower"`
}

//EquipmentType_count  struct
type EquipmentType_count struct {
	Values []EquipmentType `json:"values"`
	Count  int             `json:"count"`
	Auth   Auth            `json:"auth"`
}
