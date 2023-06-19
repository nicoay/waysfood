package repositories

import (
	"mytask/models"

	"gorm.io/gorm"
)

type CartRepo interface {
	CreateCarts(Carts models.Carts) (models.Carts, error)
	DeleteCarts(Carts models.Carts, ID int) (models.Carts, error)
	FindCarts() ([]models.Carts, error)
	GetCarts(ID int) (models.Carts, error)
}

func RepositoryCart(db *gorm.DB) *repositories {
	return &repositories{db}
}

func (r *repositories) FindCarts() ([]models.Carts, error) {
	var carts []models.Carts
	err := r.db.Find(&carts).Error

	return carts, err
}

func (r *repositories) GetCarts(ID int) (models.Carts, error) {
	var carts models.Carts
	err := r.db.First(&carts, ID).Error

	return carts, err
}

func (r *repositories) CreateCarts(Carts models.Carts) (models.Carts, error) {
	err := r.db.Create(&Carts).Error

	return Carts, err
}

func (r *repositories) DeleteCarts(Carts models.Carts, ID int) (models.Carts, error) {
	err := r.db.Delete(&Carts, ID).Scan(&Carts).Error

	return Carts, err
}
