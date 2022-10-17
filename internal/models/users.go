package models

//User struct
type User struct {
	Id           int      `json:"id"`
	UserName     string   `json:"username"`
	OrgInfo      OrgInfo  `json:"orginfo"`
	Lang         Lang     `json:"lang"`
	ChangePass   bool     `json:"changepass"`
	Position     Position `json:"position"`
	UserFullName string   `json:"userfullname"`
	Created      string   `json:"created"`
	Closed       *string  `json:"closed"`
}

//AddUser struct
type AddUser struct {
	UserName     string   `json:"username"`
	OrgInfo      OrgInfo  `json:"orginfo"`
	Lang         Lang     `json:"lang"`
	ChangePass   bool     `json:"changepass"`
	Position     Position `json:"position"`
	UserFullName string   `json:"userfullname"`
	Created      string   `json:"created"`
	Closed       *string  `json:"closed"`
}

//User_count  struct
type User_count struct {
	Values []User `json:"values"`
	Count  int    `json:"count"`
	Auth   Auth   `json:"auth"`
}
