package dbops

import "database/sql"

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dbConn, err = sql.Open("mysql", "LemonLoser:1413139430@tcp(localhost:3306)/video_server?charset=utf-8")
	if err != nil {
		panic(err.Error())
	}
}
