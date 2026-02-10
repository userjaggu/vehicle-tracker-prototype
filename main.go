// Vehicle Tracker Prototype
//
// A minimal backend that accepts live GPS updates from vehicles
// and exposes the latest known positions via a simple HTTP API.
//
// Usage:
//
//	go run main.go
package main

import (
	"log"

	"github.com/jaggu/vehicle-tracker-prototype/server"
)

const defaultPort = 8080

func main() {
	if err := server.Run(defaultPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
