package main

import (
	"net/http"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)	

var db *sql.DB

func main () {
	var err error

	db, err = sql.Open("sqlite3", "./chatbot.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fs := http.FileServer(http.Dir("./public/"))

	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", connexion)
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/ws", ws_con)

	log.Println("Démarré sur le port 3333")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
