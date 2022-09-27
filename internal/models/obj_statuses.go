package models

//ObjStatus struct
type ObjStatus struct {
	Id            int    `json:"id"`
	ObjStatusName string `json:"objstatusname"`
}

//AddObjStatus struct
type AddObjStatus struct {
	ObjStatusName string `json:"objstatusname"`
}

//ObjStatus_count  struct
type ObjStatus_count struct {
	Values []ObjStatus `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
