package main

import (
	"log"
	"strings"
)

type JSONMessage struct {
	Id  string `json:"id"`
	Destinataire    string `json:"destinataire"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

func loadMessageFromDB(userID string) ([]JSONMessage, error) {
	rows, err := db.Query("SELECT dest, date, contenu FROM messages WHERE id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []JSONMessage
	for rows.Next() {
		var msg JSONMessage
		if err := rows.Scan(&msg.Destinataire, &msg.Date, &msg.Content); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
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
