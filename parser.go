package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// readLine reads one line and strips the trailing \r\n.
func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

// readCommand reads one RESP command, e.g. ["SET", "name", "siva"].
func readCommand(r *bufio.Reader) ([]string, error) {
	line, err := readLine(r)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 || line[0] != '*' {
		return nil, fmt.Errorf("protocol error: expected '*', got %q", line)
	}
	count, err := strconv.Atoi(line[1:])
	if err != nil || count < 1 {
		return nil, fmt.Errorf("protocol error: bad array length %q", line)
	}

	// 2. For each item: a "$len" line, then the data line
	args := make([]string, 0, count)
	for i := 0; i < count; i++ {
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}
		if len(line) == 0 || line[0] != '$' {
			return nil, fmt.Errorf("protocol error: expected '$', got %q", line)
		}
		if _, err := strconv.Atoi(line[1:]); err != nil {
			return nil, fmt.Errorf("protocol error: bad bulk length %q", line)
		}

		data, err := readLine(r)
		if err != nil {
			return nil, err
		}
		args = append(args, data)
	}
	return args, nil
}
