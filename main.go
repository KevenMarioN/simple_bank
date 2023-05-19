package main

import (
	"database/sql"
	"log"

	"github.com/KevenMarioN/simple_bank/api"
	db "github.com/KevenMarioN/simple_bank/db/sqlc"
	"github.com/KevenMarioN/simple_bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
