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

func handleCmd(userID, cmd, msg string, ws *websocket.Conn) {
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
	case "close":
		// Send a close message before closing the connection
		err := ws.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.CloseNormalClosure,
				"Closing the connection upon user request"),
		)
		if err != nil {
			log.Println("write close:", err)
		}
		return
	default:
		log.Println("Invalid command")
	}
}

func renderTodoList(userID string) []byte {
	// Create a string with all the todos
	output := "Current Todos: \n\n"
	for i, todo := range TodoList[userID] {
		output += fmt.Sprintf("%d. %s\n", i+1, todo)
	}
	output += "\n----------------------------------------"
	return []byte(output)
}

func handleConnection(ws *websocket.Conn) {
	defer ws.Close()
	clientID := ws.RemoteAddr()
	userID := fmt.Sprintf("%s", clientID)
	for {
		//read data from ws
		log.Printf("Waiting for input from %s", clientID)
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s from %s", message, clientID)
		input := string(message)
		cmd := getCmd(input)
		msg := getMessage(input)
		handleCmd(userID, cmd, msg, ws)
		resp := renderTodoList(userID)
		//write ws data
		err = ws.WriteMessage(mt, resp)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func ws(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	go handleConnection(ws)
}

func Run() {
	bindAddress := "localhost:8448"
	r := gin.Default()
	r.GET("/todo", ws)
	r.Run(bindAddress)
}
