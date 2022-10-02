package models

//ObjTransVolt struct
type ObjTransVolt struct {
	Id        int       `json:"id"`
	ObjId     int       `json:"objid"`
	ObjTypeId int       `json:"objtypeid"`
	ObjName   string    `json:"objname"`
	TransVolt TransVolt `json:"transvolt"`
	Startdate string    `json:"startdate"`
	Enddate   *string   `json:"enddate"`
}

//AddObjTransVolt struct
type AddObjTransVolt struct {
	ObjId     int       `json:"objid"`
	ObjTypeId int       `json:"objtypeid"`
	ObjName   string    `json:"objname"`
	TransVolt TransVolt `json:"transvolt"`
	Startdate string    `json:"startdate"`
	Enddate   *string   `json:"enddate"`
}

//ObjTransVolt_count  struct
type ObjTransVolt_count struct {
	Values []ObjTransVolt `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}
