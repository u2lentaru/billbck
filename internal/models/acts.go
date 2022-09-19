package models

//Act struct
type Act struct {
	Id        int     `json:"id"`
	ActType   ActType `json:"acttype"`
	ActNumber string  `json:"actnumber"`
	ActDate   string  `json:"actdate"`
	Object    Object  `json:"object"`
	Staff     Staff   `json:"staff"`
	Notes     *string `json:"notes"`
	Activated *string `json:"activated"`
}

//AddAct struct
type AddAct struct {
	ActType   ActType `json:"acttype"`
	ActNumber string  `json:"actnumber"`
	ActDate   string  `json:"actdate"`
	Object    Object  `json:"object"`
	Staff     Staff   `json:"staff"`
	Notes     *string `json:"notes"`
	Activated *string `json:"activated"`
}

//Act_count  struct
type Act_count struct {
	Values []Act `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}
