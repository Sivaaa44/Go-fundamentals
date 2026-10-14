package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var store = NewStore()
var aof *AOF // nil until persistence is opened

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
		if _, err := conn.Write([]byte(reply)); err != nil {
			return // client is gone
		}
	}
}

func createServer(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Println("listening on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go clientRoutine(conn)
	}
}

func main() {
	// 1. rebuild state from disk
	if err := replay("appendonly.aof"); err != nil {
		log.Fatal(err)
	}

	// 2. now start logging new writes
	var err error
	aof, err = OpenAOF("appendonly.aof")
	if err != nil {
		log.Fatal(err)
	}

	// 3. background expiry
	store.StartSweeper(time.Second)

	// 4. Ctrl+C: close the file cleanly, then exit
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		fmt.Println("\nshutting down")
		aof.Close()
		os.Exit(0)
	}()

	// 5. serve
	if err := createServer(":6379"); err != nil {
		log.Fatal(err)
	}
}
