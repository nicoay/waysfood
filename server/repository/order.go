package repositories

import (
	"mytask/models"

	"gorm.io/gorm"
)

type OrderRepo interface {
	CreateOrder(Order models.Order) (models.Order, error)
	DeleteOrder(Order models.Order, ID int) (models.Order, error)
	DeleteAllOrder(BID int) (models.Order, error)
	FindOrder() ([]models.Order, error)
	GetOrderByUserSeller(BID int, SID int) (models.Order, error)
	GetOrderByUserProduct(BID int, SID int) (models.Order, error)
	GetOrderbyUser(ID int) ([]models.Order, error)
	GetOrder(ID int) (models.Order, error)
	UpdateOrder(Order models.Order) (models.Order, error)
	GetByOrderProduct(BID int, SID int) (models.Order, error)
}

func RepositoryOrder(db *gorm.DB) *repositories {
	return &repositories{db}
}

func (r *repositories) FindOrder() ([]models.Order, error) {
	var Order []models.Order
	err := r.db.Preload("Buyer").Preload("Seller").Find(&Order).Error

	return Order, err
}
func (r *repositories) GetOrderbyUser(ID int) ([]models.Order, error) {
	var Order []models.Order
	err := r.db.Where("buyer_Id = ?", ID).Preload("Buyer").Preload("Product").Preload("Seller").Find(&Order).Error

	return Order, err
}

func (r *repositories) GetOrderByUserSeller(BID int, SID int) (models.Order, error) {
	var Order models.Order
	err := r.db.Where("buyer_Id = ? AND seller_Id != ?", BID, SID).First(&Order).Error

	return Order, err
}
func (r *repositories) GetOrderByUserProduct(BID int, IDP int) (models.Order, error) {
	var Order models.Order
	err := r.db.Where("buyer_Id = ? AND product_Id = ?", BID, IDP).First(&Order).Error

	return Order, err
}

func (r *repositories) GetOrder(ID int) (models.Order, error) {
	var Order models.Order
	err := r.db.First(&Order, ID).Error

	return Order, err
}

func (r *repositories) UpdateOrder(Order models.Order) (models.Order, error) {
	err := r.db.Save(&Order).Error
	return Order, err
}

func (r *repositories) GetByOrderProduct(BID int, IDP int) (models.Order, error) {
	var Order models.Order
	err := r.db.First(&Order, BID, IDP).Error

	return Order, err
}

func (r *repositories) CreateOrder(Order models.Order) (models.Order, error) {
	err := r.db.Preload("Buyer").Create(&Order).Error

	return Order, err
}

func (r *repositories) DeleteOrder(Order models.Order, ID int) (models.Order, error) {
	err := r.db.Delete(&Order, ID).Scan(&Order).Error

	return Order, err
}
func (r *repositories) DeleteAllOrder(ID int) (models.Order, error) {
	var Order models.Order
	err := r.db.Where("buyer_Id = ?", ID).Delete(&Order).Error

	return Order, err
}
