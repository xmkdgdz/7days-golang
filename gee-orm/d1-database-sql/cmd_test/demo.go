package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // 导入时会注册 sqlite3 的驱动
)

func main() {
	db, _ := sql.Open("sqlite3", "gee.db")
	defer func() { _ = db.Close() }()
	_, _ = db.Exec("DROP TABLE IF EXISTS user")
	_, _ = db.Exec("CREATE TABLE User(Name text);")
	result, err := db.Exec("INSERT INTO User(`Name`) VALUES(?), (?)", "Tom", "Jack")
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}
	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err != nil {
		log.Println(err)
	} else {
		log.Println(name)
	}
}
