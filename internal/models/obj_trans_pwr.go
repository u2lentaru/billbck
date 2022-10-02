package models

//ObjTransPwr struct
type ObjTransPwr struct {
	Id        int      `json:"id"`
	ObjId     int      `json:"objid"`
	ObjTypeId int      `json:"objtypeid"`
	ObjName   string   `json:"objname"`
	TransPwr  TransPwr `json:"transpwr"`
	Startdate string   `json:"startdate"`
	Enddate   *string  `json:"enddate"`
}

//AddObjTransPwr struct
type AddObjTransPwr struct {
	ObjId     int      `json:"objid"`
	ObjTypeId int      `json:"objtypeid"`
	ObjName   string   `json:"objname"`
	TransPwr  TransPwr `json:"transpwr"`
	Startdate string   `json:"startdate"`
	Enddate   *string  `json:"enddate"`
}

//ObjTransPwr_count  struct
type ObjTransPwr_count struct {
	Values []ObjTransPwr `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
