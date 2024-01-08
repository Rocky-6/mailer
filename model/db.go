package model

type User struct {
	Email     string `dynamodbav:"mail_address"`
	Name      string `dynamodbav:"name"`
	Age       int    `dynamodbav:"age"`
	Residence string `dynamodbav:"residence"`
	Gender    int    `dynamodbav:"gender"`
}

type Filter struct {
	MinAge    int
	MaxAge    int
	Residence string
	Gender    int
}
