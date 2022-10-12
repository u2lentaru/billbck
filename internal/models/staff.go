package models

//Staff struct
type Staff struct {
	Id        int     `json:"id"`
	StaffName string  `json:"staffname"`
	OrgInfo   OrgInfo `json:"orginfo"`
	Phone     *string `json:"phone"`
	Notes     *string `json:"notes"`
}

//AddStaff struct
type AddStaff struct {
	StaffName string  `json:"staffname"`
	OrgInfo   OrgInfo `json:"orginfo"`
	Phone     *string `json:"phone"`
	Notes     *string `json:"notes"`
}

//Staff_count  struct
type Staff_count struct {
	Values []Staff `json:"values"`
	Count  int     `json:"count"`
	Auth   Auth    `json:"auth"`
}
