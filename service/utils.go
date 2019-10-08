package service

import (
	"github.com/fox-one/gin-contrib/session"
	"github.com/jinzhu/gorm"
)

func setupDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&Property{}).Error; err != nil {
		return err
	}
	return nil
}

func init() {
	session.RegisterSetdb(setupDB)
}
