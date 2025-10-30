package repository

import (
	"fmt"
	"log"
	"memo/config"
	"memo/repository/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

/* 初始化数据库 */
func Load() {
	// 连接数据库
	c := config.Instance().MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		c.Username, c.Password, c.Host, c.Port, c.DBName)
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db = db_

	// 自动建表
	err = db.AutoMigrate(&model.UserModel{}, &model.MemoModel{})
	if err != nil {
		log.Fatal(err)
	}
}
