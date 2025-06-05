package model

import (
	"github.com/hhr0815hhr/gint/internal/database/mysql"
)

func init() {
	autoMigrate()
}

func autoMigrate() {
	mysql.DB.AutoMigrate()
}
