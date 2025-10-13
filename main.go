package main

import (
	"context"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = ":8080"
)

func main() {
	var err error
	var conn *pgxpool.Pool
	conn, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
