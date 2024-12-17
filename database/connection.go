package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gatsu420/ngetes/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
)

var (
	dbConfig struct {
		host     string
		port     int
		database string
		user     string
		password string
	}

	cacheConfig struct {
		host     string
		port     int
		database int
		password string
	}
)

func init() {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Logger.Error("failed to read config file", zap.Error(err))
	}
	viper.AutomaticEnv()

	dbConfig.host = viper.GetString("POSTGRES_HOST")
	dbConfig.port = viper.GetInt("POSTGRES_PORT")
	dbConfig.database = viper.GetString("POSTGRES_DB")
	dbConfig.user = viper.GetString("POSTGRES_USER")
	dbConfig.password = viper.GetString("POSTGRES_PASSWORD")

	cacheConfig.host = viper.GetString("REDIS_HOST")
	cacheConfig.port = viper.GetInt("REDIS_PORT")
	cacheConfig.database = viper.GetInt("REDIS_DB")
	cacheConfig.password = viper.GetString("REDIS_PASSWORD")
}

func DBConn() (*bun.DB, error) {
	dsn := fmt.Sprintf(`postgres://%v:%v@%v:%v/%v?sslmode=disable`,
		dbConfig.user,
		dbConfig.password,
		dbConfig.host,
		dbConfig.port,
		dbConfig.database,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook())

	if err := checkConn(db); err != nil {
		logger.Logger.Error("failed to check DB connection", zap.Error(err))
		return nil, err
	}

	return db, nil
}

func checkConn(db *bun.DB) error {
	n := 0

	return db.NewSelect().ColumnExpr("1").Scan(context.Background(), &n)
}

func RedisConn() (*redis.Client, error) {
	viper.SetConfigFile("./.env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Error("failed to read config file", zap.Error(err))
	}

	cache := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", cacheConfig.host, cacheConfig.port),
		Password: cacheConfig.password,
		DB:       cacheConfig.database,
	})

	if err := checkRedisConn(cache); err != nil {
		logger.Logger.Error("failed to check redis connection", zap.Error(err))
		return nil, err
	}

	return cache, nil
}

func checkRedisConn(rdb *redis.Client) error {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}
