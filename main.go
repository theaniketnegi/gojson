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

	parser, err := NewParser(reader)
	if err != nil {
		log.Fatal("Error creating parser: ", err)
	}
	value, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
