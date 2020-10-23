package main

import (
	"frontend"
	"log"
)

func main() {
	if err := frontend.Run(); err != nil {
		log.Fatalf("run frontend server: %v\n", err)
	}
}
