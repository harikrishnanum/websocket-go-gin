package server

import "strings"

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
