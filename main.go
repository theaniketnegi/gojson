package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Provide a file name")
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lexer := InitLexer(reader)

	tokens, err := lexer.Tokenize()

	if err != nil {
		log.Fatal(err)
	}

	for _, token := range tokens {
		fmt.Printf("Key: %v, Value: %v\n", token.Key, token.Value)
	}
}
