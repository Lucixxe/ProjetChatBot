package main

// fonction en rapport à l'utilisateur

import (
	"errors"
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Utilisateur struct {
	id string
}

func existe_utilisateur(id string) (bool, error) {
	stmt, err := db.Prepare("select count(*) from comptes where id = ?;")
	if err != nil {
		return true, err
	}

	var number int
	err = stmt.QueryRow(id).Scan(&number)

	if err == sql.ErrNoRows || number == 0 {
		return false, nil
	}

	return true, err
}

func ajout_utilisateur(id string, password string) *Utilisateur {
	insert, err := db.Prepare("insert into comptes (id, mdp, description) VALUES (?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = insert.Exec(id, password, "")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(id, "a été inscrit")
	return &Utilisateur{
		id,
	}
}

func verification_mdp(id string, password string) (bool, error) {
	stmt, err := db.Prepare("select id, mdp from comptes where id = ?;")
	if err != nil {
		return false, err
	}

	rows, err := stmt.Query(id)
	defer rows.Close()
	for rows.Next() {
		var user_id string
		var pass string
		rows.Scan(&user_id, &pass)
		if user_id == id {
			return password == pass, nil
		}
	}

	return false, errors.New("utilisateur inexistant")
}
