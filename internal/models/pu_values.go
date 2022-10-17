package models

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
