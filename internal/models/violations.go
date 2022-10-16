package models

//Violation struct
type Violation struct {
	Id            int    `json:"id"`
	ViolationName string `json:"violationname"`
}

//AddViolation struct
type AddViolation struct {
	ViolationName string `json:"violationname"`
}

//Violation_count  struct
type Violation_count struct {
	Values []Violation `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
