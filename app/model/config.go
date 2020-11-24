package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	MysqlDsn = `user:YbcuGH4Ks@z6$@tcp(3.15.137.94)/bank?parseTime=True&charset=utf8`
	//MysqlDsn = `user:YbcuGH4Ks@z6$@tcp(127.0.0.1)/cocreate?parseTime=True&charset=utf8` //localhost
	//MysqlDsn = `root:YbcuGH4Ks@z6$@tcp(localhost:3306)/bank?parseTime=True&charset=utf8`
)

func init() {
	var err error
	for {
		DB, err = gorm.Open(mysql.Open(MysqlDsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
		if err == nil {
			break
		}
	}
	DB.AutoMigrate(new(User), new(Kategori), new(Detail_category), new(Admin), new(Posting), new(Comment), new(UserTemporary), new(UserReject), new(GrupProject), new(Project), new(Task), new(CommentTask), new(Expertise), new(Exprience), new(Portofolio))
}
