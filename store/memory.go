// Package store provides in-memory storage for vehicle locations.
//
// Design decisions:
//   - Uses a sync.RWMutex for safe concurrent access.
//   - Each vehicle's entry is overwritten on every update (latest-only).
//   - No persistence — data is lost when the process exits.
package store

import (
	"sync"

	"github.com/jaggu/vehicle-tracker-prototype/model"
)

// MemoryStore holds the latest known location for each vehicle.
type MemoryStore struct {
	mu        sync.RWMutex
	locations map[string]model.Location
}

// New creates and returns an empty MemoryStore.
func New() *MemoryStore {
	return &MemoryStore{
		locations: make(map[string]model.Location),
	}
}

// UpdateLocation stores (or overwrites) the latest location for a vehicle.
func (s *MemoryStore) UpdateLocation(loc model.Location) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.locations[loc.VehicleID] = loc
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
