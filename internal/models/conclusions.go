package models

//Conclusion struct
type Conclusion struct {
	Id             int    `json:"id"`
	ConclusionName string `json:"conclusionname"`
}

//AddConclusion struct
type AddConclusion struct {
	ConclusionName string `json:"conclusionname"`
}

//Conclusion_count  struct
type Conclusion_count struct {
	Values []Conclusion `json:"values"`
	Count  int          `json:"count"`
	Auth   Auth         `json:"auth"`
}
