package model

type MailMessege struct {
	Sender    string
	Recipient string
	Charset   string
	Subject   string
	Body      string
	Param     *UserParam
}

type UserParam struct {
	Name      string
	Age       int
	Residence string
	Gender    int
}
