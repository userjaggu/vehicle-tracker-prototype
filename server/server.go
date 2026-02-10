// Package server wires together the HTTP routes and starts the server.
package server

import (
	"fmt"
	"net/http"

	"github.com/jaggu/vehicle-tracker-prototype/handler"
	"github.com/jaggu/vehicle-tracker-prototype/store"
)

// Run starts the HTTP server on the given port.
//
// It creates the in-memory store, registers routes, and blocks
// until the server is shut down or encounters a fatal error.
func Run(port int) error {

	// ── Create shared in-memory store ───────────────────────────────
	s := store.New()

	// ── Register routes ─────────────────────────────────────────────
	mux := http.NewServeMux()
	mux.HandleFunc("/location", handler.PostLocation(s))
	mux.HandleFunc("/vehicles", handler.GetVehicles(s))

	// ── Start listening ─────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Vehicle Tracker server listening on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}
