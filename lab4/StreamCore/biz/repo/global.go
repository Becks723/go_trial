package repo

import (
	"StreamCore/biz/repo/model"
	"StreamCore/config"
	"fmt"
	"log"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	// connect to mysql database
	c := config.Instance().MySql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.Username, c.Password, c.Host, c.Port, c.DBName)

	params := url.Values{}
	params.Set("charset", "utf8mb4")
	params.Set("parseTime", "true")

	dsn += fmt.Sprintf("?%s", params.Encode())
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error opening mysql: ", err.Error())
	}
	db = db_

	// auto migrate
	err = db_.AutoMigrate(&model.UserModel{}, &model.VideoModel{})
	if err != nil {
		log.Fatal("error migrating: ", err.Error())
	}
}
