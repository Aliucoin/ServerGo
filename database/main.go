package database

import (
	"database/sql"
	"fmt"
	"log"

	"runtime"
	"server-go/common"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func InitDB() {
	config := common.Config
	DB = bun.NewDB(sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(config.DB.IP),
		pgdriver.WithUser(config.DB.User),
		pgdriver.WithPassword(config.DB.Password),
		pgdriver.WithDatabase(config.DB.Name),
		pgdriver.WithTLSConfig(nil),
	)), pgdialect.New())

	if config.Debug {
		DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	fmt.Println("Max open conns:", maxOpenConns)
	DB.SetMaxOpenConns(maxOpenConns)
	DB.SetMaxIdleConns(maxOpenConns)

	// create database structure if doesn't exist
	if err := CreateSchemas(); err != nil {
		log.Println("Failed to create schema")
		log.Panic(err)
	}
}
