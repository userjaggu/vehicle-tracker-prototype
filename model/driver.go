// Package model defines the data structures used across the application.
package model

// Driver represents a transit vehicle driver's user profile.
//
// A driver is associated with zero or one vehicle at a time.
// The VehicleID field is optional — it is set when the driver is
// actively operating a vehicle and cleared when the trip ends.
type Driver struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	VehicleID string `json:"vehicle_id,omitempty"`
}
