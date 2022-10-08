package models

//Pu struct
type Pu struct {
	Id            int     `json:"id"`
	Object        Object  `json:"object"`
	PuObjectType  int     `json:"puobjecttype"`
	PuType        PuType  `json:"putype"`
	PuNumber      string  `json:"punumber"`
	InstallDate   string  `json:"installdate"`
	InitialValue  int     `json:"initialvalue"`
	CheckInterval int     `json:"checkinterval"`
	DevStopped    bool    `json:"devstopped"`
	Startdate     string  `json:"startdate"`
	Enddate       *string `json:"enddate"`
	Pid           *int    `json:"pid"`
}

//AddPu struct
type AddPu struct {
	Object        Object `json:"object"`
	PuObjectType  int    `json:"puobjecttype"`
	PuType        PuType `json:"putype"`
	PuNumber      string `json:"punumber"`
	InstallDate   string `json:"installdate"`
	InitialValue  int    `json:"initialvalue"`
	CheckInterval int    `json:"checkinterval"`
	DevStopped    bool   `json:"devstopped"`
	Startdate     string `json:"startdate"`
	Pid           *int   `json:"pid"`
}

//Pu_count  struct
type Pu_count struct {
	Values []Pu `json:"values"`
	Count  int  `json:"count"`
	Auth   Auth `json:"auth"`
}
