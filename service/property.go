package service

import (
	"fmt"
	"time"

	"github.com/fox-one/fox-notifier/session"
	"github.com/fox-one/mixin-sdk/utils"
)

// Property property
type Property struct {
	Key       string `gorm:"type:varchar(128);PRIMARY_KEY"`
	Value     string `gorm:"type:varchar(256);"`
	UpdatedAt time.Time
}

// TableName table
func (Property) TableName() string {
	return "properties"
}

// ReadProperty read property
func ReadProperty(s *session.Session, key string) (string, error) {
	p := Property{Key: key}
	db := s.MysqlRead().Where(p).First(&p)
	if db.RecordNotFound() {
		return "", nil
	}

	return p.Value, db.Error
}

// WriteProperty write property
func WriteProperty(s *session.Session, key, value string) error {
	p := Property{Key: key}
	return s.MysqlWrite().Where(p).Assign(Property{Value: value}).FirstOrCreate(&p).Error
}

// ReadPropertyAsInt64 read property as int64
func ReadPropertyAsInt64(s *session.Session, key string) (int64, error) {
	value, err := ReadProperty(s, key)
	if err != nil {
		return 0, err
	}

	return utils.ParseInt64(value), nil
}

// WriteInt64Property write int64 property
func WriteInt64Property(s *session.Session, key string, value int64) error {
	return WriteProperty(s, key, fmt.Sprint(value))
}
