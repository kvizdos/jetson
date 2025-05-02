package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <file> <key> <value>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	searchKey := os.Args[2]
	searchValue := os.Args[3]

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	matches := 0
	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		var obj map[string]interface{}
		if err := json.Unmarshal(line, &obj); err != nil {
			continue
		}

		if val, ok := obj[searchKey]; ok {
			if str, ok := val.(string); ok && strings.Contains(str, searchValue) {
				matches++
			}
		}
	}

	fmt.Printf("Naive scan found %d matches\n", matches)
}
