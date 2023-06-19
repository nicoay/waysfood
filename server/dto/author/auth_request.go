package authdto

type RegisterRequest struct {
	FullName string `json:"fullname" gorm:"type: varchar(255)" `
	Email    string `json:"email" gorm:"type: varchar(255)" `
	Password string `json:"password" gorm:"type: varchar(255)" `
	Phone    string `json:"phone" gorm:"type: varchar(255)" `
	Address  string `json:"address" gorm:"type: varchar(255)" `
	Role     string `json:"role" grom:"type: varchar(255)"`
	Gender   string `json:"gender" grom:"type: varchar(255)"`
}

type LoginRequest struct {
	Email    string `json:"email" gorm:"type: varchar(255)" validate:"required"`
	Password string `json:"password" gorm:"type: varchar(255)" validate:"required"`
}
