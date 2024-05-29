package main

import (
	"log"
	"strings"
)

func saveMessage(pseudo string, destinataire string, date string, contenu string) {
	insert, err := db.Prepare("insert into messages (id, dest, date, contenu) VALUES (?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = insert.Exec(pseudo, destinataire, date, contenu)
	if err != nil {
		log.Fatal(err)
	}

}

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
	log.Println("Content:", messageContent)

	return messageContent
}

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
	log.Println(messageDate)

	return messageDate
}
