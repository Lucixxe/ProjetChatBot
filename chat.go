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
	"math/rand"

	"github.com/gorilla/websocket"
	"github.com/ollama/ollama/api"
)

type Message struct {
	Contenu string `json:"contenu"`
	Kind    string `json:"type"`
}

/*
chat gère la requête HTTP pour afficher la page de chat.

Paramètres : - w : l'interface utilisée par un gestionnaire HTTP pour construire une réponse HTTP
    		 - r : la requête HTTP entrante reçue par le serveur
*/
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

/*
handleError gère les erreurs en enregistrant un message d'erreur 
personnalisé et en terminant le programme si une erreur est présente.

Paramètres: - err: l'erreur à traiter
		    - errorMessage: un message selon l'erreur
*/
func handleError(err error, errorMessage string) {
	if err != nil {
		log.Fatalf("%s: %v", errorMessage, err)
	}
}

/*
sendHistory envoie l'historique des messages au client via une connexion WebSocket.

Paramètres: - c: la connexion WebSocket
			- messages_for_history: la liste des messages à envoyer

Retour: - error: renvoie une erreur si la sérialisation ou l'envoi échoue
*/
func sendHistory(c *websocket.Conn, messages_for_history []JSONMessage) error {
	initialMessages, err := json.Marshal(messages_for_history)
	handleError(err, "Error marshalling messages")

	return c.WriteMessage(websocket.TextMessage, initialMessages)
}

/*
selectWelcomeMessage sélectionne un message de bienvenue aléatoire selon que l'utilisateur
se connecte pour la première fois ou non.

Paramètres: - isFirstConnection: un booléen qui indique si c'est la première connexion de 
              					 l'utilisateur

Retour: - string: un message de bienvenue aléatoire.
*/
func selectWelcomeMessage(isFirstConnection bool) string {
	welcomeMessages_newUser := []string{
		"Bonjour, heureux de vous rencontrer. Comment puis-je vous aider aujourd'hui ?",
		"Bienvenue ! En quoi puis-je vous aider ?",
		"Bonjour, enchanté de faire votre connaissance. Que puis-je faire pour vous ?",
	}

	welcomeMessages_existingUser := []string{
		"Bonjour, heureux de vous revoir ! Comment puis-je vous aider aujourd'hui ?",
		"Bonjour, je suis ravi de vous retrouver ! Comment puis-je vous assister aujourd'hui ?",
		"Bonjour, c'est un plaisir de vous revoir ! En quoi puis-je vous être utile aujourd'hui ?",
	}

	rand.Seed(time.Now().UnixNano())
	if isFirstConnection {
        index := rand.Intn(len(welcomeMessages_newUser))
		return welcomeMessages_newUser[index]
	}
    index := rand.Intn(len(welcomeMessages_existingUser))
    return welcomeMessages_existingUser[index]
}

/*
sendWelcomeMessage envoie un message de bienvenue au client via une connexion WebSocket.

Paramètres: - c: la connexion WebSocket
			- message: le message de bienvenue à envoyer

Retour: - error: renvoie une erreur si la sérialisation ou l'envoi échoue
*/
func sendWelcomeMessage(c *websocket.Conn, message string) error {
	welcomeMSG, err := json.Marshal(&Message{
		message,
		"message",
	})
	handleError(err, "Error marshalling message")

	return c.WriteMessage(websocket.TextMessage, welcomeMSG)
}

/*
sendMessagesToClient envoie un message au client via une connexion WebSocket.

Paramètres: - c: la connexion WebSocket
			- messageContent: le contenu du message à envoyer
			- messageType: le type de message à envoyer

Retour: - error: renvoie une erreur si la sérialisation ou l'envoi échoue
*/
func sendMessagesToClient(c *websocket.Conn, messageContent string, messageType string) error {
	//TODO: check le type du message
	initialMessages, err := json.Marshal(&Message{
		messageContent,
	    messageType,
	})
	handleError(err, "Error marshalling message")

	return c.WriteMessage(websocket.TextMessage, initialMessages)
}

/*
handleClientCommunication gère l'échange de messages via une connexion WebSocket entre 
le client et le serveur.

Paramètres: - c: la connexion WebSocket
			- client: le client qui s'est connecté
			- messages: une liste contenant les messages envoyés

Valeur de retour: - error: renvoie une erreur si la communication avec le client échoue
*/
func handleClientCommunication(c *websocket.Conn, client *api.Client, messages []api.Message) error {
	pending_msg := ""
	c.SetReadLimit(5000)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error : ", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		//log.Printf("rcv : %s\n", message)

		// Sauvegarde des messages dans la base de données
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
			err = sendMessagesToClient(c, m.Message.Content, "message")
			handleError(err, "Error sending message to client")

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

				err = sendMessagesToClient(c, "#fin#", "message")	
				handleError(err, "Error sending message to client")
			}

			return nil
		})

		handleError(err, "Error sending message to client")
	}

	return nil
}


var upgrader = websocket.Upgrader{}

/*
ws_con gère la communication WebSocket entre le client et le serveur.

Paramètres: - w : l'interface utilisée par un gestionnaire HTTP pour construire une réponse HTTP
    		- r : la requête HTTP entrante reçue par le serveur
*/
func ws_con(w http.ResponseWriter, r *http.Request) {
	client, err := api.ClientFromEnvironment()
	handleError(err, "Error getting client")

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
	handleError(err, "Error loading messages from database")
	
	// Envoi des messages de l'historique dans la websocket
	err = sendHistory(c, messages_for_history)
	handleError(err, "Error sending history")

	// Check si l'utilisateur se connecte pour la première fois ou non
	isFirstConnection := false
	if len(messages_for_history) == 1 && messages_for_history[0].Content == "New user" {
		isFirstConnection = true
	}

	// Séléction du message en conséquence et envoie dans la websocket
	welcomeMessage := selectWelcomeMessage(isFirstConnection)
	err = sendWelcomeMessage(c, welcomeMessage)
	handleError(err, "Error sending welcome message")

	saveMessage(user.id, "user", time.Now().Format("02/01/2006 15:04"), welcomeMessage)

	err = sendMessagesToClient(c, "#fin#", "message")	
	handleError(err, "Error sending message to client")

	// "début" de la conversation
	err = handleClientCommunication(c, client, messages)
	handleError(err, "Error handling client communication")
}
