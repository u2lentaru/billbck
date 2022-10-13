package models

//SubBank struct
type SubBank struct {
	Id        int    `json:"id"`
	Sub       SB     `json:"subj"`
	Bank      Bank   `json:"bank"`
	AccNumber string `json:"accnumber"`
	Active    bool   `json:"active"`
}

//AddSubBank struct
type AddSubBank struct {
	Sub       SB     `json:"subj"`
	Bank      Bank   `json:"bank"`
	AccNumber string `json:"accnumber"`
	Active    bool   `json:"active"`
}

//SubBank_count  struct
type SubBank_count struct {
	Values []SubBank `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}
