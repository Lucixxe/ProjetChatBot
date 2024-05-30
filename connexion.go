package main

// page de connexion et d'inscription

import (
	"net/http"
	"html/template"
	"log"
	"strings"
)

const shift = 3	// decalage de 3 vers la droite

func caesar_cipher(password string) string {
	var encrypted_password strings.Builder

	for _, char := range password {
        if char >= 'a' && char <= 'z' {
            encrypted_password.WriteRune('a' + (char-'a'+rune(shift))%26)
        } else if char >= 'A' && char <= 'Z' {
            encrypted_password.WriteRune('A' + (char-'A'+rune(shift))%26)
        } else if char >= '0' && char <= '9' {
            encrypted_password.WriteRune('0' + (char-'0'+rune(shift))%10)
        } else {
            encrypted_password.WriteRune(char)
        }
    }

    return encrypted_password.String()
}

func decrypt_caesar_cipher(encrypted_password string) string {
	var decrypted_password strings.Builder

    for _, char := range encrypted_password {
        if char >= 'a' && char <= 'z' {
            decrypted_password.WriteRune('a' + (char-'a'-rune(shift)+26)%26)
        } else if char >= 'A' && char <= 'Z' {
            decrypted_password.WriteRune('A' + (char-'A'-rune(shift)+26)%26)
        } else if char >= '0' && char <= '9' {
            decrypted_password.WriteRune('0' + (char-'0'-rune(shift)+26)%10)
        }else {
            decrypted_password.WriteRune(char)
        }
    }

    return decrypted_password.String()
}


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
		user = ajout_utilisateur(username, caesar_cipher(password))
		// redirige sur le chat
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
	} else {
		// login
		// la combinaison pseudo, mot de passe est dans la base de données ?
		ok, err := verification_mdp(username, caesar_cipher(password))
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
			// impossible ! si faux verification_mdp renvoi toujours une erreure !
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


