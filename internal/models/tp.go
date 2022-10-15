package models

//Tp struct
type Tp struct {
	Id     int    `json:"id"`
	TpName string `json:"tpname"`
	GRp    GRp    `json:"grp"`
}

//AddTp struct
type AddTp struct {
	TpName string `json:"tpname"`
	GRp    GRp    `json:"grp"`
}

//Tp_count  struct
type Tp_count struct {
	Values []Tp `json:"values"`
	Count  int  `json:"count"`
	Auth   Auth `json:"auth"`
}
