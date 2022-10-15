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

//Vals  struct
type Vals struct {
	Vals []string `json:"vals"`
}

//SB struct
type SB struct {
	Id     int    `json:"id"`
	SBName string `json:"subname"`
}

//Voltage struct
type Voltage struct {
	Id           int    `json:"id"`
	VoltageName  string `json:"voltagename"`
	VoltageValue int    `json:"voltagevalue"`
}

//AddVoltage struct
type AddVoltage struct {
	VoltageName  string `json:"voltagename"`
	VoltageValue int    `json:"voltagevalue"`
}

//Voltage_count  struct
type Voltage_count struct {
	Values []Voltage `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}

//Uzo struct
type Uzo struct {
	Id       int    `json:"id"`
	UzoName  string `json:"uzoname"`
	UzoValue int    `json:"uzovalue"`
}

//AddUzo struct
type AddUzo struct {
	UzoName  string `json:"uzoname"`
	UzoValue int    `json:"uzovalue"`
}

//Uzo_count  struct
type Uzo_count struct {
	Values []Uzo `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}

//LoginForm struct
type LoginForm struct {
	Form   string `json:"form"`
	Rights string `json:"rights"`
	UserId int    `json:"userid"`
}

//IdClose struct
type IdClose struct {
	Id          int         `json:"id"`
	CloseDate   string      `json:"closedate"`
	ContractMot ContractMot `json:"contractmot"`
	MotNotes    string      `json:"motnotes"`
}

//PuValue struct
type PuValue struct {
	Id        int    `json:"id"`
	PuId      int    `json:"puid"`
	ValueDate string `json:"valuedate"`
	// PuValue   float32 `json:"puvalue,string"` //float32->string, struct->json
	PuValue string `json:"puvalue"`
}

//PuValue_count  struct
type PuValue_count struct {
	Values []PuValue `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}

//Balance struct
type Balance struct {
	Id         int    `json:"id"`
	PId        int    `json:"pid"`
	BName      string `json:"bname"`
	BTypeId    int    `json:"btypeid"`
	BTypeName  string `json:"btypename"`
	ChildCount int    `json:"childcount"`
	ReqId      string `json:"reqid"`
}

//Balance_count  struct
type Balance_count struct {
	Values []Balance `json:"values"`
	Count  int       `json:"count"`
}

//BalanceTab struct
type BalanceTab struct {
	Id        int     `json:"id"`
	PId       int     `json:"pid"`
	BName     string  `json:"bname"`
	BTypeId   int     `json:"btypeid"`
	BTypeName string  `json:"btypename"`
	Sum       float64 `json:"sum"`
	ReqId     string  `json:"reqid"`
}

//BalanceTab_sum  struct
type BalanceTab_sum struct {
	Values []BalanceTab `json:"values"`
	InSum  float64      `json:"insum"`
	OutSum float64      `json:"outsum"`
	Count  int          `json:"count"`
}

//Violation struct
type Violation struct {
	Id            int    `json:"id"`
	ViolationName string `json:"violationname"`
}

//AddViolation struct
type AddViolation struct {
	ViolationName string `json:"violationname"`
}

//Violation_count  struct
type Violation_count struct {
	Values []Violation `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
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

//TransType struct
type TransType struct {
	Id            int    `json:"id"`
	TransTypeName string `json:"transtypename"`
	Ratio         int    `json:"ratio"`
	Class         int    `json:"class"`
	MaxCurr       int    `json:"maxcurr"`
	NomCurr       int    `json:"nomcurr"`
}

//AddTransType struct
type AddTransType struct {
	TransTypeName string `json:"transtypename"`
	Ratio         int    `json:"ratio"`
	Class         int    `json:"class"`
	MaxCurr       int    `json:"maxcurr"`
	NomCurr       int    `json:"nomcurr"`
}

//TransType_count  struct
type TransType_count struct {
	Values []TransType `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}

//TransVolt struct
type TransVolt struct {
	Id            int       `json:"id"`
	TransVoltName string    `json:"transvoltname"`
	TransType     TransType `json:"transtype"`
	CheckDate     *string   `json:"checkdate"`
	NextCheckDate *string   `json:"nextcheckdate"`
	ProdDate      *string   `json:"proddate"`
	Serial1       *string   `json:"serial1"`
	Serial2       *string   `json:"serial2"`
	Serial3       *string   `json:"serial3"`
}

//AddTransVolt struct
type AddTransVolt struct {
	TransVoltName string    `json:"transvoltname"`
	TransType     TransType `json:"transtype"`
	CheckDate     *string   `json:"checkdate"`
	NextCheckDate *string   `json:"nextcheckdate"`
	ProdDate      *string   `json:"proddate"`
	Serial1       *string   `json:"serial1"`
	Serial2       *string   `json:"serial2"`
	Serial3       *string   `json:"serial3"`
}

//TransVolt_count  struct
type TransVolt_count struct {
	Values []TransVolt `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}

//TransPwrType struct
type TransPwrType struct {
	Id                int     `json:"id"`
	TransPwrTypeName  string  `json:"transpwrtypename"`
	ShortCircuitPower float32 `json:"shortcircuitpower"`
	IdlingLossPower   float32 `json:"idlinglosspower"`
	NominalPower      int     `json:"nominalpower"`
}

//AddTransPwrType struct
type AddTransPwrType struct {
	TransPwrTypeName  string  `json:"transpwrtypename"`
	ShortCircuitPower float32 `json:"shortcircuitpower"`
	IdlingLossPower   float32 `json:"idlinglosspower"`
	NominalPower      int     `json:"nominalpower"`
}

//TransPwrType_count  struct
type TransPwrType_count struct {
	Values []TransPwrType `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}

//TransPwr struct
type TransPwr struct {
	Id           int          `json:"id"`
	TransPwrName string       `json:"transpwrname"`
	TransPwrType TransPwrType `json:"tranpwrstype"`
}

//AddTransPwr struct
type AddTransPwr struct {
	TransPwrName string       `json:"transpwrname"`
	TransPwrType TransPwrType `json:"tranpwrstype"`
}

//TransPwr_count  struct
type TransPwr_count struct {
	Values []TransPwr `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}

//PuValueAskue struct
type PuValueAskue struct {
	Valid     bool    `json:"valid"`
	PuId      int     `json:"puid"`
	PuNumber  string  `json:"punumber"`
	ValueDate string  `json:"valuedate"`
	PuValue   float32 `json:"puvalue"`
	Notes     string  `json:"notes"`
}

//PuValueAskue_count  struct
type PuValueAskue_count struct {
	Values []PuValueAskue `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}

//AskueFile struct
type AskueFile struct {
	AskueFile []byte `json:"askuefile"`
	Sheet     string `json:"sheet"`
	AskueType int    `json:"askuetype"`
}

//AskueLoadRes struct
type AskueLoadRes struct {
	Processed int `json:"processed"`
	Denied    int `json:"denied"`
}

// //DistributionZone struct
// type DistributionZone struct {
// 	Id                   int    `json:"id"`
// 	DistributionZoneName string `json:"distributionzonename"`
// }
// //AddDistributionZone struct
// type AddDistributionZone struct {
// 	DistributionZoneName string `json:"distributionzonename"`
// }
// //DistributionZone_count  struct
// type DistributionZone_count struct {
// 	Values []DistributionZone `json:"values"`
// 	Count  int                `json:"count"`
// 	Auth   Auth               `json:"auth"`
// }
