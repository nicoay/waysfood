package repositories

import (
	"mytask/models"

	"gorm.io/gorm"
)

type ProductRepo interface {
	FindProduct() ([]models.Product, error)
	FindProductPartner(ID int) ([]models.Product, error)
	GetProduct(ID int) (models.Product, error)
	CreateProduct(Product models.Product) (models.Product, error)
	UpdateProduct(Product models.Product) (models.Product, error)
	DeleteProduct(Product models.Product, ID int) (models.Product, error)
}

func RepositoryProduct(db *gorm.DB) *repositories {
	return &repositories{db}
}

func (r *repositories) FindProduct() ([]models.Product, error) {
	var Products []models.Product
	err := r.db.Preload("User").Find(&Products).Error

	return Products, err
}

func (r *repositories) FindProductPartner(ID int) ([]models.Product, error) {
	var Products []models.Product
	err := r.db.Preload("User").Find(&Products, "user_id = ?", ID).Error

	return Products, err
}

func (r *repositories) GetProduct(ID int) (models.Product, error) {
	var Products models.Product
	err := r.db.First(&Products, ID).Error

	return Products, err
}

func (r *repositories) CreateProduct(Product models.Product) (models.Product, error) {
	err := r.db.Create(&Product).Error

	return Product, err
}

func (r *repositories) UpdateProduct(Product models.Product) (models.Product, error) {
	err := r.db.Save(&Product).Error
	return Product, err
}

func (r *repositories) DeleteProduct(Product models.Product, ID int) (models.Product, error) {
	err := r.db.Delete(&Product, ID).Scan(&Product).Error

	return Product, err
}
