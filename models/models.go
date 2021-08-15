package models

type StkPush struct {
	BusinessShortCode string
	Password          string
	Timestamp         string
	TransactionType   string
	Amount            string
	PartyA            string
	PartyB            string
	PhoneNumber       string
	CallBackURL       string
	AccountReference  string
	TransactionDesc   string
}

type StkPushQuery struct {
	BusinessShortCode string 
	Password string 
	Timestamp string 
	CheckoutRequestID string
}

type Wallet struct {
	ID           int // phoneNumber
	Balance      int64
	Transactions []Transaction `gorm:"foreignKey:PhoneNumber;references:ID"`
}
type Transaction struct {
	ID          int `gorm:"primaryKey;autoIncrement:true"`
	PhoneNumber int64
	Amount      int64
	CheckoutID  string
	Status      string
}
