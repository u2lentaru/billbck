package models

//Result struct
type Result struct {
	Id         int    `json:"id"`
	ResultName string `json:"resultname"`
}

//AddResult struct
type AddResult struct {
	ResultName string `json:"resultname"`
}

//Result_count  struct
type Result_count struct {
	Values []Result `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}
