package server

import (
	"github.com/gin-gonic/gin"
)

type Todo struct {
	UserID string   `json:"userID"`
	Task   []string `json:"task"`
}

func getAllTodos(c *gin.Context) {
	var todos []Todo
	for userID, list := range TodoList {
		todos = append(todos, Todo{UserID: userID, Task: list})
	}
	c.JSON(200, todos)
}

func Run() {
	bindAddress := "localhost:8448"
	r := gin.Default()
	r.GET("/todo", ws)           // websocket endpoint
	r.GET("/todos", getAllTodos) // REST endpoint
	r.Run(bindAddress)
}

func init() {
	TodoList = make(map[string][]string)
}
