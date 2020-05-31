package mysql_transaction

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func InitDatabase() {

	DbParam1 := DbParam{}

	DbParam1.User = "root"
	DbParam1.Password = "88888888"
	DbParam1.Host = "127.0.0.1"
	//DbParam1.Port = "4000"
	DbParam1.Port = "3306"
	DbParam1.Schema = "test"

	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DbParam1.User, DbParam1.Password, DbParam1.Host, DbParam1.Port, DbParam1.Schema,
	)

	var err error
	DB, err = gorm.Open("mysql", path)
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	//log.Println(DB)
	//DB.LogMode(true)
	DB.LogMode(false)
	DB.SingularTable(true)

	//var l L
	//DB.SetLogger(l)

}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LstdFlags)
	InitDatabase()

	ResetTable()
	//DB.LogMode(false)

}
