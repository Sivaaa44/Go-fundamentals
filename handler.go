// handler.go
package main

import (
	"fmt"
	"strings"
)

// handleCommand runs one parsed command and returns the RESP reply.
func handleCommand(args []string) string {
	switch strings.ToUpper(args[0]) {

	case "PING":
		return "+PONG\r\n"

	case "SET":
		if len(args) != 3 {
			return "-ERR wrong number of arguments for 'set'\r\n"
		}
		StoreMap.Set(args[1], args[2])
		return "+OK\r\n"

	case "GET":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for 'get'\r\n"
		}
		val, ok := StoreMap.Get(args[1])
		if !ok {
			return "$-1\r\n" // null: key doesn't exist
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)

	case "DEL":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for 'del'\r\n"
		}
		if StoreMap.Del(args[1]) {
			return ":1\r\n"
		}
		return ":0\r\n"
	}

	return fmt.Sprintf("-ERR unknown command '%s'\r\n", args[0])
}
