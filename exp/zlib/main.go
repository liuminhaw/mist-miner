package main

import (
	"fmt"

	"github.com/liuminhaw/mist-miner/shelf"
)

func main() {
	stuff := shelf.Stuff{
		Hash:     "",
		Module:   "",
		Identity: "",
	}

	if err := stuff.Read(); err != nil {
		fmt.Printf("Error reading stuff: %s\n", err)
	}
	fmt.Printf("Stuff resource: %s\n", string(stuff.Resource))
}
