package server

import "strings"

var TodoList map[string][]string

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

func updateTodoList(userID, task string) {
	tmpList := TodoList[userID]
	TodoList[userID] = []string{}
	for _, val := range tmpList {
		if val == task {
			continue
		}
		TodoList[userID] = append(TodoList[userID], val)
	}
}
