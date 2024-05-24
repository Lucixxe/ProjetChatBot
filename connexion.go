package main

import (
	"net/http"
	"html/template"
	"log"
	"strings"
)

func connexion_query (w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("pseudo")
	password := r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
		log.Println("mdp ou pseudo trop court")
		// retourne erreur
		return
	}

	password = strings.Trim(password, " ")
	username = strings.Trim(username, " ")

	ok, err := existe_utilisateur(username)
	if err != nil {
		log.Println(err)
		return	// message d'erreur
	}
	if ok {
		ok, err = verification_mdp(username, password)
		if err != nil {
			log.Println(err)
			return	// message d'erreur
		}

		if ok {
			// mot de passe correct redirection
			log.Println("connexion de ", username)
			http.Redirect(w, r, "/chat",  http.StatusSeeOther)
		} else {
			return // message d'erreur mot de passe incorrect
		}
	} else {
		// valeur de retour ignoré pour l'instant
		ajout_utilisateur(username, password)
		http.Redirect(w, r, "/chat",  http.StatusSeeOther) // redirige vers le chat
	}

}

func connexion (w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		connexion_query(w, r)
	}

	tmp := template.Must(template.ParseFiles("pages/connexion.tmpl"))
	err := tmp.ExecuteTemplate(w, "connexion", nil)
	if err != nil {
		log.Println(err)
	}
}


