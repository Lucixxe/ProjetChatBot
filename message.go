package main

import (
	"log"
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"database/sql"
)

type JSONMessage struct {
	Id  string `json:"id"`
	Destinataire    string `json:"destinataire"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

func exportMessagesToJSON(db *sql.DB, filename string) error {
	rows, err := db.Query("SELECT id, dest, date, contenu FROM messages")
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var messages []JSONMessage

	for rows.Next() {
		var message JSONMessage
		err := rows.Scan(&message.Id, &message.Destinataire, &message.Date, &message.Content)
		if err != nil {
			log.Println(err)
			continue
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %v", err)
	}

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshalling failed: %v", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("writing to file failed: %v", err)
	}

	fmt.Printf("Exported messages to %s successfully.\n", filename)
	return nil
}


func saveMessage(pseudo string, destinataire string, date string, contenu string) {
	contenu = strings.ReplaceAll(contenu, "\n", " \n ")

	insert, err := db.Prepare("insert into messages (id, dest, date, contenu) VALUES (?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = insert.Exec(pseudo, destinataire, date, contenu)
	if err != nil {
		log.Fatal(err)
	}

}

/*
extract_content_from_message
	Paramètres: message, une chaîne de caractères de type {"message":"Bonjour","date":"18/06/2024 09:19"}
	Valeur de retour: Une chaîne de caractères contenant le contenu du message
*/
func extract_content_from_message(message string) string {
	messagePrefix := `"message":"`
	startIndex := strings.Index(message, messagePrefix)

	if startIndex == -1 {
		log.Println("Message non trouvé")
	}

	startIndex += len(messagePrefix)
	endIndex := strings.Index(message[startIndex:], `"`)

	if endIndex == -1 {
		log.Println("Message non trouvé")
	}

	messageContent := message[startIndex : startIndex+endIndex]

	return messageContent
}

/*
extract_date_from_message
	Paramètres: message, une chaîne de caractères de type {"message":"Bonjour\n","date":"18/06/2024 09:19"}
	Valeur de retour: Une chaîne de caractères contenant la date du message
*/
func extract_date_from_message(message string) string {
	datePrefix := `"date":"`
	startIndex := strings.Index(message, datePrefix)

	if startIndex == -1 {
		log.Println("Date non trouvée")
		return ""
	}

	startIndex += len(datePrefix)
	endIndex := strings.Index(message[startIndex:], `"`)

	if endIndex == -1 {
		log.Println("Date non trouvée")
		return ""
	}

	messageDate := message[startIndex : startIndex+endIndex]

	return messageDate
}
