// Package handler implements the HTTP handlers for the vehicle tracker API.
//
// Handlers:
//   - PostLocation:  accepts a vehicle's GPS update       (POST /location)
//   - GetVehicles:   returns all known vehicle locations   (GET  /vehicles)
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jaggu/vehicle-tracker-prototype/model"
	"github.com/jaggu/vehicle-tracker-prototype/store"
)

// PostLocation handles POST /location.
//
// It expects a JSON body with vehicle_id, latitude, longitude, and timestamp.
// On success it stores the location and responds with {"status": "ok"}.
func PostLocation(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ── Only accept POST ────────────────────────────────────────
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "Only POST is allowed")
			return
		}

		// ── Decode request body ─────────────────────────────────────
		var loc model.Location
		if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		// ── Validate required fields ────────────────────────────────
		if loc.VehicleID == "" {
			writeError(w, http.StatusBadRequest, "vehicle_id is required")
			return
		}
		if loc.Latitude == 0 && loc.Longitude == 0 {
			writeError(w, http.StatusBadRequest, "latitude and longitude are required")
			return
		}

		// ── Store the location ──────────────────────────────────────
		s.UpdateLocation(loc)

		// ── Respond ─────────────────────────────────────────────────
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
