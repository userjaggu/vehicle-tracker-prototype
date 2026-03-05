package store_test

import (
	"testing"
	"time"

	"github.com/jaggu/vehicle-tracker-prototype/model"
	"github.com/jaggu/vehicle-tracker-prototype/store"
)

func TestMemoryStore_UpdateAndRetrieve(t *testing.T) {
	s := store.New()

	loc := model.Location{
		VehicleID: "bus-1",
		Latitude:  17.385,
		Longitude: 78.4867,
		Timestamp: 1707350000,
	}

	s.UpdateLocation(loc)
	all := s.GetAllLocations()

	if len(all) != 1 {
		t.Fatalf("expected 1 location, got %d", len(all))
	}
	if all[0].VehicleID != "bus-1" {
		t.Errorf("vehicle_id = %q, want %q", all[0].VehicleID, "bus-1")
	}
}

func TestMemoryStore_OverwriteLocation(t *testing.T) {
	s := store.New()

	s.UpdateLocation(model.Location{VehicleID: "bus-1", Latitude: 17.0, Longitude: 78.0})
	s.UpdateLocation(model.Location{VehicleID: "bus-1", Latitude: 18.0, Longitude: 79.0})

	all := s.GetAllLocations()
	if len(all) != 1 {
		t.Fatalf("expected 1 location after overwrite, got %d", len(all))
	}
	if all[0].Latitude != 18.0 {
		t.Errorf("latitude = %f, want 18.0 (should be overwritten)", all[0].Latitude)
	}
}

func TestMemoryStore_ActiveLocations(t *testing.T) {
	s := store.New()

	// This vehicle was just reported — it should be active
	s.UpdateLocation(model.Location{VehicleID: "bus-fresh", Latitude: 17.3, Longitude: 78.4})

	active := s.GetActiveLocations(5 * time.Minute)
	if len(active) != 1 {
		t.Fatalf("expected 1 active location, got %d", len(active))
	}

	// With a zero threshold, even a just-reported vehicle is stale
	stale := s.GetActiveLocations(0)
	if len(stale) != 0 {
		t.Errorf("expected 0 locations with zero threshold, got %d", len(stale))
	}
}

func TestMemoryStore_ActiveVehicleCount(t *testing.T) {
	s := store.New()

	s.UpdateLocation(model.Location{VehicleID: "bus-1", Latitude: 17.0, Longitude: 78.0})
	s.UpdateLocation(model.Location{VehicleID: "bus-2", Latitude: 17.1, Longitude: 78.1})

	if got := s.ActiveVehicleCount(5 * time.Minute); got != 2 {
		t.Errorf("active vehicle count = %d, want 2", got)
	}
	if got := s.TotalVehicleCount(); got != 2 {
		t.Errorf("total vehicle count = %d, want 2", got)
	}
}

func TestMemoryStore_ExpandedFields(t *testing.T) {
	s := store.New()

	loc := model.Location{
		VehicleID: "bus-42",
		TripID:    "route_5_0830",
		RouteID:   "5",
		Latitude:  -1.2921,
		Longitude: 36.8219,
		Bearing:   180.0,
		Speed:     8.5,
		Accuracy:  12.0,
		Timestamp: 1752566400,
	}

	s.UpdateLocation(loc)
	all := s.GetAllLocations()

	if all[0].TripID != "route_5_0830" {
		t.Errorf("trip_id = %q, want %q", all[0].TripID, "route_5_0830")
	}
	if all[0].Bearing != 180.0 {
		t.Errorf("bearing = %f, want 180.0", all[0].Bearing)
	}
	if all[0].Speed != 8.5 {
		t.Errorf("speed = %f, want 8.5", all[0].Speed)
	}
}
