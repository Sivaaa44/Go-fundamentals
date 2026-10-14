package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// logWrite appends a successful write command to the AOF (if enabled).
func logWrite(args []string) {
	if aof != nil {
		if err := aof.Append(args); err != nil {
			fmt.Println("aof write error:", err)
		}
	}
}

// handleCommand runs one parsed command and returns the RESP reply.
func handleCommand(args []string) string {
	switch strings.ToUpper(args[0]) {

	case "PING":
		return "+PONG\r\n"

	case "SET":
		// SET key value [EX seconds]
		if len(args) != 3 && len(args) != 5 {
			return "-ERR wrong number of arguments for 'set'\r\n"
		}
		var ttl time.Duration
		if len(args) == 5 {
			if strings.ToUpper(args[3]) != "EX" {
				return "-ERR syntax error\r\n"
			}
			secs, err := strconv.Atoi(args[4])
			if err != nil || secs <= 0 {
				return "-ERR invalid expire time in 'set' command\r\n"
			}
			ttl = time.Duration(secs) * time.Second
		}
		store.Set(args[1], args[2], ttl)
		logWrite(args)
		return "+OK\r\n"

	case "GET":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for 'get'\r\n"
		}
		val, ok := store.Get(args[1])
		if !ok {
			return "$-1\r\n"
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)

	case "DEL":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for 'del'\r\n"
		}
		if store.Del(args[1]) {
			logWrite(args)
			return ":1\r\n"
		}
		return ":0\r\n"

	case "EXISTS":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for 'exists'\r\n"
		}
		if _, ok := store.Get(args[1]); ok {
			return ":1\r\n"
		}
		return ":0\r\n"

	case "TTL":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for 'ttl'\r\n"
		}
		return fmt.Sprintf(":%d\r\n", store.TTL(args[1]))
	}

	return fmt.Sprintf("-ERR unknown command '%s'\r\n", args[0])
}
