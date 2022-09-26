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

//FormType struct
type FormType struct {
	Id            int    `json:"id"`
	FormTypeName  string `json:"formtypename"`
	FormTypeDescr string `json:"formtypedescr"`
}

//AddFormType struct
type AddFormType struct {
	FormTypeName  string `json:"formtypename"`
	FormTypeDescr string `json:"formtypedescr"`
}

//FormType_count  struct
type FormType_count struct {
	Values []FormType `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}

//Json_ids  struct
type Json_ids struct {
	Ids []int `json:"ids"`
}

//SubType struct
type SubType struct {
	Id           int    `json:"id"`
	SubTypeName  string `json:"subtypename"`
	SubTypeDescr string `json:"subtypedescr"`
}

//SubType_count  struct
type SubType_count struct {
	Values []SubType `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}

//AddSubType struct
type AddSubType struct {
	SubTypeName  string `json:"subtypename"`
	SubTypeDescr string `json:"subtypedescr"`
}

//Position struct
type Position struct {
	Id           int    `json:"id"`
	PositionName string `json:"positionname"`
}

//AddPosition struct
type AddPosition struct {
	PositionName string `json:"positionname"`
}

//Position_count  struct
type Position_count struct {
	Values []Position `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}

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

//AddForm struct
type AddForm struct {
	Form string `json:"form"`
}

//Vals  struct
type Vals struct {
	Vals []string `json:"vals"`
}

//OrgInfo struct
type OrgInfo struct {
	Id          int    `json:"id"`
	OIName      string `json:"oiname"`
	OIBin       string `json:"oibin"`
	OIAddr      string `json:"oiaddr"`
	OIBank      Bank   `json:"oibank"`
	OIAccNumber string `json:"oiaccnumber"`
	OIFName     string `json:"oifname"`
}

//AddOrgInfo struct
type AddOrgInfo struct {
	OIName      string `json:"oiname"`
	OIBin       string `json:"oibin"`
	OIAddr      string `json:"oiaddr"`
	OIBank      Bank   `json:"oibank"`
	OIAccNumber string `json:"oiaccnumber"`
	OIFName     string `json:"oifname"`
}

//OrgInfo_count  struct
type OrgInfo_count struct {
	Values []OrgInfo `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}

//SB struct
type SB struct {
	Id     int    `json:"id"`
	SBName string `json:"subname"`
}

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

//Ksk struct
type Ksk struct {
	Id         int    `json:"id"`
	KskName    string `json:"kskname"`
	KskAddress string `json:"kskaddress"`
	KskHead    string `json:"kskhead"`
	KskPhone   string `json:"kskphone"`
}

//AddKsk struct
type AddKsk struct {
	KskName    string `json:"kskname"`
	KskAddress string `json:"kskaddress"`
	KskHead    string `json:"kskhead"`
	KskPhone   string `json:"kskphone"`
}

//Ksk_count  struct
type Ksk_count struct {
	Values []Ksk `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}

//Sector struct
type Sector struct {
	Id         int    `json:"id"`
	SectorName string `json:"sectorname"`
}

//AddSector struct
type AddSector struct {
	SectorName string `json:"sectorname"`
}

//Sector_count  struct
type Sector_count struct {
	Values []Sector `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}

//InputType struct
type InputType struct {
	Id            int    `json:"id"`
	InputTypeName string `json:"inputtypename"`
}

//AddInputType struct
type AddInputType struct {
	InputTypeName string `json:"inputtypename"`
}

//InputType_count  struct
type InputType_count struct {
	Values []InputType `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}

//Reliability struct
type Reliability struct {
	Id              int    `json:"id"`
	ReliabilityName string `json:"reliabilityname"`
}

//AddReliability struct
type AddReliability struct {
	ReliabilityName string `json:"reliabilityname"`
}

//Reliability_count  struct
type Reliability_count struct {
	Values []Reliability `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
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

//Eso struct
type Eso struct {
	Id      int    `json:"id"`
	EsoName string `json:"esoname"`
}

//AddEso struct
type AddEso struct {
	EsoName string `json:"esoname"`
}

//Eso_count  struct
type Eso_count struct {
	Values []Eso `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}

//ObjType struct
type ObjType struct {
	Id          int    `json:"id"`
	ObjTypeName string `json:"objtypename"`
}

//AddObjType struct
type AddObjType struct {
	ObjTypeName string `json:"objtypename"`
}

//ObjType_count  struct
type ObjType_count struct {
	Values []ObjType `json:"values"`
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

//PuType struct
type PuType struct {
	Id         int    `json:"id"`
	PuTypeName string `json:"putypename"`
}

//AddPuType struct
type AddPuType struct {
	PuTypeName string `json:"putypename"`
}

//PuType_count  struct
type PuType_count struct {
	Values []PuType `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}

//PaymentType struct
type PaymentType struct {
	Id              int    `json:"id"`
	PaymentTypeName string `json:"paymenttypename"`
}

//AddPaymentType struct
type AddPaymentType struct {
	PaymentTypeName string `json:"paymenttypename"`
}

//PaymentType_count  struct
type PaymentType_count struct {
	Values []PaymentType `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}

//TariffGroup struct
type TariffGroup struct {
	Id              int    `json:"id"`
	TariffGroupName string `json:"tariffgroupname"`
}

//AddTariffGroup struct
type AddTariffGroup struct {
	TariffGroupName string `json:"tariffgroupname"`
}

//TariffGroup_count  struct
type TariffGroup_count struct {
	Values []TariffGroup `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}

//Tariff struct
type Tariff struct {
	Id          int         `json:"id"`
	TariffName  string      `json:"tariffname"`
	TariffGroup TariffGroup `json:"tariffgroup"`
	Norma       float32     `json:"norma"`
	Tariff      float32     `json:"tariff"`
	Startdate   string      `json:"startdate"`
	//*string correct scan null data!
	Enddate *string `json:"enddate"`
}

//AddTariff struct
type AddTariff struct {
	TariffName  string      `json:"tariffname"`
	TariffGroup TariffGroup `json:"tariffgroup"`
	Norma       float32     `json:"norma"`
	Tariff      float32     `json:"tariff"`
	Startdate   string      `json:"startdate"`
}

//Tariff_count  struct
type Tariff_count struct {
	Values []Tariff `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}

//Rp struct
type Rp struct {
	Id             int     `json:"id"`
	RpName         string  `json:"tguname"`
	InvNumber      string  `json:"invnumber"`
	InputVoltage   Voltage `json:"inputvoltage"`
	OutputVoltage1 Voltage `json:"outputvoltage1"`
	OutputVoltage2 Voltage `json:"outputvoltage2"`
	Tp             Tp      `json:"tp"`
}

//AddRp struct
type AddRp struct {
	RpName         string  `json:"tguname"`
	InvNumber      string  `json:"invnumber"`
	InputVoltage   Voltage `json:"inputvoltage"`
	OutputVoltage1 Voltage `json:"outputvoltage1"`
	OutputVoltage2 Voltage `json:"outputvoltage2"`
	Tp             Tp      `json:"tp"`
}

//Rp_count  struct
type Rp_count struct {
	Values []Rp `json:"values"`
	Count  int  `json:"count"`
	Auth   Auth `json:"auth"`
}

//GRp struct
type GRp struct {
	Id      int    `json:"id"`
	GRpName string `json:"grpname"`
}

//AddGRp struct
type AddGRp struct {
	GRpName string `json:"grpname"`
}

//GRp_count  struct
type GRp_count struct {
	Values []GRp `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}

//Tp struct
type Tp struct {
	Id     int    `json:"id"`
	TpName string `json:"tpname"`
	GRp    GRp    `json:"grp"`
}

//AddTp struct
type AddTp struct {
	TpName string `json:"tpname"`
	GRp    GRp    `json:"grp"`
}

//Tp_count  struct
type Tp_count struct {
	Values []Tp `json:"values"`
	Count  int  `json:"count"`
	Auth   Auth `json:"auth"`
}

//LoginForm struct
type LoginForm struct {
	Form   string `json:"form"`
	Rights string `json:"rights"`
	UserId int    `json:"userid"`
}

//Street struct
type Street struct {
	Id         int     `json:"id"`
	StreetName string  `json:"streetname"`
	Created    string  `json:"created"`
	Closed     *string `json:"closed"`
	City       City    `json:"city"`
}

//AddStreet struct
type AddStreet struct {
	StreetName string `json:"streetname"`
	Created    string `json:"created"`
	City       City   `json:"city"`
}

//Street_count  struct
type Street_count struct {
	Values []Street `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}

//StreetClose struct
type StreetClose struct {
	Id    int    `json:"id"`
	Close string `json:"close"`
}

//UpdStreet struct
type UpdStreet struct {
	Id         int    `json:"id"`
	StreetName string `json:"streetname"`
	Created    string `json:"created"`
}

//House struct
type House struct {
	Id             int          `json:"id"`
	BuildingType   BuildingType `json:"buildingtype"`
	Street         Street       `json:"street"`
	HouseNumber    string       `json:"housenumber"`
	BuildingNumber *string      `json:"buildingnumber"`
	RP             Rp           `json:"tgu"`
	Area           Area         `json:"area"`
	Ksk            Ksk          `json:"ksk"`
	Sector         Sector       `json:"sector"`
	Connector      Connector    `json:"connector"`
	InputType      InputType    `json:"inputtype"`
	Reliability    Reliability  `json:"reliability"`
	Voltage        Voltage      `json:"voltage"`
	Notes          *string      `json:"notes"`
}

type AddHouse struct {
	BuildingType   BuildingType `json:"buildingtype"`
	Street         Street       `json:"street"`
	HouseNumber    string       `json:"housenumber"`
	BuildingNumber *string      `json:"buildingnumber"`
	RP             Rp           `json:"tgu"`
	Area           Area         `json:"area"`
	Ksk            Ksk          `json:"ksk"`
	Sector         Sector       `json:"sector"`
	Connector      Connector    `json:"connector"`
	InputType      InputType    `json:"inputtype"`
	Reliability    Reliability  `json:"reliability"`
	Voltage        Voltage      `json:"voltage"`
	Notes          *string      `json:"notes"`
}

//UpdHouse struct
type UpdHouse struct {
	Id             int    `json:"id"`
	BuildingTypeId int    `json:"buildingtypeid"`
	StreetId       int    `json:"streetid"`
	HouseNumber    string `json:"housenumber"`
	BuildingNumber string `json:"buildingnumber"`
	RPId           int    `json:"rpid"`
	AreaId         int    `json:"areaid"`
	KskId          int    `json:"kskid"`
	SectorId       int    `json:"sectorid"`
	ConnectorId    int    `json:"connectorid"`
	InputTypeId    int    `json:"inputtypeid"`
	ReliabilityId  int    `json:"reliabilityid"`
	VoltageId      int    `json:"voltageid"`
}

//House_count  struct
type House_count struct {
	Values []House `json:"values"`
	Count  int     `json:"count"`
	Auth   Auth    `json:"auth"`
}

//Object struct
type Object struct {
	Id              int             `json:"id"`
	ObjectName      string          `json:"objectname"`
	House           House           `json:"house"`
	FlatNumber      *string         `json:"flatnumber"`
	ObjType         ObjType         `json:"objtype"`
	RegQty          int             `json:"regqty"`
	Uzo             Uzo             `json:"uzo"`
	TariffGroup     TariffGroup     `json:"tariffgroup"`
	CalculationType CalculationType `json:"calculationtype"`
	ObjStatus       ObjStatus       `json:"objstatus"`
	Notes           *string         `json:"notes"`
	MffId           *int            `json:"mffid"`
}

//AddObject struct
type AddObject struct {
	ObjectName      string          `json:"objectname"`
	House           House           `json:"house"`
	FlatNumber      *string         `json:"flatnumber"`
	ObjType         ObjType         `json:"objtype"`
	RegQty          int             `json:"regqty"`
	Uzo             Uzo             `json:"uzo"`
	TariffGroup     TariffGroup     `json:"tariffgroup"`
	CalculationType CalculationType `json:"calculationtype"`
	ObjStatus       ObjStatus       `json:"objstatus"`
	Notes           *string         `json:"notes"`
	MffId           *int            `json:"mffid"`
}

//Object_count  struct
type Object_count struct {
	Values []Object `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}

//ObjContract struct
type ObjContract struct {
	Id        int      `json:"id"`
	Object    Object   `json:"object"`
	Contract  Contract `json:"contract"`
	ObjTypeId int      `json:"objtypeid"`
	Startdate string   `json:"startdate"`
	Enddate   *string  `json:"enddate"`
}

//AddObjContract struct
type AddObjContract struct {
	Object    Object   `json:"object"`
	Contract  Contract `json:"contract"`
	ObjTypeId int      `json:"objtypeid"`
	Startdate string   `json:"startdate"`
	Enddate   *string  `json:"enddate"`
}

//ObjContract_count  struct
type ObjContract_count struct {
	Values []ObjContract `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}

//IdClose struct
type IdClose struct {
	Id          int         `json:"id"`
	CloseDate   string      `json:"closedate"`
	ContractMot ContractMot `json:"contractmot"`
	MotNotes    string      `json:"motnotes"`
}

//Pu struct
type Pu struct {
	Id            int     `json:"id"`
	Object        Object  `json:"object"`
	PuObjectType  int     `json:"puobjecttype"`
	PuType        PuType  `json:"putype"`
	PuNumber      string  `json:"punumber"`
	InstallDate   string  `json:"installdate"`
	InitialValue  int     `json:"initialvalue"`
	CheckInterval int     `json:"checkinterval"`
	DevStopped    bool    `json:"devstopped"`
	Startdate     string  `json:"startdate"`
	Enddate       *string `json:"enddate"`
	Pid           *int    `json:"pid"`
}

//AddPu struct
type AddPu struct {
	Object        Object `json:"object"`
	PuObjectType  int    `json:"puobjecttype"`
	PuType        PuType `json:"putype"`
	PuNumber      string `json:"punumber"`
	InstallDate   string `json:"installdate"`
	InitialValue  int    `json:"initialvalue"`
	CheckInterval int    `json:"checkinterval"`
	DevStopped    bool   `json:"devstopped"`
	Startdate     string `json:"startdate"`
	Pid           *int   `json:"pid"`
}

//Pu_count  struct
type Pu_count struct {
	Values []Pu `json:"values"`
	Count  int  `json:"count"`
	Auth   Auth `json:"auth"`
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

//SubPu struct
type SubPu struct {
	Id    int `json:"id"`
	ParId int `json:"parid"`
	SubId int `json:"subid"`
}

//AddSubPu struct
type AddSubPu struct {
	ParId int `json:"parid"`
	SubId int `json:"subid"`
}

//SubPu_count  struct
type SubPu_count struct {
	Values []SubPu `json:"values"`
	Count  int     `json:"count"`
	Auth   Auth    `json:"auth"`
}

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

//TguType struct
type TguType struct {
	Id          int    `json:"id"`
	TguTypeName string `json:"tgutypename"`
}

//Tgu struct
type Tgu struct {
	Id             int     `json:"id"`
	PId            int     `json:"pid"`
	TguName        string  `json:"tguname"`
	TguType        TguType `json:"tgutype"`
	InvNumber      *string `json:"invnumber"`
	InputVoltage   Voltage `json:"inputvoltage"`
	OutputVoltage1 Voltage `json:"outputvoltage1"`
	OutputVoltage2 Voltage `json:"outputvoltage2"`
}

//AddTgu struct
type AddTgu struct {
	PId            int     `json:"pid"`
	TguName        string  `json:"tguname"`
	TguType        TguType `json:"tgutype"`
	InvNumber      string  `json:"invnumber"`
	InputVoltage   Voltage `json:"inputvoltage"`
	OutputVoltage1 Voltage `json:"outputvoltage1"`
	OutputVoltage2 Voltage `json:"outputvoltage2"`
}

//Tgu_count  struct
type Tgu_count struct {
	Values []Tgu `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}

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

//ShutdownType struct
type ShutdownType struct {
	Id               int    `json:"id"`
	ShutdownTypeName string `json:"shutdowntypename"`
}

//AddShutdownType struct
type AddShutdownType struct {
	ShutdownTypeName string `json:"shutdowntypename"`
}

//ShutdownType_count  struct
type ShutdownType_count struct {
	Values []ShutdownType `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}

//Reason struct
type Reason struct {
	Id         int    `json:"id"`
	ReasonName string `json:"reasonname"`
}

//AddReason struct
type AddReason struct {
	ReasonName string `json:"reasonname"`
}

//Reason_count  struct
type Reason_count struct {
	Values []Reason `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
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

//SealType struct
type SealType struct {
	Id           int    `json:"id"`
	SealTypeName string `json:"sealtypename"`
}

//AddSealType struct
type AddSealType struct {
	SealTypeName string `json:"sealtypename"`
}

//SealType_count  struct
type SealType_count struct {
	Values []SealType `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}

//SealColour struct
type SealColour struct {
	Id             int    `json:"id"`
	SealColourName string `json:"sealcolourname"`
}

//AddSealColour struct
type AddSealColour struct {
	SealColourName string `json:"sealcolourname"`
}

//SealColour_count  struct
type SealColour_count struct {
	Values []SealColour `json:"values"`
	Count  int          `json:"count"`
	Auth   Auth         `json:"auth"`
}

//SealStatus struct
type SealStatus struct {
	Id             int    `json:"id"`
	SealStatusName string `json:"sealstatusname"`
}

//AddSealStatus struct
type AddSealStatus struct {
	SealStatusName string `json:"sealstatusname"`
}

//SealStatus_count  struct
type SealStatus_count struct {
	Values []SealStatus `json:"values"`
	Count  int          `json:"count"`
	Auth   Auth         `json:"auth"`
}

//Payment struct
type Payment struct {
	Id           int         `json:"id"`
	PaymentDate  string      `json:"paymentdate"`
	Contract     Contract    `json:"contract"`
	Object       Object      `json:"object"`
	PaymentType  PaymentType `json:"paymenttype"`
	ChargeType   ChargeType  `json:"chargetype"`
	Cashdesk     Cashdesk    `json:"cashdesk"`
	BundleNumber int         `json:"bundlenumber"`
	Amount       float32     `json:"amount"`
}

//AddPayment struct
type AddPayment struct {
	PaymentDate  string      `json:"paymentdate"`
	Contract     Contract    `json:"contract"`
	Object       Object      `json:"object"`
	PaymentType  PaymentType `json:"paymenttype"`
	ChargeType   ChargeType  `json:"chargetype"`
	Cashdesk     Cashdesk    `json:"cashdesk"`
	BundleNumber int         `json:"bundlenumber"`
	Amount       float32     `json:"amount"`
}

//Payment_count  struct
type Payment_count struct {
	Values []Payment `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}

//Result struct
type Result struct {
	Id         int    `json:"id"`
	ResultName string `json:"resultname"`
}

//AddResult struct
type AddResult struct {
	ResultName string `json:"resultname"`
}

//Result_count  struct
type Result_count struct {
	Values []Result `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}

//ServiceType struct
type ServiceType struct {
	Id              int    `json:"id"`
	ServiceTypeName string `json:"servicetypename"`
}

//AddServiceType struct
type AddServiceType struct {
	ServiceTypeName string `json:"servicetypename"`
}

//ServiceType_count struct
type ServiceType_count struct {
	Values []ServiceType `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}

//RequestType struct
type RequestType struct {
	Id              int         `json:"id"`
	RequestTypeName string      `json:"requesttypename"`
	RequestKind     RequestKind `json:"requestkind"`
}

//AddRequestType struct
type AddRequestType struct {
	RequestTypeName string      `json:"requesttypename"`
	RequestKind     RequestKind `json:"requestkind"`
}

//RequestType_count struct
type RequestType_count struct {
	Values []RequestType `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}

//RequestKind struct
type RequestKind struct {
	Id              int    `json:"id"`
	RequestKindName string `json:"requestkindname"`
}

//AddRequestKind struct
type AddRequestKind struct {
	RequestKindName string `json:"requestkindname"`
}

//RequestKind_count struct
type RequestKind_count struct {
	Values []RequestKind `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}

//Request struct
type Request struct {
	Id            int         `json:"id"`
	RequestNumber string      `json:"requestnumber"`
	RequestDate   string      `json:"requestdate"`
	Contract      Contract    `json:"contract"`
	ServiceType   ServiceType `json:"servicetype"`
	RequestType   RequestType `json:"requesttype"`
	RequestKind   RequestKind `json:"requestkind"`
	ClaimType     ClaimType   `json:"claimtype"`
	TermDate      string      `json:"termdate"`
	Executive     string      `json:"executive"`
	Accept        string      `json:"accept"`
	Notes         *string     `json:"notes"`
	Result        Result      `json:"result"`
	Act           Act         `json:"act"`
	Object        Object      `json:"object"`
}

//AddRequest struct
type AddRequest struct {
	RequestNumber string      `json:"requestnumber"`
	RequestDate   string      `json:"requestdate"`
	Contract      Contract    `json:"contract"`
	ServiceType   ServiceType `json:"servicetype"`
	RequestType   RequestType `json:"requesttype"`
	RequestKind   RequestKind `json:"requestkind"`
	ClaimType     ClaimType   `json:"claimtype"`
	TermDate      string      `json:"termdate"`
	Executive     string      `json:"executive"`
	Accept        string      `json:"accept"`
	Notes         *string     `json:"notes"`
	Result        Result      `json:"result"`
	Act           Act         `json:"act"`
	Object        Object      `json:"object"`
}

//Request_count struct
type Request_count struct {
	Values []Request `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}

//ObjStatus struct
type ObjStatus struct {
	Id            int    `json:"id"`
	ObjStatusName string `json:"objstatusname"`
}

//AddObjStatus struct
type AddObjStatus struct {
	ObjStatusName string `json:"objstatusname"`
}

//ObjStatus_count  struct
type ObjStatus_count struct {
	Values []ObjStatus `json:"values"`
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

//Period struct
type Period struct {
	Id         int    `json:"id"`
	PeriodName string `json:"periodname"`
	Startdate  string `json:"startdate"`
	Enddate    string `json:"enddate"`
}

//AddPeriod struct
type AddPeriod struct {
	PeriodName string `json:"periodname"`
	Startdate  string `json:"startdate"`
	Enddate    string `json:"enddate"`
}

//Period_count  struct
type Period_count struct {
	Values []Period `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
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

//TransCurr struct
type TransCurr struct {
	Id            int       `json:"id"`
	TransCurrName string    `json:"transcurrname"`
	TransType     TransType `json:"transtype"`
	CheckDate     *string   `json:"checkdate"`
	NextCheckDate *string   `json:"nextcheckdate"`
	ProdDate      *string   `json:"proddate"`
	Serial1       *string   `json:"serial1"`
	Serial2       *string   `json:"serial2"`
	Serial3       *string   `json:"serial3"`
}

//AddTransCurr struct
type AddTransCurr struct {
	TransCurrName string    `json:"transcurrname"`
	TransType     TransType `json:"transtype"`
	CheckDate     *string   `json:"checkdate"`
	NextCheckDate *string   `json:"nextcheckdate"`
	ProdDate      *string   `json:"proddate"`
	Serial1       *string   `json:"serial1"`
	Serial2       *string   `json:"serial2"`
	Serial3       *string   `json:"serial3"`
}

//TransCurr_count  struct
type TransCurr_count struct {
	Values []TransCurr `json:"values"`
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

//ObjTransCurr struct
type ObjTransCurr struct {
	Id        int       `json:"id"`
	ObjId     int       `json:"objid"`
	ObjTypeId int       `json:"objtypeid"`
	ObjName   string    `json:"objname"`
	TransCurr TransCurr `json:"transcurr"`
	Startdate string    `json:"startdate"`
	Enddate   *string   `json:"enddate"`
}

//AddObjTransCurr struct
type AddObjTransCurr struct {
	ObjId     int       `json:"objid"`
	ObjTypeId int       `json:"objtypeid"`
	ObjName   string    `json:"objname"`
	TransCurr TransCurr `json:"transcurr"`
	Startdate string    `json:"startdate"`
	Enddate   *string   `json:"enddate"`
}

//ObjTransCurr_count  struct
type ObjTransCurr_count struct {
	Values []ObjTransCurr `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}

//ObjTransVolt struct
type ObjTransVolt struct {
	Id        int       `json:"id"`
	ObjId     int       `json:"objid"`
	ObjTypeId int       `json:"objtypeid"`
	ObjName   string    `json:"objname"`
	TransVolt TransVolt `json:"transvolt"`
	Startdate string    `json:"startdate"`
	Enddate   *string   `json:"enddate"`
}

//AddObjTransVolt struct
type AddObjTransVolt struct {
	ObjId     int       `json:"objid"`
	ObjTypeId int       `json:"objtypeid"`
	ObjName   string    `json:"objname"`
	TransVolt TransVolt `json:"transvolt"`
	Startdate string    `json:"startdate"`
	Enddate   *string   `json:"enddate"`
}

//ObjTransVolt_count  struct
type ObjTransVolt_count struct {
	Values []ObjTransVolt `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}

//ObjTransPwr struct
type ObjTransPwr struct {
	Id        int      `json:"id"`
	ObjId     int      `json:"objid"`
	ObjTypeId int      `json:"objtypeid"`
	ObjName   string   `json:"objname"`
	TransPwr  TransPwr `json:"transpwr"`
	Startdate string   `json:"startdate"`
	Enddate   *string  `json:"enddate"`
}

//AddObjTransPwr struct
type AddObjTransPwr struct {
	ObjId     int      `json:"objid"`
	ObjTypeId int      `json:"objtypeid"`
	ObjName   string   `json:"objname"`
	TransPwr  TransPwr `json:"transpwr"`
	Startdate string   `json:"startdate"`
	Enddate   *string  `json:"enddate"`
}

//ObjTransPwr_count  struct
type ObjTransPwr_count struct {
	Values []ObjTransPwr `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}

//ObjLine struct
type ObjLine struct {
	Id              int             `json:"id"`
	ObjId           int             `json:"objid"`
	ObjTypeId       int             `json:"objtypeid"`
	ObjName         string          `json:"objname"`
	CableResistance CableResistance `json:"cableresistance"`
	LineLength      float32         `json:"linelength"`
	Startdate       string          `json:"startdate"`
	Enddate         *string         `json:"enddate"`
}

//AddObjLine struct
type AddObjLine struct {
	ObjId           int             `json:"objid"`
	ObjTypeId       int             `json:"objtypeid"`
	ObjName         string          `json:"objname"`
	CableResistance CableResistance `json:"cableresistance"`
	LineLength      float32         `json:"linelength"`
	Startdate       string          `json:"startdate"`
	Enddate         *string         `json:"enddate"`
}

//ObjLine_count  struct
type ObjLine_count struct {
	Values []ObjLine `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
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

//PuValueType struct
type PuValueType struct {
	Id              int    `json:"id"`
	PuValueTypeName string `json:"puvaluetypename"`
}

//AddPuValueType struct
type AddPuValueType struct {
	PuValueTypeName string `json:"puvaluetypename"`
}

//PuValueType_count  struct
type PuValueType_count struct {
	Values []PuValueType `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}

//OrderType struct
type OrderType struct {
	Id            int    `json:"id"`
	OrderTypeName string `json:"ordertypename"`
}

//AddOrderType struct
type AddOrderType struct {
	OrderTypeName string `json:"ordertypename"`
}

//OrderType_count  struct
type OrderType_count struct {
	Values []OrderType `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
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
