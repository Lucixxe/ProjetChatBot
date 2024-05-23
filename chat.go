package main

import (
	"net/http"
	"log"
	"time"
	"bytes"

	"github.com/gorilla/websocket"
)

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
	c.SetReadLimit(2000)
	c.SetReadDeadline(time.Now().Add(60 * time.Second))
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
		err = c.WriteMessage(mt, append([]byte("i received the message : "), message...))
		if err != nil {
			log.Println("write :", err)
			break
		}
	}
}


