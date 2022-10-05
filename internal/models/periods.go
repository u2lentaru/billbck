package models

//Period struct
type Period struct {
	Id         int    `json:"id"`
	PeriodName string `json:"periodname"`
	Startdate  string `json:"startdate"`
	Enddate    string `json:"enddate"`
}

//AddPeriod struct
type AddPeriod struct {
	PeriodName string `json:"periodname"`
	Startdate  string `json:"startdate"`
	Enddate    string `json:"enddate"`
}

//Period_count  struct
type Period_count struct {
	Values []Period `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}
