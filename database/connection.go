package database

import (
	"fmt"
	"log"
	"onlibrary/config"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)



var onceDb sync.Once

var instance *gorm.DB

func GetInstance() *gorm.DB {
	onceDb.Do(func() {
		databaseConfig := config.DatabaseNew().(*config.DatabaseConfig)
		dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",databaseConfig.Mysql.DbUsername, databaseConfig.Mysql.DbPassword,databaseConfig.Mysql.DbPort,databaseConfig.Mysql.DbDatabase)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		// db, err := gorm.Open("mysql", fmt.Sprintf("host:%s port:%s user=%s dbname=%s password=%s",
		// 		databaseConfig.Mysql.DbHost,databaseConfig.Mysql.DbPort,databaseConfig.Mysql.DbUsername, databaseConfig.Mysql.DbUsername))
		if err != nil {
			log.Fatalf("Could not connect database : %v", err)
		}

		instance = db
	})
	return instance
}