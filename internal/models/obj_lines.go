package models

//ObjLine struct
type ObjLine struct {
	Id              int             `json:"id"`
	ObjId           int             `json:"objid"`
	ObjTypeId       int             `json:"objtypeid"`
	ObjName         string          `json:"objname"`
	CableResistance CableResistance `json:"cableresistance"`
	LineLength      float32         `json:"linelength"`
	Startdate       string          `json:"startdate"`
	Enddate         *string         `json:"enddate"`
}

//AddObjLine struct
type AddObjLine struct {
	ObjId           int             `json:"objid"`
	ObjTypeId       int             `json:"objtypeid"`
	ObjName         string          `json:"objname"`
	CableResistance CableResistance `json:"cableresistance"`
	LineLength      float32         `json:"linelength"`
	Startdate       string          `json:"startdate"`
	Enddate         *string         `json:"enddate"`
}

//ObjLine_count  struct
type ObjLine_count struct {
	Values []ObjLine `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}
