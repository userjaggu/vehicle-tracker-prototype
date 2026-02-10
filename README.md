# Vehicle Tracker — Minimal Backend Prototype

> A minimal backend that accepts live vehicle GPS updates and exposes the latest
> known positions using in-memory storage only.

---

## What Is This?

A proof-of-concept HTTP server for the **OneBusAway Vehicle Tracker** project.  
Drivers (or test scripts) send GPS coordinates, and the server stores the latest
known location for each vehicle in memory.

---

## What Problem Does It Solve?

Transit agencies need to know **where their vehicles are right now**.  
This prototype demonstrates the simplest possible version of that:

1. A vehicle sends its GPS position to the server.
2. Anyone can query the server to see all known vehicle positions.

---

## What It Intentionally Does NOT Do (Yet)

| Feature               | Status       |
|----------------------|--------------|
| Authentication        | ❌ Not included |
| Database / persistence | ❌ In-memory only |
| Offline buffering     | ❌ Not included |
| GTFS-Realtime output  | ❌ Not included |
| Android app           | ❌ Not included |
| Logging framework     | ❌ Not included |
| Production deployment | ❌ Not included |

These are all **future work** — this repo focuses purely on the core location flow.

---

## Project Structure

```
vehicle-tracker-prototype/
├── main.go              # Entry point — starts the server on :8080
├── server/
│   └── server.go        # Wires routes and starts HTTP listener
├── handler/
│   ├── location.go      # POST /location  — accept a GPS update
│   ├── vehicles.go      # GET  /vehicles  — return all known locations
│   └── helpers.go       # Shared JSON response helpers
├── model/
│   └── vehicle.go       # Location data structure
├── store/
│   └── memory.go        # Thread-safe in-memory storage (map + mutex)
├── go.mod
└── README.md
```

---

## How to Run

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)

### Start the Server

```bash
go run main.go
```

Output:

```
Vehicle Tracker server listening on http://localhost:8080
```

---

## API Reference

### `POST /location`

Send a vehicle's current GPS position.

**Request:**

```bash
curl -X POST http://localhost:8080/location \
  -H "Content-Type: application/json" \
  -d '{
    "vehicle_id": "bus-42",
    "latitude": 17.385,
    "longitude": 78.4867,
    "timestamp": 1707350000
  }'
```

**Response:**

```json
{ "status": "ok" }
```

**Validation rules:**

- `vehicle_id` must be a non-empty string.
- `latitude` and `longitude` must be provided.
- `timestamp` is optional (but recommended).

---

### `GET /vehicles`

Retrieve the latest known location for all vehicles.

**Request:**

```bash
curl http://localhost:8080/vehicles
```

**Response:**

```json
{
  "vehicles": [
    {
      "vehicle_id": "bus-42",
      "latitude": 17.385,
      "longitude": 78.4867,
      "timestamp": 1707350000
    }
  ]
}
```

- Returns an empty array `[]` if no vehicles have reported yet.
- Only the **latest** location per vehicle is stored.

---

## Quick Test (Copy-Paste)

```bash
# 1. Start the server
go run main.go &

# 2. Send a location update
curl -s -X POST http://localhost:8080/location \
  -H "Content-Type: application/json" \
  -d '{"vehicle_id":"bus-42","latitude":17.385,"longitude":78.4867,"timestamp":1707350000}'

# 3. Send another vehicle
curl -s -X POST http://localhost:8080/location \
  -H "Content-Type: application/json" \
  -d '{"vehicle_id":"bus-99","latitude":17.400,"longitude":78.500,"timestamp":1707350100}'

# 4. Fetch all vehicles
curl -s http://localhost:8080/vehicles | python3 -m json.tool
```

---

## Design Decisions

| Decision | Rationale |
|----------|-----------|
| **Go + standard library** | Zero dependencies, single binary, easy to read |
| **In-memory map** | Simplest storage; persistence is a future concern |
| **sync.RWMutex** | Safe concurrent reads/writes without a database |
| **No authentication** | Prototype scope — will be added later |
| **Latest-only storage** | Each vehicle overwrites its previous position |

---

## License

This project is part of the [OneBusAway](https://onebusaway.org/) ecosystem.
