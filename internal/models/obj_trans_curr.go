package models

//ObjTransCurr struct
type ObjTransCurr struct {
	Id        int       `json:"id"`
	ObjId     int       `json:"objid"`
	ObjTypeId int       `json:"objtypeid"`
	ObjName   string    `json:"objname"`
	TransCurr TransCurr `json:"transcurr"`
	Startdate string    `json:"startdate"`
	Enddate   *string   `json:"enddate"`
}

//AddObjTransCurr struct
type AddObjTransCurr struct {
	ObjId     int       `json:"objid"`
	ObjTypeId int       `json:"objtypeid"`
	ObjName   string    `json:"objname"`
	TransCurr TransCurr `json:"transcurr"`
	Startdate string    `json:"startdate"`
	Enddate   *string   `json:"enddate"`
}

//ObjTransCurr_count  struct
type ObjTransCurr_count struct {
	Values []ObjTransCurr `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}
