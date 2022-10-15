package models

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
