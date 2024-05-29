package main

import (
	"net/http"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)	

var db *sql.DB
var user *Utilisateur

func main () {
	var err error

	user = nil
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
	http.HandleFunc("/disconnect", deconnexion)

	log.Println("Démarré sur le port 3333")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
