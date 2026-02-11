# Vehicle Tracker Backend Prototype

## What I Built

This is a simple backend server that receives GPS location updates from vehicles and stores them. The idea is:
- Vehicles (or an Android app) send their GPS coordinates to the server
- The server stores the latest position for each vehicle
- Anyone can ask the server "where are all the vehicles?" and get back all the locations



---

## Why I'm Building This

The OneBusAway project wants real-time vehicle tracking. Before building the full Android app and a complex system, I'm starting with this minimal backend to prove the basic concept works:
1. Can we receive GPS data?
2. Can we store it properly?
3. Can we retrieve it when asked?

This prototype answers "yes" to all three.

---

## What I'm NOT Doing (Not Yet)

Right now, I'm **only** focusing on the core flow. These are future improvements:

- **Authentication**: Anyone can send locations (fine for prototype)
- **Database**: Data only stays in memory (resets when server restarts)
- **Offline mode**: If the app is offline, no data is sent (simple for now)
- **GTFS-Realtime**: Not converting data to that standard yet
- **Logging**: Not tracking server activity in files
- **Android app**: Will build this next (this is the backend only)

---

## How the Code is Organized

```
vehicle-tracker-prototype/
├── main.go              # The entry point - runs: go run main.go
├── server/
│   └── server.go        # Sets up the routes and starts listening on port 8080
├── handler/
│   ├── location.go      # Handles: POST /location (receives GPS updates)
│   ├── vehicles.go      # Handles: GET /vehicles (returns all stored locations)
│   ├── dashboard.go     # Serves the web dashboard at http://localhost:8080
│   └── helpers.go       # Shared code for sending JSON responses
├── model/
│   └── vehicle.go       # Defines what a "Location" looks like (GPS point)
├── store/
│   └── memory.go        # In-memory map to store all vehicle locations
├── go.mod               # Go module definition (zero external dependencies!)
└── README.md            # This file
```

---

## Getting Started

### What You Need

- Go 1.21 or higher ([download here](https://go.dev/dl/))

### Running on Your Machine

```bash
# Clone/download the code
cd vehicle-tracker-prototype

# Build it
go build -o vehicle-tracker

# Run it
./vehicle-tracker

# In another terminal, test it
curl http://localhost:8080/vehicles
```

Then open your browser and visit:
```
http://localhost:8080
```

You'll see an interactive dashboard to test the API and view all tracked vehicles.

---

## How to Use It

There are two ways to interact with the system:

### Option 1: Web Dashboard

Visit `http://localhost:8080` in your browser. You can:
- View all tracked vehicles in real-time
- Send test locations using the form
- Click quick test buttons to test different scenarios
- See error handling in action

### Option 2: Command Line (curl)

**Send a vehicle location:**

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

**The server responds with:**
```json
{ "status": "ok" }
```

**Get all vehicle locations:**

```bash
curl http://localhost:8080/vehicles
```

**The server responds with:**
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

---

## What I Tested

I tested this thoroughly with `curl` (a command-line tool that simulates what an app would do):

### Test Results (All Passed )

1. **Empty server**: Getting locations when no data sent → Returns empty array
2. **Sending one vehicle**: POST bus-42 location → Stored correctly
3. **Sending multiple vehicles**: POST bus-99 location → Both vehicles stored
4. **Overwriting location**: POST new position for bus-42 → Old position replaced, bus-99 unchanged
5. **Missing vehicle_id**: POST without vehicle_id → Server rejects with error
6. **Missing coordinates**: POST without lat/lng → Server rejects with error
7. **Invalid JSON**: POST malformed data → Server rejects with error
8. **Wrong HTTP method**: GET to /location endpoint → Server rejects with error

The server properly validates input and rejects bad data while accepting valid data.

---

## How It Works Behind the Scenes

### When a Vehicle Sends GPS Data

```
1. Vehicle (or curl) sends: POST /location with GPS coordinates
2. Server receives the request
3. Server checks: "Is this valid JSON?" → If no, reject
4. Server checks: "Does it have vehicle_id?" → If no, reject
5. Server checks: "Does it have lat/lng?" → If no, reject
6. Server stores: vehicle_id → Location{lat, lng, timestamp}
7. Server responds: {"status": "ok"}
```

### When Someone Asks for All Locations

```
1. Client sends: GET /vehicles
2. Server receives the request
3. Server gets all stored locations from memory
4. Server packages them into JSON format
5. Server responds with: {"vehicles": [...all locations...]}
```

### Storage

The server keeps one "map" in memory. Think of it like a table:

```
vehicle_id  | latitude | longitude | timestamp
------------|----------|-----------|----------
bus-42      | 17.385   | 78.4867   | 1707350000
bus-99      | 17.400   | 78.500    | 1707350100
```

When a vehicle sends an update for `bus-42`, the row for `bus-42` is overwritten with the new data.

---




