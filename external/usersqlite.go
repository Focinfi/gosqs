package external

import (
	"github.com/Focinfi/gosqs/config"
	"github.com/jinzhu/gorm"
	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	if config.Config.Env.IsTest() {
		return
	}
	db, err := gorm.Open("sqlite3", "gosqs.db")
	if err != nil {
		panic(err)
	}

	gormDB = db
}

// SQLite wraps a mysql service
type SQLite struct {
	*gorm.DB
}

// NewSQLite create a new mysql client
func NewSQLite() *MySQL {
	return &MySQL{DB: gormDB}
}

// UserSQLite stores a user into MySQL
type UserSQLite struct {
	gorm.Model
	UniqueID string `json:"unique_id"`
}

// GetUserIDByUniqueID implements UserStore
func (db *SQLite) GetUserIDByUniqueID(uniqueID string) (int64, error) {
	user := UserMySQL{}
	if err := db.Where("unique_id = ?", uniqueID).First(&user).Error; err != nil {
		return -1, err
	}
	return int64(user.ID), nil
}

// CreateUserByUniqueID implements UserStore
func (db *SQLite) CreateUserByUniqueID(uniqueID string) (int64, error) {
	user := UserMySQL{UniqueID: uniqueID}
	if err := db.Create(user).Error; err != nil {
		return -1, err
	}
	return int64(user.ID), nil
}
