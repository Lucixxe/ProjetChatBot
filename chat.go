package main

// communication websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/ollama/ollama/api"
)

type Message struct {
	Contenu string `json:"contenu"`
	Kind    string `json:"type"`
}

func chat(w http.ResponseWriter, r *http.Request) {
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmp := template.Must(template.ParseFiles("pages/chat.tmpl"))
	err := tmp.ExecuteTemplate(w, "chat", struct {
		Pseudo string
	}{
		user.id,
	})
	if err != nil {
		log.Println(err)
	}
}

var upgrader = websocket.Upgrader{}

/*
Communication entre les WebSockets
*/
func ws_con(w http.ResponseWriter, r *http.Request) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade : ", err)
	}
	defer c.Close()

	messages := []api.Message{
		api.Message{
			Role:    "system",
			Content: "Tu es l'assistant d'une personne âgée, tu dois la motiver et la conseiller à faire des activités sociales, intellectuelles ou physiques. Les réponses doivent être concises. Tu es empathique, discret et calme, compréhensif, encourageant, informatif et fiable. Tu ne dois cependant pas être trop formel sans non plus être trop amical.",
		},
	}

	// Chargement des messages depuis la base de données
	messages_for_history, err := loadMessageFromDB(user.id)
	if err != nil {
		log.Println("Error loading messages:", err)
		return
	}

	// Envoi des messages
	initialMessages, err := json.Marshal(messages_for_history)
	if err != nil {
		log.Println("Error marshalling messages:", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, initialMessages)
	if err != nil {
		log.Println("Error sending initial messages:", err)
		return
	}

	pending_msg := ""
	c.SetReadLimit(5000)
	//c.SetReadDeadline(time.Now().Add(120 * time.Second))
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error : ", err)
			}

			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		log.Printf("rcv : %s\n", message)

		// save user message into DB
		userMessage := string(message)
		date := extract_date_from_message(userMessage)
		content := extract_content_from_message(userMessage)
		saveMessage(user.id, "assistant", date, content)

		messages = append(messages, api.Message{
			Role:    "user",
			Content: string(message),
		})

		req := &api.ChatRequest{
			Model:    "llama3",
			Messages: messages,
		}
		err = client.Chat(context.Background(), req, func(m api.ChatResponse) error {
			json_data, err := json.Marshal(&Message{
				m.Message.Content,
				"message",
			})
			if err != nil {
				log.Fatal(err)
			}
			err = c.WriteMessage(mt, json_data)

			if err != nil {
				log.Fatal(err)
			}
			pending_msg += m.Message.Content

			if m.Done {
				current_time := time.Now()
				formatted_date := current_time.Format("02/01/2006 15:04")
				saveMessage(user.id, "user", formatted_date, pending_msg)

				messages = append(messages, api.Message{
					Role:    "assistant",
					Content: pending_msg,
				})
				pending_msg = ""
				json_data, err = json.Marshal(&Message{
					"#fin#",
					"message",
				})
				if err != nil {
					log.Fatal(err)
				}
				err = c.WriteMessage(mt, json_data)
				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
