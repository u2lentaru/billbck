package models

//Charge struct
type Charge struct {
	Id         int        `json:"id"`
	ChargeDate string     `json:"chargedate"`
	Contract   Contract   `json:"contract"`
	Object     Object     `json:"object"`
	ObjTypeId  int        `json:"objtypeid"`
	Pu         Pu         `json:"pu"`
	ChargeType ChargeType `json:"chargetype"`
	Qty        float32    `json:"qty"`
	TransLoss  float32    `json:"transloss"`
	Lineloss   float32    `json:"lineloss"`
	Startdate  string     `json:"startdate"`
	Enddate    string     `json:"enddate"`
}

//AddCharge struct
type AddCharge struct {
	ChargeDate string     `json:"chargedate"`
	Contract   Contract   `json:"contract"`
	Object     Object     `json:"object"`
	ObjTypeId  int        `json:"objtypeid"`
	Pu         Pu         `json:"pu"`
	ChargeType ChargeType `json:"chargetype"`
	Qty        float32    `json:"qty"`
	TransLoss  float32    `json:"transloss"`
	Lineloss   float32    `json:"lineloss"`
	Startdate  string     `json:"startdate"`
	Enddate    string     `json:"enddate"`
}

//Charge_count  struct
type Charge_count struct {
	Values []Charge `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}
