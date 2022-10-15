package models

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
