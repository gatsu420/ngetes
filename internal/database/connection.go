package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func DBConn() (*bun.DB, error) {
	viper.SetConfigFile("./.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	dsn := "postgres://" + viper.GetString("db_user") + ":" + viper.GetString("db_password") + "@" + viper.GetString("db_addr") + "/" + viper.GetString("db_database") + "?sslmode=disable"
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
