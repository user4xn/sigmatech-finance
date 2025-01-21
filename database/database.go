package database

import (
	"clean-arch/pkg/config"
	"clean-arch/pkg/util"
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	dbConn  *gorm.DB
	rdb     *redis.Client
	once    sync.Once
	onceRdb sync.Once
)

func CreateConnection() {
	// Create database configuration information
	conf := dbConfig{
		User: config.MysqlUser(),
		Pass: config.MysqlPass(),
		Host: config.MysqlHost(),
		Port: config.MysqlPort(),
		Name: config.MysqlDBName(),
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

func BeginTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.WithContext(ctx).Begin()
}

func GetRedisClient() *redis.Client {

	onceRdb.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     util.GetEnv("REDIS_HOST", "localhost:6379"),
			Password: util.GetEnv("REDIS_PASSWORD", ""),
			DB:       0,
		})
	})

	return rdb
}
