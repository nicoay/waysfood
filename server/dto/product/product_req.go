package productdto

type CreateProductReq struct {
	Title  string `json:"title" form:"title" gorm:"type: varchar(255)"`
	Price  int64  `json:"price" form:"price" gorm:"type: int"`
	Image  string `json:"image" form:"image" gorm:"type: varchar(255)"`
	Qty    int    `json:"qty" form:"qty" gorm:"type: int"`
	UserId int    `json:"user_id" form:"user_id" gorm:"type: int"`
}

type UpdateProductRequest struct {
	Title string `json:"title" form:"title" gorm:"type: varchar(255)"`
	Price int64  `json:"price" form:"price" gorm:"type: int"`
	Image string `json:"image" form:"image" gorm:"type: varchar(255)"`
	Qty   int    `json:"qty" form:"qty" gorm:"type: int"`
}
