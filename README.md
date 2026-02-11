# Vehicle Tracker Backend Prototype

## What I Built

This is a simple backend server that receives GPS location updates from vehicles and stores them. The idea is:
- Vehicles (or an Android app) send their GPS coordinates to the server
- The server stores the latest position for each vehicle
- Anyone can ask the server "where are all the vehicles?" and get back all the locations

Think of it like a restaurant keeping a notepad: drivers call in to say where they are, and the restaurant writes it down. When a customer asks "where is my delivery?", the restaurant reads the notepad and tells them.

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
- That's it! No npm, no pip, no databases to install

### Starting the Server

```bash
go run main.go
```

You should see:
```
Vehicle Tracker server listening on http://localhost:8080
```

The server is now running and ready to receive requests.

---

## How to Use It

There are two endpoints:

### 1. Send a Vehicle Location

**Send GPS coordinates from a vehicle:**

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

**What each field means:**
- `vehicle_id`: A unique name for the vehicle (e.g., "bus-42", "truck-101")
- `latitude`: The north-south position (decimal degrees)
- `longitude`: The east-west position (decimal degrees)
- `timestamp`: When this location was recorded (Unix timestamp, optional)

---

### 2. Get All Vehicle Locations

**Ask the server for all known vehicle positions:**

```bash
curl http://localhost:8080/vehicles
```

**The server responds with something like:**
```json
{
  "vehicles": [
    {
      "vehicle_id": "bus-42",
      "latitude": 17.385,
      "longitude": 78.4867,
      "timestamp": 1707350000
    },
    {
      "vehicle_id": "bus-99",
      "latitude": 17.400,
      "longitude": 78.500,
      "timestamp": 1707350100
    }
  ]
}
```

If no vehicles have sent data yet, the response is:
```json
{ "vehicles": [] }
```

---

## Testing It Yourself

Here's a complete example you can run to test everything:

```bash
# 1. Start the server in the background
go run main.go &

# 2. Send bus-42's location
curl -s -X POST http://localhost:8080/location \
  -H "Content-Type: application/json" \
  -d '{"vehicle_id":"bus-42","latitude":17.385,"longitude":78.4867,"timestamp":1707350000}'

# 3. Send bus-99's location
curl -s -X POST http://localhost:8080/location \
  -H "Content-Type: application/json" \
  -d '{"vehicle_id":"bus-99","latitude":17.400,"longitude":78.500,"timestamp":1707350100}'

# 4. Ask for all vehicles
curl -s http://localhost:8080/vehicles | python3 -m json.tool

# 5. Update bus-42's location (send new GPS)
curl -s -X POST http://localhost:8080/location \
  -H "Content-Type: application/json" \
  -d '{"vehicle_id":"bus-42","latitude":18.000,"longitude":79.000,"timestamp":1707360000}'

# 6. Check that bus-42 updated and bus-99 is still there
curl -s http://localhost:8080/vehicles | python3 -m json.tool
```

---

## What I Tested

I tested this thoroughly with `curl` (a command-line tool that simulates what an app would do):

### Test Results (All Passed ✅)

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

When a vehicle sends an update for `bus-42`, the row for `bus-42` is overwritten with the new data. The old data is gone (since it's in-memory only).

---

## Why I Made These Design Choices

### Go + Standard Library (No External Dependencies)
- Go has excellent built-in support for HTTP and JSON
- No need to install packages with npm or pip
- Compiles to a single binary file that just works
- Easy to deploy anywhere

### In-Memory Storage
- Simplest way to store data
- No database setup needed
- Data resets when server restarts (fine for a prototype)
- Will switch to a real database later if needed

### Thread-Safe Access with Mutex
- Multiple requests can happen at the same time
- A "lock" ensures they don't corrupt each other's data
- This is important for production reliability

### Only Two Endpoints
- Keeps it simple and focused
- Does one thing well: accept and return GPS data
- Easy to test and reason about

---

## What Happens Next

After this backend works, the plan is:

1. **Build Android App**: Create an app that reads the phone's GPS and sends data to this server every 10-15 seconds
2. **Add Persistence**: Switch from in-memory storage to a real database
3. **Add Authentication**: Only allow trusted apps to send data
4. **Add More Features**: Like tracking history, speed, direction, etc.

But for now, this backend does what it needs to do: accept GPS updates and return them on demand.

---



---

## Running on Your Machine

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

That's it! No installation, no configuration needed.
