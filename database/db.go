package database

import (
	"database/sql"
	"fmt"
	"log"

	env "github.com/edwinbustillos/api-go-gin/utils"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	env.Load()
	host := env.GetEnv("HOST_DB", "")
	port := env.GetEnv("PORT_DB", "5432")
	user := env.GetEnv("USER_DB", "")
	password := env.GetEnv("PASSWORD_DB", "")
	dbname := env.GetEnv("DBNAME_DB", "")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")
	return db
}
