package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	/* setting CheckOrigin to always return true can have security implications,
	 * as it allows WebSocket connections from any origin.
	 * In a production environment, implement a proper origin check
	 * based on your specific requirements and security considerations.
	 */
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleCmd(userID, cmd, msg string, ws *websocket.Conn) bool {
	// return true if the connection should be closed
	switch cmd {
	case "add":
		if msg == "" {
			log.Println("Empty task...")
			break
		}
		TodoList[userID] = append(TodoList[userID], msg)
	case "done":
		updateTodoList(userID, msg)
	case "clear":
		TodoList[userID] = []string{}
	case "list":
		log.Printf("Sending todo list to %s", userID)
	case "close":
		log.Printf("Closing connection for %s upon user request", userID)
		// Send a close message before closing the connection
		err := ws.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.CloseNormalClosure,
				"Closing the connection upon user request"),
		)
		if err != nil {
			log.Println(err.Error())
		}
		delete(TodoList, userID)
		return true
	default:
		log.Println("Invalid command")
		return false
	}
	resp := renderTodoList(userID)
	err := ws.WriteMessage(websocket.TextMessage, resp)
	if err != nil {
		log.Println(err.Error())
	}
	return false
}

func handleConnection(ws *websocket.Conn) {
	defer ws.Close()
	clientID := ws.RemoteAddr()
	userID := fmt.Sprintf("%s", clientID)
	for {
		//read data from ws
		log.Printf("Waiting for input from %s", clientID)
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s from %s", message, clientID)
		input := string(message)
		cmd := getCmd(input)
		msg := getMessage(input)
		isClose := handleCmd(userID, cmd, msg, ws)
		if isClose {
			break
		}
	}
	log.Printf("Closed connection for %s", clientID)
}

func ws(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	log.Println("Upgrading to websockets...")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	go handleConnection(ws)
}
