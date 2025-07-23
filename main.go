package main

import (
	"database/sql"
	"go-bank/api"
	db "go-bank/db/sqlc"
	"go-bank/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("Failed to connect to db")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server")
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server")
	}
}
