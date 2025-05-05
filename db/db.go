package db

import (
	"database/sql"
	"fmt"
	"log"
)

func LocalDbConnect() (*sql.DB, error) {
	log.Println("LocalDbConnect(+)")
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "root", "localhost", "3306", "dev")
	lDb, lErr := sql.Open("mysql", connString)
	if lErr != nil {
		log.Println("Open connection failed:", lErr.Error())
		return nil, lErr
	}

	if err := lDb.Ping(); err != nil {
		log.Println("Ping failed:", err.Error())
		return nil, err
	}

	log.Println("LocalDbConnect(-)")
	return lDb, nil
}
