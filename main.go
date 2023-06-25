package main

import "websocket-go-gin/server"

func main() {
	server.TodoList = make(map[string][]string)
	server.Run()
}
