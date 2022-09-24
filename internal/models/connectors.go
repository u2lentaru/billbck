package models

//Connector struct
type Connector struct {
	Id            int    `json:"id"`
	ConnectorName string `json:"connectorname"`
}

//AddConnector struct
type AddConnector struct {
	ConnectorName string `json:"connectorname"`
}

//Connector_count  struct
type Connector_count struct {
	Values []Connector `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
