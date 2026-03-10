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

	// Create shared in-memory store
	s := store.New()

	// Register routes
	mux := http.NewServeMux()

	// --- Driver-facing endpoints ---
	mux.HandleFunc("/location", handler.PostLocation(s))        // legacy endpoint
	mux.HandleFunc("/api/v1/locations", handler.PostLocation(s)) // matches mentor spec

	// --- Driver profile management ---
	mux.HandleFunc("GET /api/v1/drivers", handler.GetDrivers(s))
	mux.HandleFunc("POST /api/v1/drivers", handler.PostDriver(s))
	mux.HandleFunc("GET /api/v1/drivers/{id}", handler.GetDriver(s))
	mux.HandleFunc("PUT /api/v1/drivers/{id}", handler.PutDriver(s))
	mux.HandleFunc("DELETE /api/v1/drivers/{id}", handler.DeleteDriver(s))

	// --- GTFS-RT feed ---
	mux.HandleFunc("/gtfs-rt/vehicle-positions", handler.GetGTFSRT(s))

	// --- Operational endpoints ---
	mux.HandleFunc("/vehicles", handler.GetVehicles(s))
	mux.HandleFunc("/api/v1/status", handler.GetStatus(s))

	// Start listening
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Vehicle Tracker server listening on http://localhost%s\n", addr)
	fmt.Printf("  POST /api/v1/locations           — submit vehicle GPS data\n")
	fmt.Printf("  GET  /gtfs-rt/vehicle-positions   — GTFS-RT protobuf feed\n")
	fmt.Printf("  GET  /gtfs-rt/vehicle-positions?format=json — feed as JSON\n")
	fmt.Printf("  GET  /vehicles                    — all vehicle locations\n")
	fmt.Printf("  GET  /api/v1/status               — system health\n")
	fmt.Printf("  GET  /api/v1/drivers              — list driver profiles\n")
	fmt.Printf("  POST /api/v1/drivers              — register a driver\n")
	fmt.Printf("  GET  /api/v1/drivers/{id}         — get a driver profile\n")
	fmt.Printf("  PUT  /api/v1/drivers/{id}         — update a driver profile\n")
	fmt.Printf("  DELETE /api/v1/drivers/{id}       — remove a driver profile\n")
	return http.ListenAndServe(addr, mux)
}
