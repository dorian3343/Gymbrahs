package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"main/api"
	"main/config"
	"net/http"
)

func main() {
	conf, err := config.ConfFromFile("config.json")
	if err != nil {
		log.Println("Something went wrong: " + err.Error())
		return
	}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.Println("Web server turned on. ")

	db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// This needs to be removed in final version
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT,
				password TEXT,
				email TEXT
			);
		`)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", api.GetRoot)
	http.HandleFunc("/auth", api.AuthHandler(conf.JwtSalt, db))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Something went wrong: " + err.Error())
	}
}
