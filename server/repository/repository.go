package repositories

import "gorm.io/gorm"

// struc save connection n can use on global
type repositories struct {
	db *gorm.DB
}
