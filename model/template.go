package model

type Template struct {
	Subject string      `json:"subject"`
	Body    string      `json:"body"`
	Params  interface{} `json:"params"`
}
