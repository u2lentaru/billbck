package models

//GRp struct
type GRp struct {
	Id      int    `json:"id"`
	GRpName string `json:"grpname"`
}

//AddGRp struct
type AddGRp struct {
	GRpName string `json:"grpname"`
}

//GRp_count  struct
type GRp_count struct {
	Values []GRp `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}
