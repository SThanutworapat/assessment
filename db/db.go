package db

import (
	"database/sql"
	"log"

	"github.com/SThanutworapat/assessment/config"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	config := config.NewConfig()
	conn, err := sql.Open("postgres", config.Database_Url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := ` CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	); `
	_, err = conn.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
	return conn
}
