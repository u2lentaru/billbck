package models

import (
	"database/sql"
)

type NullInt32 sql.NullInt32
type NullString sql.NullString

//Json_id  struct
type Json_id struct {
	Id int `json:"id"`
}

//Json_sum  struct
type Json_sum struct {
	Sum float32 `json:"sum"`
}

//Auth  struct
type Auth struct {
	Create   bool  `json:"create"`
	Read     bool  `json:"read"`
	Update   bool  `json:"update"`
	Delete   bool  `json:"delete"`
	ActTypes []int `json:"acttypes"`
}

//Json_ids  struct
type Json_ids struct {
	Ids []int `json:"ids"`
}

//AddForm struct
type AddForm struct {
	Form string `json:"form"`
}

//IdClose struct
type IdClose struct {
	Id          int         `json:"id"`
	CloseDate   string      `json:"closedate"`
	ContractMot ContractMot `json:"contractmot"`
	MotNotes    string      `json:"motnotes"`
}

//Lang struct
type Lang struct {
	Id        int    `json:"id"`
	LangName  string `json:"langname"`
	LangDescr string `json:"langdescr"`
}

//AddLang struct
type AddLang struct {
	LangName  string `json:"langname"`
	LangDescr string `json:"langdescr"`
}

//Lang_count  struct
type Lang_count struct {
	Values []Lang `json:"values"`
	Count  int    `json:"count"`
}
