package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// GLOBAL STORE
var StoreMap = SharedStore{Store: make(map[string]string)}

// var store = SharedStore({StoreMap, mu})

func clientRoutine(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		cmd, err := readCommand(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Println("error:", err)
				conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
			}
			return
		}
		reply := handleCommand(cmd)
		conn.Write([]byte(reply))
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
