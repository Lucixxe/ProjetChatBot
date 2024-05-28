package main

import (
	"net/http"
)

func deconnexion (w http.ResponseWriter, r *http.Request) {
	// déconnecte l'utilisateur
	user = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
