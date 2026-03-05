package handler

import (
	"net/http"
	"time"

	"github.com/jaggu/vehicle-tracker-prototype/model"
	"github.com/jaggu/vehicle-tracker-prototype/store"
)

// statusResponse is the JSON shape returned by GET /api/v1/status.
type statusResponse struct {
	Status              string `json:"status"`
	ActiveVehicles      int    `json:"active_vehicles"`
	TotalVehicles       int    `json:"total_vehicles"`
	StalenessThreshold  string `json:"staleness_threshold"`
	ServerTimeUTC       string `json:"server_time_utc"`
	FeedEndpoint        string `json:"feed_endpoint"`
	FeedEndpointJSON    string `json:"feed_endpoint_json"`
}

// GetStatus handles GET /api/v1/status.
//
// It returns basic system health information: how many vehicles are
// actively reporting, the staleness threshold in use, and the feed URL.
func GetStatus(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "Only GET is allowed")
			return
		}

		resp := statusResponse{
			Status:             "ok",
			ActiveVehicles:     s.ActiveVehicleCount(model.DefaultStalenessThreshold),
			TotalVehicles:      s.TotalVehicleCount(),
			StalenessThreshold: model.DefaultStalenessThreshold.String(),
			ServerTimeUTC:      time.Now().UTC().Format(time.RFC3339),
			FeedEndpoint:       "/gtfs-rt/vehicle-positions",
			FeedEndpointJSON:   "/gtfs-rt/vehicle-positions?format=json",
		}

		writeJSON(w, http.StatusOK, resp)
	}
}
