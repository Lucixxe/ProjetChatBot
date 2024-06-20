package main

// page de connexion et d'inscription

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
	action := r.FormValue("action")

	// à faire, définir des limites sur les champs
	if len(username) == 0 || len(password) == 0 {
		log.Println("mdp ou pseudo trop court")
		// retourne erreur
		return false
	}

	// on enlève les espaces superflus
	password = strings.Trim(password, " ")
	username = strings.Trim(username, " ")

	if action == "register" {
		// register
		// l'utilisateur existe déjà ?
		ok, err := existe_utilisateur(username)
		if err != nil {
			log.Println(err)
			return false
		}
		if ok {
			// oui
			log.Println("l'utilisateur existe déjà")
			return false
		}

		// ajoute l'utilisateur
		user = ajout_utilisateur(username, password)
		// redirige sur le chat
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
	} else {
		// login
		// la combinaison pseudo, mot de passe est dans la base de données ?
		ok, err := verification_mdp(username, password)
		if err != nil {
			// non, mdp incorrect ou utilisateur inexistant
			log.Println(err)
			return false
		}
		if ok {
			// oui, redirige sur le chat
			user = &Utilisateur{ username }
			http.Redirect(w, r, "/chat", http.StatusSeeOther)
		} else {
			// impossible ! si faux verification_mdp renvoi toujours une erreur !
			return false
		}
	}
	return true
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


