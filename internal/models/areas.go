package models

//Area struct
type Area struct {
	Id         int    `json:"id"`
	AreaNumber string `json:"areanumber"`
	AreaName   string `json:"areaname"`
}

//AddArea struct
type AddArea struct {
	AreaNumber string `json:"areanumber"`
	AreaName   string `json:"areaname"`
}

//Area_count  struct
type Area_count struct {
	Values []Area `json:"values"`
	Count  int    `json:"count"`
	Auth   Auth   `json:"auth"`
}
