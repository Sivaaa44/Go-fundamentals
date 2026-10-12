package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

// GLOBAL STORE
var StoreMap = make(map[string]string)
var mu sync.Mutex

// var store = SharedStore({StoreMap, mu})

func clientRoutine(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			return
		}
		fmt.Println(line)
	}
}

func createServer() {
	listen, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("error listening to port 6379 :", err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("error accepting client connection to port", err)
		}
		go clientRoutine(conn)
	}
}

func main() {
	createServer()
	fmt.Println("main function")
}
