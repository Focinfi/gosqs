package external

import (
	"github.com/jinzhu/gorm"
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var gormDB *gorm.DB

func init() {
	// if config.Config.Env.IsTest() {
	// 	return
	// }
	// dbCfg := config.Config.SQLDB
	// db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC", dbCfg.User, dbCfg.Password, dbCfg.Protocol, dbCfg.Host, dbCfg.Port, dbCfg.Name))
	// if err != nil {
	// 	panic(err)
	// }

	// gormDB = db
}

// MySQL wraps a mysql service
type MySQL struct {
	*gorm.DB
}

// NewMySQL create a new mysql client
func NewMySQL() *MySQL {
	return &MySQL{DB: gormDB}
}

// UserMySQL stores a user into MySQL
type UserMySQL struct {
	gorm.Model
	UniqueID string `json:"unique_id"`
}

// GetUserIDByUniqueID implements UserStore
func (db *MySQL) GetUserIDByUniqueID(uniqueID string) (int64, error) {
	user := UserMySQL{}
	if err := db.Where("unique_id = ?", uniqueID).First(&user).Error; err != nil {
		return -1, err
	}
	return int64(user.ID), nil
}

// CreateUserByUniqueID implements UserStore
func (db *MySQL) CreateUserByUniqueID(uniqueID string) (int64, error) {
	user := UserMySQL{UniqueID: uniqueID}
	if err := db.Create(user).Error; err != nil {
		return -1, err
	}
	return int64(user.ID), nil
}
