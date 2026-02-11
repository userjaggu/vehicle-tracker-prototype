// Package model defines the data structures used across the application.
package model

// Location represents a single GPS update from a vehicle.
//
// Fields:
//   VehicleID:  unique identifier for the vehicle (e.g. "bus-42")
//    Latitude:   GPS latitude in decimal degrees
//    Longitude:  GPS longitude in decimal degrees
//   Timestamp:  Unix timestamp of when the location was recorded
type Location struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}
