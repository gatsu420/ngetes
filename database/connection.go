package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gatsu420/ngetes/config"
	"github.com/gatsu420/ngetes/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
)

func DBConn(config *config.Config) (*bun.DB, error) {
	dsn := fmt.Sprintf(`postgres://%v:%v@%v:%v/%v?sslmode=disable`,
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresDB,
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

func RedisConn(config *config.Config) (*redis.Client, error) {
	viper.SetConfigFile("./.env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Error("failed to read config file", zap.Error(err))
	}

	cache := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       config.RedisDB,
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
