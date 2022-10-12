package models

//ShutdownType struct
type ShutdownType struct {
	Id               int    `json:"id"`
	ShutdownTypeName string `json:"shutdowntypename"`
}

//AddShutdownType struct
type AddShutdownType struct {
	ShutdownTypeName string `json:"shutdowntypename"`
}

//ShutdownType_count  struct
type ShutdownType_count struct {
	Values []ShutdownType `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}
