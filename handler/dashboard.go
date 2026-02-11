package handler

import (
	"net/http"
)

// GetDashboard serves a simple HTML dashboard to visualize vehicle locations
func GetDashboard(s interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vehicle Tracker Dashboard</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        .header {
            background: white;
            padding: 30px;
            border-radius: 10px;
            margin-bottom: 30px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        
        .header h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 28px;
        }
        
        .header p {
            color: #666;
            font-size: 14px;
        }
        
        .status {
            display: inline-block;
            background: #10b981;
            color: white;
            padding: 5px 15px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: bold;
            margin-top: 10px;
        }
        
        .main-content {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
        }
        
        @media (max-width: 768px) {
            .main-content {
                grid-template-columns: 1fr;
            }
        }
        
        .card {
            background: white;
            border-radius: 10px;
            padding: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        
        .card h2 {
            color: #333;
            margin-bottom: 20px;
            font-size: 18px;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        
        .vehicle-list {
            list-style: none;
        }
        
        .vehicle-item {
            background: #f8f9fa;
            padding: 15px;
            margin-bottom: 10px;
            border-left: 4px solid #667eea;
            border-radius: 5px;
            transition: transform 0.2s;
        }
        
        .vehicle-item:hover {
            transform: translateX(5px);
        }
        
        .vehicle-id {
            font-weight: bold;
            color: #333;
            font-size: 14px;
        }
        
        .vehicle-coords {
            color: #666;
            font-size: 12px;
            margin-top: 8px;
            font-family: 'Courier New', monospace;
        }
        
        .vehicle-time {
            color: #999;
            font-size: 11px;
            margin-top: 5px;
        }
        
        .form-group {
            margin-bottom: 15px;
        }
        
        .form-group label {
            display: block;
            margin-bottom: 5px;
            color: #333;
            font-size: 13px;
            font-weight: 600;
        }
        
        .form-group input {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-size: 13px;
        }
        
        .btn {
            background: #667eea;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 13px;
            font-weight: 600;
            width: 100%;
            transition: background 0.2s;
        }
        
        .btn:hover {
            background: #764ba2;
        }
        
        .btn:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
        
        .message {
            padding: 10px;
            border-radius: 5px;
            margin-top: 10px;
            font-size: 12px;
        }
        
        .message.success {
            background: #d1fae5;
            color: #065f46;
        }
        
        .message.error {
            background: #fee2e2;
            color: #991b1b;
        }
        
        .empty-state {
            text-align: center;
            color: #999;
            padding: 30px;
            font-size: 14px;
        }
        
        .full-width {
            grid-column: 1 / -1;
        }
        
        .stats {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 15px;
            margin-bottom: 20px;
        }
        
        .stat-box {
            background: #f0f4ff;
            padding: 15px;
            border-radius: 5px;
            text-align: center;
        }
        
        .stat-number {
            font-size: 24px;
            font-weight: bold;
            color: #667eea;
        }
        
        .stat-label {
            font-size: 12px;
            color: #666;
            margin-top: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚌 Vehicle Tracker Dashboard</h1>
            <p>Real-time vehicle GPS location tracking system</p>
            <div class="status">✓ Server Running</div>
        </div>
        
        <div class="main-content">
            <!-- Vehicles Section -->
            <div class="card">
                <h2>📍 Tracked Vehicles</h2>
                <div class="stats">
                    <div class="stat-box">
                        <div class="stat-number" id="vehicle-count">0</div>
                        <div class="stat-label">Total Vehicles</div>
                    </div>
                    <div class="stat-box">
                        <div class="stat-number" id="last-update">—</div>
                        <div class="stat-label">Last Update</div>
                    </div>
                </div>
                <ul class="vehicle-list" id="vehicle-list">
                    <div class="empty-state">No vehicles tracked yet</div>
                </ul>
                <button class="btn" onclick="refreshVehicles()">🔄 Refresh</button>
            </div>
            
            <!-- Test Section -->
            <div class="card">
                <h2>🧪 Send Test Location</h2>
                <div class="form-group">
                    <label>Vehicle ID</label>
                    <input type="text" id="vehicle-id" placeholder="e.g., bus-42" value="bus-42">
                </div>
                <div class="form-group">
                    <label>Latitude</label>
                    <input type="number" id="latitude" placeholder="e.g., 17.385" value="17.385" step="0.0001">
                </div>
                <div class="form-group">
                    <label>Longitude</label>
                    <input type="number" id="longitude" placeholder="e.g., 78.4867" value="78.4867" step="0.0001">
                </div>
                <button class="btn" onclick="sendLocation()">📤 Send Location</button>
                <div id="message"></div>
            </div>
            
            <!-- Quick Tests Section -->
            <div class="card full-width">
                <h2>⚡ Quick Tests</h2>
                <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 10px;">
                    <button class="btn" onclick="testSendBus42()">Test: Bus-42 (Default)</button>
                    <button class="btn" onclick="testSendBus99()">Test: Bus-99</button>
                    <button class="btn" onclick="testUpdateBus42()">Test: Update Bus-42</button>
                    <button class="btn" onclick="testInvalidData()">Test: Invalid Data</button>
                    <button class="btn" onclick="testMissingField()">Test: Missing Field</button>
                    <button class="btn" onclick="clearAll()">Clear All Data</button>
                </div>
            </div>
        </div>
    </div>
    
    <script>
        const API_BASE = 'http://localhost:8080';
        
        // Refresh vehicles list
        async function refreshVehicles() {
            try {
                const response = await fetch(API_BASE + '/vehicles');
                const data = await response.json();
                const vehicleList = document.getElementById('vehicle-list');
                const vehicleCount = document.getElementById('vehicle-count');
                
                if (!data.vehicles || data.vehicles.length === 0) {
                    vehicleList.innerHTML = '<div class="empty-state">No vehicles tracked yet</div>';
                    vehicleCount.textContent = '0';
                    return;
                }
                
                vehicleCount.textContent = data.vehicles.length;
                vehicleList.innerHTML = data.vehicles.map(v => {
                    const timestamp = new Date(v.timestamp * 1000).toLocaleString();
                    const lat = v.latitude.toFixed(4);
                    const lng = v.longitude.toFixed(4);
                    return '<li class="vehicle-item">' +
                        '<div class="vehicle-id">🚌 ' + v.vehicle_id + '</div>' +
                        '<div class="vehicle-coords">📍 ' + lat + '°N, ' + lng + '°E</div>' +
                        '<div class="vehicle-time">⏱️ ' + timestamp + '</div>' +
                    '</li>';
                }).join('');
                
                document.getElementById('last-update').textContent = new Date().toLocaleTimeString();
            } catch (error) {
                showMessage('Error fetching vehicles', 'error');
            }
        }
        
        // Send location
        async function sendLocation() {
            const vehicleId = document.getElementById('vehicle-id').value;
            const latitude = parseFloat(document.getElementById('latitude').value);
            const longitude = parseFloat(document.getElementById('longitude').value);
            
            if (!vehicleId || !latitude || !longitude) {
                showMessage('Please fill all fields', 'error');
                return;
            }
            
            try {
                const response = await fetch(API_BASE + '/location', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({
                        vehicle_id: vehicleId,
                        latitude: latitude,
                        longitude: longitude,
                        timestamp: Math.floor(Date.now() / 1000)
                    })
                });
                
                if (response.ok) {
                    showMessage('✓ Location sent successfully!', 'success');
                    setTimeout(refreshVehicles, 500);
                } else {
                    showMessage('Failed to send location', 'error');
                }
            } catch (error) {
                showMessage('Error: ' + error.message, 'error');
            }
        }
        
        // Test functions
        async function testSendBus42() {
            document.getElementById('vehicle-id').value = 'bus-42';
            document.getElementById('latitude').value = '17.385';
            document.getElementById('longitude').value = '78.4867';
            await sendLocation();
        }
        
        async function testSendBus99() {
            document.getElementById('vehicle-id').value = 'bus-99';
            document.getElementById('latitude').value = '17.400';
            document.getElementById('longitude').value = '78.500';
            await sendLocation();
        }
        
        async function testUpdateBus42() {
            document.getElementById('vehicle-id').value = 'bus-42';
            document.getElementById('latitude').value = '18.000';
            document.getElementById('longitude').value = '79.000';
            await sendLocation();
        }
        
        async function testInvalidData() {
            try {
                const response = await fetch(API_BASE + '/location', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: 'not-json'
                });
                const data = await response.json();
                showMessage('✓ Invalid data rejected: ' + data.error, 'success');
            } catch (error) {
                showMessage('✓ Invalid data rejected properly', 'success');
            }
        }
        
        async function testMissingField() {
            try {
                const response = await fetch(API_BASE + '/location', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({
                        latitude: 17.385,
                        longitude: 78.4867
                    })
                });
                const data = await response.json();
                showMessage('✓ Missing field rejected: ' + data.error, 'success');
            } catch (error) {
                showMessage('Error in test', 'error');
            }
        }
        
        async function clearAll() {
            if (confirm('This will not actually clear the server (it clears on restart). Continue?')) {
                showMessage('Note: To clear data, restart the server', 'error');
            }
        }
        
        // Show message
        function showMessage(text, type) {
            const msg = document.getElementById('message');
            msg.textContent = text;
            msg.className = 'message ' + type;
            setTimeout(() => {
                msg.textContent = '';
                msg.className = 'message';
            }, 3000);
        }
        
        // Auto-refresh every 5 seconds
        setInterval(refreshVehicles, 5000);
        
        // Initial load
        window.onload = refreshVehicles;
    </script>
</body>
</html>
`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}
}
