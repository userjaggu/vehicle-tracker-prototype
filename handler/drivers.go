package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jaggu/vehicle-tracker-prototype/model"
	"github.com/jaggu/vehicle-tracker-prototype/store"
)

// driversResponse is the JSON shape returned by GET /api/v1/drivers.
type driversResponse struct {
	Drivers []model.Driver `json:"drivers"`
}

// GetDrivers handles GET /api/v1/drivers.
//
// It returns all registered driver profiles.
func GetDrivers(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		drivers := s.GetAllDrivers()
		if drivers == nil {
			drivers = []model.Driver{}
		}

		writeJSON(w, http.StatusOK, driversResponse{Drivers: drivers})
	}
}

// PostDriver handles POST /api/v1/drivers.
//
// It expects a JSON body with at minimum an id and name.
// On success it stores the driver profile and responds with the created profile.
func PostDriver(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var d model.Driver
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if d.ID == "" {
			writeError(w, http.StatusBadRequest, "id is required")
			return
		}
		if d.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}

		created := s.UpsertDriver(d)
		status := http.StatusOK
		if created {
			status = http.StatusCreated
		}
		writeJSON(w, status, d)
	}
}

// GetDriver handles GET /api/v1/drivers/{id}.
//
// It returns the profile for a single driver.
func GetDriver(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "driver id is required")
			return
		}

		d, ok := s.GetDriver(id)
		if !ok {
			writeError(w, http.StatusNotFound, "driver not found")
			return
		}

		writeJSON(w, http.StatusOK, d)
	}
}

// PutDriver handles PUT /api/v1/drivers/{id}.
//
// It replaces the entire driver profile for the given ID.
// The id in the URL path takes precedence over any id in the body.
func PutDriver(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "driver id is required")
			return
		}

		var d model.Driver
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if d.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}

		// Path ID takes precedence.
		d.ID = id

		s.UpsertDriver(d)
		writeJSON(w, http.StatusOK, d)
	}
}

// DeleteDriver handles DELETE /api/v1/drivers/{id}.
//
// It removes the driver profile for the given ID.
func DeleteDriver(s *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "driver id is required")
			return
		}

		if !s.DeleteDriver(id) {
			writeError(w, http.StatusNotFound, "driver not found")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
