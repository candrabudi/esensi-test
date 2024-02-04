package database

import (
	"esensi-test/pkg/util"
	"sync"

	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"
)

var (
	dbConn *gorm.DB
	once   sync.Once
	rdb    *redis.Client
)

func CreateConnection() {
	// Create database configuration information
	conf := dbConfig{
		User: util.GetEnv("MYSQL_USER", "root"),
		Pass: util.GetEnv("MYSQL_PASS", ""),
		Host: util.GetEnv("MYSQL_HOST", "localhost"),
		Port: util.GetEnv("MYSQL_PORT", "3306"),
		Name: util.GetEnv("MYSQL_DB_NAME", "user_db"),
	}

	mysql := mysqlConfig{dbConfig: conf}
	// Create only one mysql Connection, not the same as mysql TCP connection
	once.Do(func() {
		mysql.Connect()
	})
}

func GetConnection() *gorm.DB {
	// Check db connection, if exist return the memory address of the db connection
	if dbConn == nil {
		CreateConnection()
	}
	return dbConn
}

func GetRedisClient() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     util.GetEnv("REDIS_HOST", "localhost:6379"),
		Password: util.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})
	return rdb
}
