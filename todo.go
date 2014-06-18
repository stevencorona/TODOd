package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

const VERSION = "1.0.0"

type TodoItem struct {
	Task string
}

type TodoList struct {
	Items []TodoItem
	sync.RWMutex
}

func main() {
	log.Println("Welcome to Todolist Server", VERSION)
	todoList := new(TodoList)
	log.Println("You have", len(todoList.Items), "in your todolist")

	server, err := net.Listen("tcp", "127.0.0.1:9000")

	if err != nil {
		log.Fatal("Could not create the server")
	}

	defer server.Close()

	taskChannel := make(chan int, 100)

	go runAnalytics(todoList, taskChannel)

	for {
		conn, _ := server.Accept()
		go handleConnection(conn, todoList, taskChannel)
	}

}

func runAnalytics(todoList *TodoList, taskChannel chan int) {
	for {
		<-taskChannel
		log.Println("Tasks: ", len(todoList.Items))
	}
}

func handleConnection(conn net.Conn, todoList *TodoList, taskChannel chan int) {
	fmt.Fprintf(conn, "Welcome to Todolist Server\n")

	reader := bufio.NewReader(conn)

	for {
		fmt.Fprintf(conn, "%d items in your TodoList\n", len(todoList.Items))
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Println("Client did something weird")
			return
		}

		line = strings.TrimSpace(line)

		switch line {
		case "add":
			task, _ := reader.ReadString('\n')
			todoList.AddTask(task)
			taskChannel <- 1
		case "list":
			json, _ := json.Marshal(todoList)
			fmt.Fprintf(conn, "%s\n", json)
		}
	}
}

func (todoList *TodoList) AddTask(task string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Phew, saved the day.")
		}
	}()

	task = strings.TrimSpace(task)
	todoList.Lock()
	defer todoList.Unlock()

	if task == "themes" {
		panic("Ack!")
	}

	item := TodoItem{task}
	todoList.Items = append(todoList.Items, item)
}
