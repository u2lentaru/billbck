package models

//Reliability struct
type Reliability struct {
	Id              int    `json:"id"`
	ReliabilityName string `json:"reliabilityname"`
}

//AddReliability struct
type AddReliability struct {
	ReliabilityName string `json:"reliabilityname"`
}

//Reliability_count  struct
type Reliability_count struct {
	Values []Reliability `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
