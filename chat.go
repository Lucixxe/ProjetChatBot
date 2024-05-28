package main

// communication websocket

import (
	"net/http"
	"log"
	"time"
	"bytes"
	"html/template"

	"github.com/gorilla/websocket"
)

func chat (w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("pages/chat.tmpl"))
	err := tmp.ExecuteTemplate(w, "chat", struct{
		Pseudo	string
	}{
		"Moi",
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
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade : ", err)
	}
	defer c.Close()

	c.SetReadLimit(5000)
	c.SetReadDeadline(time.Now().Add(120 * time.Second))
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
		err = c.WriteMessage(mt, []byte("salut"))
	}
}

