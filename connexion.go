package main

import (
	"net/http"
	"html/template"
	"log"
	"strings"
)

func connexion_query (w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	username := r.FormValue("pseudo")
	password := r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
		log.Println("mdp ou pseudo trop court")
		// retourne erreur
		return false
	}

	password = strings.Trim(password, " ")
	username = strings.Trim(username, " ")

	ok, err := existe_utilisateur(username)
	if err != nil {
		log.Println(err)
		return false // message d'erreur
	}
	if ok {
		ok, err = verification_mdp(username, password)
		if err != nil {
			log.Println(err)
			return false // message d'erreur
		}

		if ok {
			// mot de passe correct redirection
			log.Println("connexion de ", username)
			http.Redirect(w, r, "/chat",  http.StatusSeeOther)
			return true
		} else {
			log.Println("mot de passe incorrect")
			return false // message d'erreur mot de passe incorrect
		}
	} else {
		// valeur de retour ignoré pour l'instant
		ajout_utilisateur(username, password)
		http.Redirect(w, r, "/chat",  http.StatusSeeOther) // redirige vers le chat
		return true
	}
	return false
}

func connexion (w http.ResponseWriter, r *http.Request) {
	no_error := true
	if r.Method == http.MethodPost {
		no_error = connexion_query(w, r)
	}

	tmp := template.Must(template.ParseFiles("pages/connexion.tmpl"))
	err := tmp.ExecuteTemplate(w, "connexion", struct{
		Error	bool
	}{
		!no_error,
	})
	if err != nil {
		log.Println(err)
	}
}


