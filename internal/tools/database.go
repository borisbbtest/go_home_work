package tools

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func PingDataBase(psqlInfo string) (bool, error) {

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return false, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}
