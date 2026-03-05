// Package model defines the data structures used across the application.
package model

import "time"

// Location represents a single GPS update from a vehicle.
//
// Fields align with the GTFS-RT VehiclePosition specification so that the
// in-memory state can be mapped directly to a protobuf FeedMessage.
type Location struct {
	VehicleID string  `json:"vehicle_id"`
	TripID    string  `json:"trip_id,omitempty"`
	RouteID   string  `json:"route_id,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Bearing   float32 `json:"bearing,omitempty"`
	Speed     float32 `json:"speed,omitempty"`
	Accuracy  float32 `json:"accuracy,omitempty"`
	Timestamp int64   `json:"timestamp"`
}

// DefaultStalenessThreshold is the maximum age of a location report before
// it is excluded from the GTFS-RT feed.  A vehicle that hasn't reported in
// this window is assumed to no longer be active.
const DefaultStalenessThreshold = 5 * time.Minute
