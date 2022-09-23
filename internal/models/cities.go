package models

//City struct
type City struct {
	Id       int    `json:"id"`
	CityName string `json:"cityname"`
}

//AddCity struct
type AddCity struct {
	CityName string `json:"cityname"`
}

//City_count  struct
type City_count struct {
	Values []City `json:"values"`
	Count  int    `json:"count"`
	Auth   Auth   `json:"auth"`
}
