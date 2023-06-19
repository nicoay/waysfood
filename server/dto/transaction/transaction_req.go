package transactiondto

type TransactionDeleteRes struct {
	ID int `json:"id" gorm:"primary_key:auto_increment"`
}
