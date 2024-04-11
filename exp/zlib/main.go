package main

import (
	"fmt"

	"github.com/liuminhaw/mist-miner/shelf"
)

func main() {
	stuff := shelf.Stuff{
		Hash:     "5f0c8365480fe959806dff9b8d08fc85e7a454523fbd8a25f184ac5c60564dc4",
		Module:   "mm-s3",
		Identity: "lmhaw",
	}

    if err := stuff.Read(); err != nil {
       fmt.Printf("Error reading stuff: %s\n", err) 
    }
	fmt.Printf("Stuff resource: %s\n", string(stuff.Resource))
}
