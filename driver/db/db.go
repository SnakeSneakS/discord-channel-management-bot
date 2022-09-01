package db

import (
	"fmt"
	"log"
	"time"

	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var dsn string

func Init(env *entity.Environment) {
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", env.Database.User, env.Database.Password, env.Database.Host, env.Database.Port, env.Database.Database)
	db, err := connect()
	if err != nil {
		log.Fatalln(err)
	}

	if err := db.AutoMigrate(&entity.DiscordChannel{}); err != nil {
		log.Fatalln(err)
	}
	if err := db.AutoMigrate(&entity.DiscordChannelUser{}); err != nil {
		log.Fatalln(err)
	}
	if err := db.AutoMigrate(&entity.DiscordChannelSetting{}); err != nil {
		log.Fatalln(err)
	}
}

// connect to db
// dsn example: "user:pass@tcp(db_host:db_port)/{database}?charset=utf8&parseTime=True&loc=Local"
func connect() (*gorm.DB, error) {
	count := 0
	log.Printf("trying to connect database. ")

	var err error
	for count < 60 {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
		count++
		log.Printf("%d\n", count)
	}
	log.Print("\n")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDB() (*gorm.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nulll")
	}
	return db, nil
}
