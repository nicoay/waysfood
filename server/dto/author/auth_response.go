package authdto

type RegResponse struct {
	Email string `json:"email" `
	Token string `gorm:"type: varchar(255)" json:"token"`
}

type LoginResponse struct {
	Email string `gorm:"type: varchar(255)" json:"email"`
	Token string `gorm:"type: varchar(255)" json:"token"`
	Role  string `gorm:"type: varchar(255)" json:"role"`
	Image string `gorm:"type: varchar(255)" json:"image"`
	ID    int    `json:"id"`
}
