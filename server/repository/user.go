package repositories

import (
	"mytask/models"

	"gorm.io/gorm"
)

// sign a Contract
type UseRepository interface {
	FindUser() ([]models.User, error)
	FindPartner() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUser(ID int) (models.User, error)
	DeleteUser(user models.User, ID int) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
}

// func Connection

func RepositoryUser(db *gorm.DB) *repositories {
	return &repositories{db}
}

func (r *repositories) CreateUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

func (r *repositories) FindUser() ([]models.User, error) {
	var users []models.User
	var role = "As User"
	err := r.db.Where("role = ?", role).Find(&users).Error

	return users, err
}
func (r *repositories) FindPartner() ([]models.User, error) {
	var users []models.User
	var role = "As Partner"
	err := r.db.Where("role = ?", role).Find(&users).Error

	return users, err
}

func (r *repositories) GetUser(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

func (r *repositories) DeleteUser(user models.User, ID int) (models.User, error) {
	err := r.db.Delete(&user, ID).Scan(&user).Error

	return user, err
}

func (r *repositories) UpdateUser(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error

	return user, err
}
