package main

import (
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
	http.HandleFunc("/", api.GetRoot)
	http.HandleFunc("/auth", api.AuthHandler(conf.JwtSalt))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Something went wrong: " + err.Error())
	}
}
