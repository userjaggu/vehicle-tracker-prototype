package handler

import (
	"net/http"

	"github.com/jaggu/vehicle-tracker-prototype/gtfsrt"
	"github.com/jaggu/vehicle-tracker-prototype/model"
	"github.com/jaggu/vehicle-tracker-prototype/store"
	"google.golang.org/protobuf/encoding/protojson"
)

// GetGTFSRT handles GET /gtfs-rt/vehicle-positions.
//
// By default it returns the binary protobuf representation
// (Content-Type: application/x-protobuf).
//
// Append ?format=json to receive the same feed as human-readable JSON
// (useful for debugging and for consumers that prefer JSON over protobuf).
//
// Only vehicles that have reported within the staleness threshold are
// included.  A feed with zero active vehicles is still valid — it returns
// a FeedMessage with an empty entity list.
func GetGTFSRT(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only accept GET
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "Only GET is allowed")
			return
		}

		// Fetch only active (non-stale) vehicles
		locations := s.GetActiveLocations(model.DefaultStalenessThreshold)

		// Build the GTFS-RT FeedMessage
		feed := gtfsrt.BuildFeed(locations)

		// Check if the caller wants JSON output for debugging
		if r.URL.Query().Get("format") == "json" {
			marshaler := protojson.MarshalOptions{
				EmitUnpopulated: false,
				Indent:          "  ",
			}
			jsonBytes, err := marshaler.Marshal(feed)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Failed to marshal feed to JSON")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonBytes)
			return
		}

		// Default: binary protobuf
		data, err := gtfsrt.Marshal(feed)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to marshal feed")
			return
		}

		w.Header().Set("Content-Type", "application/x-protobuf")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
