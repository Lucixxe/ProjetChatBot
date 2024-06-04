package main

import (
	"log"
	"net/http"
)

func deconnexion(w http.ResponseWriter, r *http.Request) {
	// déconnecte l'utilisateur
	log.Println("déconnexion de ", user.id)
	delete := r.URL.Query().Get("delete")
	log.Println("delete = ", delete)
	if delete != "" && delete == "true" {
		log.Println("suppression des données de ", user.id)
		supprime_compte(user)
	}
	user = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
