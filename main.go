package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

var todoList []string

func getCmd(input string) string {
	inputArr := strings.Split(input, " ")
	if len(inputArr) == 0 {
		return ""
	}
	return inputArr[0]
}

func getMessage(input string) string {
	inputArr := strings.Split(input, " ")
	return strings.Join(inputArr[1:], " ")
}

func updateTodoList(input string) {
	tmpList := todoList
	todoList = []string{}
	for _, val := range tmpList {
		if val == input {
			continue
		}
		todoList = append(todoList, val)
	}
}

func ws(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()
	for {
		//read data from ws
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		input := string(message)
		log.Printf("recv: %s", message)
		cmd := getCmd(input)
		msg := getMessage(input)
		switch cmd {
		case "add":
			if msg == "" {
				log.Println("Empty task...")
				break
			}
			todoList = append(todoList, msg)
		case "done":
			updateTodoList(msg)
		case "clear":
			todoList = []string{}
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
		output := "Current Todos: \n\n"
		for i, todo := range todoList {
			output += fmt.Sprintf("%d. %s\n", i+1, todo)
		}
		output += "\n----------------------------------------"
		message = []byte(output)

		//write ws data
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	bindAddress := "localhost:8448"
	r := gin.Default()
	r.GET("/todo", ws)
	r.Run(bindAddress)
}
