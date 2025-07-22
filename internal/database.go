package internal

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/raingrave/apirest/configs"
)

var DB *sqlx.DB

func ConnectDB() {
	var err error
	DB, err = sqlx.Connect("postgres", configs.EnvDatabaseURL())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connected")
}
