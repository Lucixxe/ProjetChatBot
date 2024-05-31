package main

import (
	"database/sql"
	"log"
	"testing"
	"time"
)

func init_db() {
	var err error
	db, err = sql.Open("sqlite3", "./chatbot.db")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreation(t *testing.T) {
	init_db()
	a := get_agenda("bot")
	if a == nil {
		t.Fatalf("La création de l'agenda a échoué")
	}
	db.Close()
}

func TestNewEntry(t *testing.T) {
	init_db()
	a := get_agenda("bot")
	a.new_entry_str(-1, time.Now().Format(time.RFC3339), "ceci est un test", false)
	if len(a.entries) == 0 {
		t.Fatalf("création de l'entrée à échoué")
	}
	db.Close()
}

func TestFetchRemoveEntry(t *testing.T) {
	init_db()
	a := get_agenda("bot")
	a.new_entry_str(-1, time.Now().Format(time.RFC3339), "ceci est un test", false)
	get := a.time_check()
	if get == nil || len(a.removed) == 0 {
		t.Fatalf("suppression de l'entrée à échoué")
	}
	db.Close()
}
