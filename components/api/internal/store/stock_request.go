package store

type StockRequest struct {
	Id       string `dynamodbav:"Id"`
	ISBN     string `dynamodbav:"ISBN"`
	Quantity int    `dynamodbav:"Quantity"`
}
