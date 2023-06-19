package repositories

import (
	"mytask/models"

	"gorm.io/gorm"
)

type TransactionsRepo interface {
	FindTransaction() ([]models.Transaction, error)
	GetUserTransaction(ID int) ([]models.Transaction, error)
	GetPartnerTransaction(ID int) ([]models.Transaction, error)
	GetTransaction(ID int) (models.Transaction, error)
	CreateTransaction(Transaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(status string, orderId int) (models.Transaction, error)
	DeleteTransaction(Transaction models.Transaction, ID int) (models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repositories {
	return &repositories{db}
}

func (r *repositories) FindTransaction() ([]models.Transaction, error) {
	var Transaction []models.Transaction
	err := r.db.Preload("Cart.Product").Preload("Buyer").Preload("Seller").Find(&Transaction).Error

	return Transaction, err
}

func (r *repositories) GetUserTransaction(ID int) ([]models.Transaction, error) {
	var Transaction []models.Transaction
	err := r.db.Where("buyer_Id = ?", ID).Preload("Cart.Product").Preload("Buyer").Preload("Seller").Find(&Transaction).Error

	return Transaction, err
}
func (r *repositories) GetPartnerTransaction(ID int) ([]models.Transaction, error) {
	var Transaction []models.Transaction
	err := r.db.Where("seller_Id = ?", ID).Preload("Cart.Product").Preload("Buyer").Preload("Seller").Find(&Transaction).Error

	return Transaction, err
}

func (r *repositories) GetTransaction(ID int) (models.Transaction, error) {
	var Transaction models.Transaction
	err := r.db.Preload("Buyer").First(&Transaction, ID).Error

	return Transaction, err
}

func (r *repositories) CreateTransaction(Transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&Transaction).Error

	return Transaction, err
}

func (r *repositories) DeleteTransaction(Transaction models.Transaction, ID int) (models.Transaction, error) {
	err := r.db.Delete(&Transaction, ID).Scan(&Transaction).Error

	return Transaction, err
}

func (r *repositories) UpdateTransaction(status string, orderId int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Buyer").First(&transaction, orderId)

	transaction.Status = status
	err := r.db.Save(&transaction).Error
	return transaction, err
}
