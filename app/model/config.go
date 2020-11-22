package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	MysqlDsn = `user:@tcp(3.15.137.94)/bank?parseTime=True&charset=utf8`
	//MysqlDsn = `root:@tcp(localhost:3306)/bank?parseTime=True&charset=utf8`
)

func init() {
	var err error
	for {
		DB, err = gorm.Open(mysql.Open(MysqlDsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
		if err == nil {
			break
		}
	}
	DB.AutoMigrate(new(User), new(Kategori), new(Detail_category), new(Admin), new(Posting), new(Comment), new(UserTemporary), new(UserReject))
}
