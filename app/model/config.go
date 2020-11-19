package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	MysqlDsn = `user:@tcp(3.15.137.94)/bank?parseTime=True&charset=utf8`
)

func init() {
	var err error

	/*db_u := "user"           //os.Getenv("DB_U")
	db_p := ""               //os.Getenv("DB_P")
	db_host := "3.15.137.94" //os.Getenv("DB_HOST")
	db_name := "bank"        //os.Getenv("DB_NAME")*/

	// DB, err = gorm.Open(mysql.Open(fmt.Sprintf("root:root@tcp(172.18.0.10:3306)/digitalent_bank")), &gorm.Config{})
	for {
		//DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s", db_u, db_p, db_host, db_name)), &gorm.Config{})
		DB, err = gorm.Open(mysql.Open(MysqlDsn), &gorm.Config{})
		if err == nil {
			break
		}
	}
	DB.AutoMigrate(new(User), new(Kategori), new(Detail_category), new(Admin), new(Posting), new(Comment), new(UserTemporary), new(UserReject))

}
