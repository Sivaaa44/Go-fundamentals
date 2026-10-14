package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type AOF struct {
	mu sync.Mutex
	f  *os.File
}

func OpenAOF(path string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &AOF{f: f}, nil
}

// Append writes the command to the file in RESP format,
// the same bytes a client would have sent.
func (a *AOF) Append(args []string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(args))
	for _, s := range args {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(s), s)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	_, err := a.f.WriteString(b.String())
	return err
}

func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.f.Close()
}

// replay re-runs every command in the file to rebuild the store.
// Call it BEFORE the global aof is set, so replayed commands aren't logged again.
func replay(path string) error {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	count := 0
	for {
		cmd, err := readCommand(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("aof: stopping replay:", err)
			break
		}
		handleCommand(cmd)
		count++
	}
	fmt.Println("aof: replayed", count, "commands")
	return nil
}
