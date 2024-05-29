package main

// communication websocket

import (
	"net/http"
	"log"
	"time"
	"bytes"
	"html/template"
	"context"
	"os"

	"github.com/ollama/ollama/api"
	"github.com/gorilla/websocket"
)

func chat (w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("pages/chat.tmpl"))
	err := tmp.ExecuteTemplate(w, "chat", struct{
		Pseudo	string
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
func ws_con (w http.ResponseWriter, r *http.Request) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade : ", err)
	}
	defer c.Close()

	namefile := user.id
	file, err := os.OpenFile(namefile + ".txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// load messages from DB
	messages := []api.Message{
		api.Message{
			Role: "system",
			Content: "Tu es l'assistant d'une personne âgée, tu dois la motiver et la conseiller à faire des activités sociales, intellectuelles ou physiques. Les réponses doivent être concise.",
		},
	}

	pending_msg := ""
	c.SetReadLimit(5000)
	c.SetReadDeadline(time.Now().Add(120 * time.Second))
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error : ", err)	
			}
			// save messages into DB
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		log.Printf("rcv : %s\n", message)

		myMessage := string(message)
		if _, err := file.WriteString("user\n" + myMessage + "\n\n"); err != nil {
			log.Println("Erreur ecriture dans le fichier : ", err)
		}

		messages = append(messages, api.Message{
			Role: "user",
			Content: string(message),
		})

		req := &api.ChatRequest{
			Model: "llama3",
			Messages: messages,
		}
		err = client.Chat(context.Background(), req, func(m api.ChatResponse) error {
			err = c.WriteMessage(mt, []byte(m.Message.Content))
			if err != nil { log.Fatal(err) }
			pending_msg += m.Message.Content
			if m.Done {
				messages = append(messages, api.Message{
					Role: "assistant",
					Content: pending_msg,
				})
				pending_msg = ""
			}
			return nil
		})
	}
}

