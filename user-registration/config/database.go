package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
	Ctx   = context.Background()
)

func InitMySQL() {
	var err error
	DB, err = gorm.Open(mysql.Open(AppConfig.MySQL.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("MySQL connected successfully")
}

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr(),
		Password: AppConfig.Redis.Password,
		DB:       AppConfig.Redis.Database,
	})

	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
}
