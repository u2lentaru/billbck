package models

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
