package users

import "gorm.io/gorm"

type User struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string
	Email string
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
