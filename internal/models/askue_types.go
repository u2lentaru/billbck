package models

//AskueType struct
type AskueType struct {
	Id              int     `json:"id"`
	AskueTypeName   string  `json:"askuetypename"`
	StartLine       int     `json:"startline"`
	PuColumn        int     `json:"pucolumn"`
	ValueColumn     int     `json:"valuecolumn"`
	DateColumn      int     `json:"datecolumn"`
	DateColumnArray *string `json:"datecolumnarray"`
}

//AddAskueType struct
type AddAskueType struct {
	AskueTypeName   string  `json:"askuetypename"`
	StartLine       int     `json:"startline"`
	PuColumn        int     `json:"pucolumn"`
	ValueColumn     int     `json:"valuecolumn"`
	DateColumn      int     `json:"datecolumn"`
	DateColumnArray *string `json:"datecolumnarray"`
}

//AskueType_count  struct
type AskueType_count struct {
	Values []AskueType `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
