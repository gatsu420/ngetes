package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	dbConfig struct {
		address  string
		database string
		user     string
		password string
	}

	redisConfig struct {
		address  string
		database int
		password string
	}
)

func init() {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error while trying to read config file: %v", err))
	}

	dbConfig.address = viper.GetString("DB_ADDR")
	dbConfig.database = viper.GetString("DB_DATABASE")
	dbConfig.user = viper.GetString("DB_USER")
	dbConfig.password = viper.GetString("DB_PASSWORD")

	redisConfig.address = viper.GetString("REDIS_ADDR")
	redisConfig.database = viper.GetInt("REDIS_DATABASE")
	redisConfig.password = viper.GetString("REDIS_PASSWORD")
}

func DBConn() (*bun.DB, error) {
	dsn := fmt.Sprintf(`postgres://%v:%v@%v/%v?sslmode=disable`,
		dbConfig.user,
		dbConfig.password,
		dbConfig.address,
		dbConfig.database,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook())

	if err := checkConn(db); err != nil {
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
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.address,
		Password: redisConfig.password,
		DB:       redisConfig.database,
	})

	if err := checkRedisConn(rdb); err != nil {
		return nil, err
	}

	return rdb, nil
}

func checkRedisConn(rdb *redis.Client) error {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}
