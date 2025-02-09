package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./tests/step4/valid.json")
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Value: %s\n", err)
				return
			}
			log.Fatal("Error reading file: ", err)
		}

		if b == '\n' {
			fmt.Printf("Value: \\n\n")
		} else if b == ' ' {
			fmt.Printf("Value: <space>\n")
		} else {
			fmt.Printf("Value: %s\n", string(b))
		}
	}
}
