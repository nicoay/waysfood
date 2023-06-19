package transactiondto

type CreateTransactionReq struct {
	Status     string `json:"status" form:"status" gorm:"type: varchar(255)"`
	Qty        int    `json:"qty" form:"qty" gorm:"type: int"`
	SellerID   int    `json:"sellerID"`
	OrderID    int    `json:"orderId" form:"orderId"`
	TotalPrice int64  `json:"totalPrice"`
}

type UpdateTransactionReq struct {
	Status  string `json:"status" form:"status" gorm:"type: varchar(255)"`
	Qty     int    `json:"qty" form:"qty" gorm:"type: int"`
	OrderID int    `json:"orderId" form:"orderId"`
}
