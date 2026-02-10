package handler

import (
	"net/http"

	"github.com/jaggu/vehicle-tracker-prototype/model"
	"github.com/jaggu/vehicle-tracker-prototype/store"
)

// vehiclesResponse is the JSON shape returned by GET /vehicles.
type vehiclesResponse struct {
	Vehicles []model.Location `json:"vehicles"`
}

// GetVehicles handles GET /vehicles.
//
// It returns the latest known GPS location for every vehicle
// that has reported at least one update.
func GetVehicles(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ── Only accept GET ─────────────────────────────────────────
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "Only GET is allowed")
			return
		}

		// ── Fetch all locations from the store ──────────────────────
		locations := s.GetAllLocations()

		// ── Always return an array, even if empty ───────────────────
		if locations == nil {
			locations = []model.Location{}
		}

		// ── Respond ─────────────────────────────────────────────────
		writeJSON(w, http.StatusOK, vehiclesResponse{Vehicles: locations})
	}
}
