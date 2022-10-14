package models

//Subject struct
type Subject struct {
	SubId        int      `json:"id"`
	SubType      SubType  `json:"subtype"`
	SubPhys      bool     `json:"subphys"`
	SubDescr     string   `json:"subdescr"`
	SubName      string   `json:"subname"`
	SubBin       string   `json:"subbin"`
	SubHeadPos   Position `json:"subheadpos"`
	SubHeadName  *string  `json:"subheadname"`
	SubAccPos    Position `json:"subaccpos"`
	SubAccName   *string  `json:"subaccname"`
	SubAddr      string   `json:"subaddr"`
	SubPhone     string   `json:"subphone"`
	SubStart     string   `json:"substart"`
	SubAccNumber string   `json:"subaccnumber"`
	Job          *string  `json:"job"`
	Email        *string  `json:"email"`
	MobPhone     *string  `json:"mobphone"`
	JobPhone     *string  `json:"jobphone"`
	Notes        *string  `json:"notes"`
}

//Subject_count  struct
type Subject_count struct {
	Values []Subject `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}

//AddSubject struct
type AddSubject struct {
	SubType      SubType  `json:"subtype"`
	SubPhys      bool     `json:"subphys"`
	SubDescr     string   `json:"subdescr"`
	SubName      string   `json:"subname"`
	SubBin       string   `json:"subbin"`
	SubHeadPos   Position `json:"subheadpos"`
	SubHeadName  *string  `json:"subheadname"`
	SubAccPos    Position `json:"subaccpos"`
	SubAccName   *string  `json:"subaccname"`
	SubAddr      string   `json:"subaddr"`
	SubPhone     string   `json:"subphone"`
	SubStart     string   `json:"substart"`
	SubAccNumber string   `json:"subaccnumber"`
	Job          *string  `json:"job"`
	Email        *string  `json:"email"`
	MobPhone     *string  `json:"mobphone"`
	JobPhone     *string  `json:"jobphone"`
	Notes        *string  `json:"notes"`
}

//SubjectClose struct
type SubjectClose struct {
	SubId    int    `json:"id"`
	SubClose string `json:"subclose"`
}
