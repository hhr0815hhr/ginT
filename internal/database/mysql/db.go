package mysql

import (
	"fmt"
	"time"

	"github.com/hhr0815hhr/gint/internal/config"
	"github.com/hhr0815hhr/gint/internal/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func init() {
	conf := config.Conf
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.Database.Prefix,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Logger.Fatal("Failed to connect to database,err: " + err.Error())
	}
	DB = db
	sqlDB, _ := db.DB()
	//设置连接池
	sqlDB.SetMaxOpenConns(conf.Database.MaxIdleConns)
	sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func ProvideDB() *gorm.DB {
	return DB
}
