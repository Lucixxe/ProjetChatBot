package main

import (
	"log"
	"strings"
	"encoding/json"
	"os"
	"io/ioutil"
)

type JSONMessage struct {
	Id  string `json:"id"`
	Role    string `json:"role"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

func appendMessageToJSONFile(message JSONMessage, filename string) error {
	var messages []JSONMessage
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, &messages)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		messages = []JSONMessage{}
	} else {
		return err
	}

	messages = append(messages, message)

	data, err = json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

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

	return messageDate
}
