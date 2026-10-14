package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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

	args := make([]string, 0, count)
	for i := 0; i < count; i++ {
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}

		if len(line) == 0 || line[0] != '$' {
			return nil, fmt.Errorf("protocol error: expected '$', got %q", line)
		}

		n, err := strconv.Atoi(line[1:])
		if err != nil || n < 0 || n > 512*1024*1024 {
			return nil, fmt.Errorf("protocol error: bad bulk length %q", line)
		}

		// read exactly n bytes of data + the trailing \r\n
		buf := make([]byte, n+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}

		if buf[n] != '\r' || buf[n+1] != '\n' {
			return nil, fmt.Errorf("protocol error: bulk string not terminated by CRLF")
		}
		args = append(args, string(buf[:n]))
	}
	return args, nil
}
