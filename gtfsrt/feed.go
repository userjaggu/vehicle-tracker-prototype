// Package gtfsrt converts in-memory vehicle locations into a valid
// GTFS-Realtime Vehicle Positions FeedMessage.
//
// The feed follows the official specification:
// https://gtfs.org/documentation/realtime/proto/
//
// Key decisions:
//   - incrementality is FULL_DATASET (every response contains all active vehicles)
//   - only VehiclePosition entities are produced (no TripUpdates or Alerts)
//   - vehicles that exceed the staleness threshold are excluded
package gtfsrt

import (
	"fmt"
	"time"

	"github.com/jaggu/vehicle-tracker-prototype/model"
	pb "github.com/jaggu/vehicle-tracker-prototype/proto/gtfsrt"
	"google.golang.org/protobuf/proto"
)

// gtfsRTVersion is the protocol version written into every FeedHeader.
const gtfsRTVersion = "2.0"

// BuildFeed creates a complete GTFS-RT FeedMessage from a slice of active
// vehicle locations.
//
// Each location becomes a FeedEntity containing a VehiclePosition with
// position (lat, lon, bearing, speed), trip descriptor, vehicle descriptor,
// and timestamp.
func BuildFeed(locations []model.Location) *pb.FeedMessage {
	now := uint64(time.Now().Unix())
	version := gtfsRTVersion
	incrementality := pb.FeedHeader_FULL_DATASET

	header := &pb.FeedHeader{
		GtfsRealtimeVersion: &version,
		Incrementality:      &incrementality,
		Timestamp:           &now,
	}

	entities := make([]*pb.FeedEntity, 0, len(locations))

	for i := range locations {
		loc := &locations[i]
		entity := buildEntity(loc)
		entities = append(entities, entity)
	}

	return &pb.FeedMessage{
		Header: header,
		Entity: entities,
	}
}

// buildEntity converts a single Location into a FeedEntity wrapping a
// VehiclePosition message.
func buildEntity(loc *model.Location) *pb.FeedEntity {
	id := fmt.Sprintf("vehicle-%s", loc.VehicleID)

	lat := float32(loc.Latitude)
	lon := float32(loc.Longitude)

	position := &pb.Position{
		Latitude:  &lat,
		Longitude: &lon,
	}

	// Only include bearing and speed when they carry meaningful data.
	if loc.Bearing != 0 {
		b := loc.Bearing
		position.Bearing = &b
	}
	if loc.Speed != 0 {
		s := loc.Speed
		position.Speed = &s
	}

	vp := &pb.VehiclePosition{
		Position: position,
		Vehicle: &pb.VehicleDescriptor{
			Id:    &loc.VehicleID,
			Label: &loc.VehicleID,
		},
	}

	// Attach timestamp if the client provided one.
	if loc.Timestamp > 0 {
		ts := uint64(loc.Timestamp)
		vp.Timestamp = &ts
	}

	// Attach trip information when available.
	if loc.TripID != "" || loc.RouteID != "" {
		td := &pb.TripDescriptor{}
		if loc.TripID != "" {
			td.TripId = &loc.TripID
		}
		if loc.RouteID != "" {
			td.RouteId = &loc.RouteID
		}
		vp.Trip = td
	}

	return &pb.FeedEntity{
		Id:      &id,
		Vehicle: vp,
	}
}

// Marshal serializes a FeedMessage to the protobuf wire format.
func Marshal(feed *pb.FeedMessage) ([]byte, error) {
	return proto.Marshal(feed)
}
