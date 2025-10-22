package model

import "gorm.io/gorm"

// AutoMigrate ensures the database schema matches the application models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
