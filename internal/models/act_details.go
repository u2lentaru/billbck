package models

//ActDetail struct
type ActDetail struct {
	Id               int          `json:"id"`
	Act              Act          `json:"act"`
	ActDetailDate    string       `json:"actdetaildate"`
	PuId             *int         `json:"puid"`
	PuType           PuType       `json:"putype"`
	PuNumber         *string      `json:"punumber"`
	InstallDate      *string      `json:"installdate"`
	InitialValue     *int         `json:"initialvalue"`
	CheckInterval    *int         `json:"checkinterval"`
	DevStopped       *bool        `json:"devstopped"`
	Startdate        *string      `json:"startdate"`
	Enddate          *string      `json:"enddate"`
	Pid              *int         `json:"pid"`
	Seal             Seal         `json:"seal"`
	SealNumber       *string      `json:"sealnumber"`
	SealDate         *string      `json:"sealdate"`
	AdPuValue        *float32     `json:"adpuvalue"`
	Conclusion       Conclusion   `json:"conclusion"`
	ConclusionNumber *string      `json:"conclusionnumber"`
	ShutdownType     ShutdownType `json:"shutdowntype"`
	Reason           Reason       `json:"reason"`
	Violation        Violation    `json:"violation"`
	Customer         *string      `json:"customer"`
	CustomerPhone    *string      `json:"customerphone"`
	CustomerPos      *string      `json:"customerpos"`
	Notes            *string      `json:"notes"`
}

//AddActDetail struct
type AddActDetail struct {
	Act              Act          `json:"act"`
	ActDetailDate    string       `json:"actdetaildate"`
	PuId             *int         `json:"puid"`
	PuType           PuType       `json:"putype"`
	PuNumber         *string      `json:"punumber"`
	InstallDate      *string      `json:"installdate"`
	InitialValue     *int         `json:"initialvalue"`
	CheckInterval    *int         `json:"checkinterval"`
	DevStopped       *bool        `json:"devstopped"`
	Startdate        *string      `json:"startdate"`
	Enddate          *string      `json:"enddate"`
	Pid              *int         `json:"pid"`
	Seal             Seal         `json:"seal"`
	SealNumber       *string      `json:"sealnumber"`
	SealDate         *string      `json:"sealdate"`
	AdPuValue        *float32     `json:"adpuvalue"`
	Conclusion       Conclusion   `json:"conclusion"`
	ConclusionNumber *string      `json:"conclusionnumber"`
	ShutdownType     ShutdownType `json:"shutdowntype"`
	Reason           Reason       `json:"reason"`
	Violation        Violation    `json:"violation"`
	Customer         *string      `json:"customer"`
	CustomerPhone    *string      `json:"customerphone"`
	CustomerPos      *string      `json:"customerpos"`
	Notes            *string      `json:"notes"`
}

//ActDetail_count  struct
type ActDetail_count struct {
	Values []ActDetail `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
