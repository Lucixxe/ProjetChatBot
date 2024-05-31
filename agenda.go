package main

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type Agenda struct {
	user_id string
	entries []*AgendaEntry
	removed []*AgendaEntry
}

func (a *Agenda) send_agenda() string {
	json_msg := `{ "type": "agenda", "events": [`
	for _, ae := range a.entries {
		json_msg += fmt.Sprintf(`{ "content": "%s", "date": "%s" },`, ae.reminder, ae.date.Format("02/01/2006 15:04"))
	}
	json_msg += `]}`
	return json_msg
}

func get_agenda(id string) *Agenda {
	a := &Agenda{
		id,
		[]*AgendaEntry{},
		[]*AgendaEntry{},
	}
	a.load()
	return a
}

func (a *Agenda) time_check() *AgendaEntry {
	now := time.Now()
	if len(a.entries) == 0 {
		return nil
	}
	if a.entries[0].date.Before(now) {
		save := a.entries[0]
		a.removed = append(a.removed, save)
		a.entries = a.entries[1:]
		return save
	}
	return nil
}

func (a *Agenda) load() {
	stmt, err := db.Prepare("select numero, date, contenu from agenda where id=?;")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query(a.user_id)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var numero int
		var date string
		var contenu string
		rows.Scan(&numero, &date, &contenu)
		a.new_entry_str(numero, date, contenu, true)
	}
	// sort entries
	a.sort_by_date()
	rows.Close()
}

func (a *Agenda) remove() {
	// remove entries in agenda.removed
	if len(a.removed) == 0 {
		return
	}
	query := "delete from agenda where "
	parameters := []any{}
	for _, ae := range a.removed {
		if len(parameters) > 0 {
			query += " OR "
		}
		parameters = append(parameters, ae.id)
		query += "(numero = ?)"
	}
	query += ";"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(parameters...)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Agenda) save() {
	// add entries not in database
	query := "insert into agenda (id, date, contenu) values "
	parameters := []any{}
	for _, ae := range a.entries {
		if !ae.in_db {
			if len(parameters) > 0 {
				query += ", "
			}
			parameters = append(parameters, a.user_id, ae.date.Format(time.RFC3339), ae.reminder)
			query += "(?, ?, ?)"
		}
	}
	if len(parameters) == 0 {
		return
	}
	query += ";"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(parameters...)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Agenda) new_entry_str(id int, date string, content string, in_db bool) {
	format_date, err := time.Parse(time.RFC3339, date)
	if err != nil {
		log.Fatal(err)
	}
	a.entries = append(a.entries, &AgendaEntry{
		id,
		format_date,
		content,
		in_db,
	})
	// sort if not in db
	if !in_db {
		a.sort_by_date()
	}
}

func (a *Agenda) sort_by_date() {
	sort.Slice(a.entries, func(i, j int) bool {
		return a.entries[i].date.Before(a.entries[j].date)
	})
}

type AgendaEntry struct {
	id       int
	date     time.Time
	reminder string
	in_db    bool
}
