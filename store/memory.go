// Package store provides in-memory storage for vehicle locations.
//
// Design decisions:
//
//	Uses a sync.RWMutex for safe concurrent access.
//	Each vehicle's entry is overwritten on every update (latest-only).
//	Supports staleness filtering for the GTFS-RT feed.
//	No persistence — data is lost when the process exits.
package store

import (
	"sync"
	"time"

	"github.com/jaggu/vehicle-tracker-prototype/model"
)

// MemoryStore holds the latest known location for each vehicle.
type MemoryStore struct {
	mu        sync.RWMutex
	locations map[string]model.Location

	// receivedAt tracks the server wall-clock time each location was
	// stored, so we can apply staleness filtering independently of the
	// client-supplied timestamp.
	receivedAt map[string]time.Time
}

// New creates and returns an empty MemoryStore.
func New() *MemoryStore {
	return &MemoryStore{
		locations:  make(map[string]model.Location),
		receivedAt: make(map[string]time.Time),
	}
}

// UpdateLocation stores (or overwrites) the latest location for a vehicle.
func (s *MemoryStore) UpdateLocation(loc model.Location) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.locations[loc.VehicleID] = loc
	s.receivedAt[loc.VehicleID] = time.Now()
}

// GetAllLocations returns a snapshot of all known vehicle locations.
func (s *MemoryStore) GetAllLocations() []model.Location {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Location, 0, len(s.locations))
	for _, loc := range s.locations {
		result = append(result, loc)
	}
	return result
}

// GetActiveLocations returns locations that were received within the given
// staleness window.  Vehicles that haven't reported in longer than the
// threshold are considered inactive and excluded from the result.
func (s *MemoryStore) GetActiveLocations(threshold time.Duration) []model.Location {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := time.Now().Add(-threshold)
	result := make([]model.Location, 0, len(s.locations))
	for id, loc := range s.locations {
		if s.receivedAt[id].After(cutoff) {
			result = append(result, loc)
		}
	}
	return result
}

// ActiveVehicleCount returns the number of vehicles that have reported
// within the given staleness window.
func (s *MemoryStore) ActiveVehicleCount(threshold time.Duration) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := time.Now().Add(-threshold)
	count := 0
	for id := range s.locations {
		if s.receivedAt[id].After(cutoff) {
			count++
		}
	}
	return count
}

// TotalVehicleCount returns the total number of vehicles that have ever
// reported a location.
func (s *MemoryStore) TotalVehicleCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.locations)
}
