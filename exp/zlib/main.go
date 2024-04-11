package main

import (
	"fmt"

	"github.com/liuminhaw/mist-miner/shelf"
)

func main() {
	stuff := shelf.Stuff{
		Hash:     "2cbe35d5c8b0de9d95ce4bf80ec03b4661ae67f1beac75016620b81447e724d9",
		Module:   "mm-s3",
		Identity: "lmhaw",
	}

	if err := stuff.Read(); err != nil {
		fmt.Printf("Error reading stuff: %s\n", err)
	}
	fmt.Printf("Stuff resource: %s\n", string(stuff.Resource))
}
