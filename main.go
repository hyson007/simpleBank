package main

import (
	"database/sql"
	"log"

	"github.com/hyson007/simpleBank/api"
	db "github.com/hyson007/simpleBank/db/sqlc"
	"github.com/hyson007/simpleBank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v\n", server)
	if err := server.Start(config.ServerAddress); err != nil {
		panic(err)
	}

}
