package main

import (
	"net/http"
	"html/template"
	"log"
)

func empty (w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("pages/vide.tmpl"))
	err := tmp.ExecuteTemplate(w, "vide", struct {
		Name	string
	}{
		"Tout le monde :)",
	})
	if err != nil {
		log.Println(err)
	}
}

func main () {
	fs := http.FileServer(http.Dir("./public/"))

	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", empty)

	log.Println("Démarré sur le port 3333")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
