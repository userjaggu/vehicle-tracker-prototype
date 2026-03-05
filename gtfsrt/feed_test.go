package gtfsrt_test

import (
	"testing"
	"time"

	"github.com/jaggu/vehicle-tracker-prototype/gtfsrt"
	"github.com/jaggu/vehicle-tracker-prototype/model"
	pb "github.com/jaggu/vehicle-tracker-prototype/proto/gtfsrt"
	"google.golang.org/protobuf/proto"
)

// TestBuildFeed_EmptyInput verifies that an empty location slice produces
// a valid FeedMessage with zero entities (this is the correct behavior —
// a feed with no active vehicles is still valid GTFS-RT).
func TestBuildFeed_EmptyInput(t *testing.T) {
	feed := gtfsrt.BuildFeed([]model.Location{})

	if feed.Header == nil {
		t.Fatal("FeedHeader must not be nil")
	}
	if v := feed.Header.GetGtfsRealtimeVersion(); v != "2.0" {
		t.Errorf("expected version 2.0, got %s", v)
	}
	if inc := feed.Header.GetIncrementality(); inc != pb.FeedHeader_FULL_DATASET {
		t.Errorf("expected FULL_DATASET, got %v", inc)
	}
	if len(feed.Entity) != 0 {
		t.Errorf("expected 0 entities, got %d", len(feed.Entity))
	}
}

// TestBuildFeed_SingleVehicle verifies the happy-path for one vehicle
// with all fields populated.
func TestBuildFeed_SingleVehicle(t *testing.T) {
	locations := []model.Location{
		{
			VehicleID: "bus-42",
			TripID:    "route_5_0830",
			RouteID:   "5",
			Latitude:  -1.2921,
			Longitude: 36.8219,
			Bearing:   180.0,
			Speed:     8.5,
			Accuracy:  12.0,
			Timestamp: 1752566400,
		},
	}

	feed := gtfsrt.BuildFeed(locations)

	if len(feed.Entity) != 1 {
		t.Fatalf("expected 1 entity, got %d", len(feed.Entity))
	}

	entity := feed.Entity[0]
	if entity.GetId() != "vehicle-bus-42" {
		t.Errorf("entity id = %q, want %q", entity.GetId(), "vehicle-bus-42")
	}

	vp := entity.GetVehicle()
	if vp == nil {
		t.Fatal("VehiclePosition must not be nil")
	}

	// Check position
	pos := vp.GetPosition()
	if pos == nil {
		t.Fatal("Position must not be nil")
	}
	if got := pos.GetLatitude(); got < -1.293 || got > -1.291 {
		t.Errorf("latitude = %f, want ~-1.2921", got)
	}
	if got := pos.GetLongitude(); got < 36.821 || got > 36.823 {
		t.Errorf("longitude = %f, want ~36.8219", got)
	}
	if got := pos.GetBearing(); got != 180.0 {
		t.Errorf("bearing = %f, want 180.0", got)
	}
	if got := pos.GetSpeed(); got != 8.5 {
		t.Errorf("speed = %f, want 8.5", got)
	}

	// Check vehicle descriptor
	vd := vp.GetVehicle()
	if vd.GetId() != "bus-42" {
		t.Errorf("vehicle.id = %q, want %q", vd.GetId(), "bus-42")
	}

	// Check trip descriptor
	td := vp.GetTrip()
	if td == nil {
		t.Fatal("TripDescriptor must not be nil when trip_id is provided")
	}
	if td.GetTripId() != "route_5_0830" {
		t.Errorf("trip.trip_id = %q, want %q", td.GetTripId(), "route_5_0830")
	}
	if td.GetRouteId() != "5" {
		t.Errorf("trip.route_id = %q, want %q", td.GetRouteId(), "5")
	}

	// Check timestamp
	if got := vp.GetTimestamp(); got != 1752566400 {
		t.Errorf("timestamp = %d, want 1752566400", got)
	}
}

// TestBuildFeed_NoTripInfo verifies that a vehicle without trip/route info
// does not include a TripDescriptor (and doesn't crash).
func TestBuildFeed_NoTripInfo(t *testing.T) {
	locations := []model.Location{
		{
			VehicleID: "van-1",
			Latitude:  17.385,
			Longitude: 78.4867,
			Timestamp: 1707350000,
		},
	}

	feed := gtfsrt.BuildFeed(locations)
	vp := feed.Entity[0].GetVehicle()
	if vp.GetTrip() != nil {
		t.Error("TripDescriptor should be nil when no trip/route info is set")
	}
}

// TestBuildFeed_MultipleVehicles verifies that the feed correctly
// contains one entity per vehicle.
func TestBuildFeed_MultipleVehicles(t *testing.T) {
	locations := []model.Location{
		{VehicleID: "bus-1", Latitude: 17.3, Longitude: 78.4, Timestamp: 1000},
		{VehicleID: "bus-2", Latitude: 17.4, Longitude: 78.5, Timestamp: 1001},
		{VehicleID: "bus-3", Latitude: 17.5, Longitude: 78.6, Timestamp: 1002},
	}

	feed := gtfsrt.BuildFeed(locations)
	if len(feed.Entity) != 3 {
		t.Errorf("expected 3 entities, got %d", len(feed.Entity))
	}

	// Verify all vehicle IDs are present
	ids := make(map[string]bool)
	for _, e := range feed.Entity {
		ids[e.GetVehicle().GetVehicle().GetId()] = true
	}
	for _, expected := range []string{"bus-1", "bus-2", "bus-3"} {
		if !ids[expected] {
			t.Errorf("missing vehicle %q in feed", expected)
		}
	}
}

// TestMarshal_RoundTrip verifies that a FeedMessage can be serialized
// to protobuf binary and deserialized back without data loss.
func TestMarshal_RoundTrip(t *testing.T) {
	locations := []model.Location{
		{
			VehicleID: "bus-42",
			TripID:    "route_5_0830",
			RouteID:   "5",
			Latitude:  -1.2921,
			Longitude: 36.8219,
			Bearing:   180.0,
			Speed:     8.5,
			Timestamp: 1752566400,
		},
	}

	feed := gtfsrt.BuildFeed(locations)

	// Marshal to binary
	data, err := gtfsrt.Marshal(feed)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("marshaled data is empty")
	}

	// Unmarshal back
	decoded := &pb.FeedMessage{}
	if err := proto.Unmarshal(data, decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify round-trip integrity
	if decoded.Header.GetGtfsRealtimeVersion() != "2.0" {
		t.Error("version lost in round-trip")
	}
	if len(decoded.Entity) != 1 {
		t.Fatalf("expected 1 entity after round-trip, got %d", len(decoded.Entity))
	}
	pos := decoded.Entity[0].GetVehicle().GetPosition()
	if pos.GetLatitude() < -1.293 || pos.GetLatitude() > -1.291 {
		t.Errorf("latitude lost precision in round-trip: %f", pos.GetLatitude())
	}
}

// TestFeedHeader_Timestamp verifies that the feed header timestamp is
// reasonably close to the current time (within 5 seconds).
func TestFeedHeader_Timestamp(t *testing.T) {
	before := uint64(time.Now().Unix())
	feed := gtfsrt.BuildFeed([]model.Location{})
	after := uint64(time.Now().Unix())

	ts := feed.Header.GetTimestamp()
	if ts < before || ts > after+1 {
		t.Errorf("header timestamp %d not in expected range [%d, %d]", ts, before, after)
	}
}
