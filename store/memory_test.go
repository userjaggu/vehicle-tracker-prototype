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

func TestMemoryStore_UpsertAndGetDriver(t *testing.T) {
	s := store.New()

	d := model.Driver{ID: "driver-1", Name: "Alice", Email: "alice@example.com"}

	created := s.UpsertDriver(d)
	if !created {
		t.Errorf("UpsertDriver should return true (created) on first insert")
	}

	got, ok := s.GetDriver("driver-1")
	if !ok {
		t.Fatal("GetDriver: expected to find driver-1")
	}
	if got.Name != "Alice" {
		t.Errorf("Name = %q, want %q", got.Name, "Alice")
	}
	if got.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", got.Email, "alice@example.com")
	}
}

func TestMemoryStore_UpdateDriver(t *testing.T) {
	s := store.New()

	s.UpsertDriver(model.Driver{ID: "driver-1", Name: "Alice"})

	// Update — should return false (not a new record)
	created := s.UpsertDriver(model.Driver{ID: "driver-1", Name: "Alice Updated", Phone: "555-1234"})
	if created {
		t.Errorf("UpsertDriver should return false (updated) when driver already exists")
	}

	got, _ := s.GetDriver("driver-1")
	if got.Name != "Alice Updated" {
		t.Errorf("Name after update = %q, want %q", got.Name, "Alice Updated")
	}
	if got.Phone != "555-1234" {
		t.Errorf("Phone after update = %q, want %q", got.Phone, "555-1234")
	}
}

func TestMemoryStore_GetAllDrivers(t *testing.T) {
	s := store.New()

	s.UpsertDriver(model.Driver{ID: "d1", Name: "Alice"})
	s.UpsertDriver(model.Driver{ID: "d2", Name: "Bob"})

	all := s.GetAllDrivers()
	if len(all) != 2 {
		t.Fatalf("expected 2 drivers, got %d", len(all))
	}
}

func TestMemoryStore_GetAllDrivers_Empty(t *testing.T) {
	s := store.New()

	all := s.GetAllDrivers()
	if len(all) != 0 {
		t.Errorf("expected 0 drivers on empty store, got %d", len(all))
	}
}

func TestMemoryStore_DeleteDriver(t *testing.T) {
	s := store.New()

	s.UpsertDriver(model.Driver{ID: "driver-1", Name: "Alice"})

	deleted := s.DeleteDriver("driver-1")
	if !deleted {
		t.Errorf("DeleteDriver should return true when driver exists")
	}

	_, ok := s.GetDriver("driver-1")
	if ok {
		t.Errorf("driver-1 should not be found after deletion")
	}
}

func TestMemoryStore_DeleteDriver_NotFound(t *testing.T) {
	s := store.New()

	deleted := s.DeleteDriver("nonexistent")
	if deleted {
		t.Errorf("DeleteDriver should return false when driver does not exist")
	}
}

func TestMemoryStore_GetDriver_NotFound(t *testing.T) {
	s := store.New()

	_, ok := s.GetDriver("nonexistent")
	if ok {
		t.Errorf("GetDriver should return false for nonexistent driver")
	}
}
