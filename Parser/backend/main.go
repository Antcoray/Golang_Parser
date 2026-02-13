package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {

	// dictionary := map[string]int{
	// 	"==": 0,
	// 	"!=": 0,
	// }

	filepath := filepath.Join("..", "..", "ExampleCode", "main.go")

	file, error := os.Open(filepath)

	if error != nil {
		log.Fatal(error)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
