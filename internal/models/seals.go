package models

//Seal struct
type Seal struct {
	Id           int        `json:"id"`
	PacketNumber *string    `json:"packetnumber"`
	Area         Area       `json:"area"`
	Staff        Staff      `json:"staff"`
	SealType     SealType   `json:"sealtype"`
	SealColour   SealColour `json:"sealcolour"`
	SealStatus   SealStatus `json:"sealstatus"`
	IssueDate    *string    `json:"issuedate"`
	ReportDate   *string    `json:"reportdate"`
}

//AddSeal struct
type AddSeal struct {
	PacketNumber *string    `json:"packetnumber"`
	Area         Area       `json:"area"`
	Staff        Staff      `json:"staff"`
	SealType     SealType   `json:"sealtype"`
	SealColour   SealColour `json:"sealcolour"`
	SealStatus   SealStatus `json:"sealstatus"`
	IssueDate    *string    `json:"issuedate"`
	ReportDate   *string    `json:"reportdate"`
}

//Seal_count  struct
type Seal_count struct {
	Values []Seal `json:"values"`
	Count  int    `json:"count"`
	Auth   Auth   `json:"auth"`
}
