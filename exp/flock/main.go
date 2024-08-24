package main

import (
	"log"
	"os"
	"time"

	"github.com/gofrs/flock"
)

func main() {
	filePath := "example.txt"
	lock := flock.New(filePath + ".lock")

	locked, err := lock.TryLock()
	if err != nil {
		log.Fatalf("Error acquiring lock: %v", err)
	}

	if !locked {
		log.Println("File is already locked, another process is writing to it.")
		return
	}
	defer lock.Unlock()

	// Simulate long-running operation
	time.Sleep(5 * time.Second)

	// file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("Writing to the file...")
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	log.Println("File written successfully.")
}
