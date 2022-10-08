package models

//RequestKind struct
type RequestKind struct {
	Id              int    `json:"id"`
	RequestKindName string `json:"requestkindname"`
}

//AddRequestKind struct
type AddRequestKind struct {
	RequestKindName string `json:"requestkindname"`
}

//RequestKind_count struct
type RequestKind_count struct {
	Values []RequestKind `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
