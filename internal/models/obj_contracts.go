package models

//ObjContract struct
type ObjContract struct {
	Id        int      `json:"id"`
	Object    Object   `json:"object"`
	Contract  Contract `json:"contract"`
	ObjTypeId int      `json:"objtypeid"`
	Startdate string   `json:"startdate"`
	Enddate   *string  `json:"enddate"`
}

//AddObjContract struct
type AddObjContract struct {
	Object    Object   `json:"object"`
	Contract  Contract `json:"contract"`
	ObjTypeId int      `json:"objtypeid"`
	Startdate string   `json:"startdate"`
	Enddate   *string  `json:"enddate"`
}

//ObjContract_count  struct
type ObjContract_count struct {
	Values []ObjContract `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
