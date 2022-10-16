package models

//Balance struct
type Balance struct {
	Id         int    `json:"id"`
	PId        int    `json:"pid"`
	BName      string `json:"bname"`
	BTypeId    int    `json:"btypeid"`
	BTypeName  string `json:"btypename"`
	ChildCount int    `json:"childcount"`
	ReqId      string `json:"reqid"`
}

//Balance_count  struct
type Balance_count struct {
	Values []Balance `json:"values"`
	Count  int       `json:"count"`
}

//BalanceTab struct
type BalanceTab struct {
	Id        int     `json:"id"`
	PId       int     `json:"pid"`
	BName     string  `json:"bname"`
	BTypeId   int     `json:"btypeid"`
	BTypeName string  `json:"btypename"`
	Sum       float64 `json:"sum"`
	ReqId     string  `json:"reqid"`
}

//BalanceTab_sum  struct
type BalanceTab_sum struct {
	Values []BalanceTab `json:"values"`
	InSum  float64      `json:"insum"`
	OutSum float64      `json:"outsum"`
	Count  int          `json:"count"`
}
