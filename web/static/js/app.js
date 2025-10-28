// Reolink Server Frontend Application

// Configuration
const API_BASE = '/api/v1';
let authToken = localStorage.getItem('authToken');
let eventWebSocket = null;

// API Helper Functions
async function apiRequest(endpoint, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    if (authToken) {
        headers['Authorization'] = `Bearer ${authToken}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers
    });

    if (response.status === 401) {
        // Token expired or invalid
        localStorage.removeItem('authToken');
        authToken = null;
        showLogin();
        throw new Error('Authentication required');
    }

    if (!response.ok) {
        const error = await response.json().catch(() => ({ message: 'Request failed' }));
        throw new Error(error.message || 'Request failed');
    }

    return response.json();
}

// Authentication
function showLogin() {
    const loginHTML = `
        <div class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
            <div class="bg-white p-8 rounded-lg shadow-xl max-w-md w-full">
                <h2 class="text-2xl font-bold mb-4">Login</h2>
                <form id="loginForm">
                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-bold mb-2">Username</label>
                        <input type="text" id="username" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700" required>
                    </div>
                    <div class="mb-6">
                        <label class="block text-gray-700 text-sm font-bold mb-2">Password</label>
                        <input type="password" id="password" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700" required>
                    </div>
                    <div class="flex items-center justify-between">
                        <button type="submit" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                            Sign In
                        </button>
                    </div>
                    <div id="loginError" class="mt-4 text-red-500 text-sm hidden"></div>
                </form>
            </div>
        </div>
    `;

    document.body.insertAdjacentHTML('beforeend', loginHTML);

    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch(`${API_BASE}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            if (!response.ok) {
                throw new Error('Invalid credentials');
            }

            const result = await response.json();
            // API returns { success: true, data: { token: "...", ... } }
            authToken = result.data.token;
            localStorage.setItem('authToken', authToken);

            // Remove login modal
            document.querySelector('.fixed').remove();

            // Initialize app
            init();
        } catch (error) {
            console.error('Login error:', error);
            document.getElementById('loginError').textContent = error.message;
            document.getElementById('loginError').classList.remove('hidden');
        }
    });
}

// System Status
async function loadSystemStatus() {
    try {
        const healthResp = await fetch('/health').then(r => r.json());
        const readyResp = await fetch('/ready').then(r => r.json());

        const health = healthResp.data;
        const ready = readyResp.data;

        document.getElementById('status').innerHTML = `
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <p class="text-sm text-gray-500">Health</p>
                    <p class="text-lg font-semibold ${health.status === 'healthy' ? 'text-green-600' : 'text-red-600'}">
                        ${health.status.toUpperCase()}
                    </p>
                </div>
                <div>
                    <p class="text-sm text-gray-500">Database</p>
                    <p class="text-lg font-semibold ${ready.components.database === 'healthy' ? 'text-green-600' : 'text-red-600'}">
                        ${ready.components.database.toUpperCase()}
                    </p>
                </div>
            </div>
        `;
    } catch (error) {
        console.error('Failed to load system status:', error);
        document.getElementById('status').innerHTML = `
            <p class="text-red-600">Failed to load system status</p>
        `;
    }
}

// Cameras
async function loadCameras() {
    try {
        const result = await apiRequest('/cameras');
        const cameras = result.data.cameras || [];

        if (cameras.length === 0) {
            document.getElementById('cameras').innerHTML = `
                <p class="text-gray-600">No cameras configured yet.</p>
            `;
            return;
        }

        document.getElementById('cameras').innerHTML = cameras.map(camera => `
            <div class="border rounded-lg p-4 hover:shadow-lg transition-shadow">
                <div class="flex justify-between items-start mb-2">
                    <h3 class="font-semibold text-lg">${camera.name}</h3>
                    <span class="px-2 py-1 text-xs rounded ${camera.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}">
                        ${camera.enabled ? 'Active' : 'Inactive'}
                    </span>
                </div>
                <p class="text-sm text-gray-600 mb-2">${camera.host}</p>
                <p class="text-xs text-gray-500 mb-3">${camera.model || 'Unknown Model'}</p>
                <div class="flex gap-2">
                    <button onclick="viewCamera('${camera.id}')" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white text-sm py-1 px-2 rounded">
                        View
                    </button>
                    <button onclick="getSnapshot('${camera.id}')" class="flex-1 bg-green-500 hover:bg-green-600 text-white text-sm py-1 px-2 rounded">
                        Snapshot
                    </button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load cameras:', error);
        document.getElementById('cameras').innerHTML = `
            <p class="text-red-600">Failed to load cameras</p>
        `;
    }
}

// Add Camera
function showAddCameraModal() {
    const modalHTML = `
        <div id="addCameraModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
            <div class="bg-white p-8 rounded-lg shadow-xl max-w-md w-full">
                <h2 class="text-2xl font-bold mb-4">Add Camera</h2>
                <form id="addCameraForm">
                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-bold mb-2">Name</label>
                        <input type="text" id="cameraName" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700" required>
                    </div>
                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-bold mb-2">Host/IP</label>
                        <input type="text" id="cameraHost" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700" placeholder="192.168.1.100" required>
                    </div>
                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-bold mb-2">Port</label>
                        <input type="number" id="cameraPort" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700" value="80">
                    </div>
                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-bold mb-2">Username</label>
                        <input type="text" id="cameraUsername" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700" value="admin" required>
                    </div>
                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-bold mb-2">Password</label>
                        <input type="password" id="cameraPassword" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700" required>
                    </div>
                    <div class="mb-4">
                        <label class="flex items-center">
                            <input type="checkbox" id="cameraUseHTTPS" class="mr-2">
                            <span class="text-sm text-gray-700">Use HTTPS</span>
                        </label>
                    </div>
                    <div class="mb-6">
                        <label class="flex items-center">
                            <input type="checkbox" id="cameraSkipVerify" class="mr-2" checked>
                            <span class="text-sm text-gray-700">Skip SSL Verification</span>
                        </label>
                    </div>
                    <div id="addCameraError" class="mb-4 text-red-500 text-sm hidden"></div>
                    <div class="flex gap-2">
                        <button type="submit" class="flex-1 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                            Add Camera
                        </button>
                        <button type="button" onclick="closeAddCameraModal()" class="flex-1 bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded">
                            Cancel
                        </button>
                    </div>
                </form>
            </div>
        </div>
    `;

    document.body.insertAdjacentHTML('beforeend', modalHTML);

    document.getElementById('addCameraForm').addEventListener('submit', async (e) => {
        e.preventDefault();

        const cameraData = {
            name: document.getElementById('cameraName').value,
            host: document.getElementById('cameraHost').value,
            port: parseInt(document.getElementById('cameraPort').value) || 80,
            username: document.getElementById('cameraUsername').value,
            password: document.getElementById('cameraPassword').value,
            use_https: document.getElementById('cameraUseHTTPS').checked,
            skip_verify: document.getElementById('cameraSkipVerify').checked
        };

        try {
            await apiRequest('/cameras', {
                method: 'POST',
                body: JSON.stringify(cameraData)
            });

            closeAddCameraModal();
            await loadCameras();
        } catch (error) {
            console.error('Failed to add camera:', error);
            document.getElementById('addCameraError').textContent = error.message;
            document.getElementById('addCameraError').classList.remove('hidden');
        }
    });
}

function closeAddCameraModal() {
    const modal = document.getElementById('addCameraModal');
    if (modal) {
        modal.remove();
    }
}

// View Camera Details
async function viewCamera(cameraId) {
    try {
        const cameraResult = await apiRequest(`/cameras/${cameraId}`);
        const statusResult = await apiRequest(`/cameras/${cameraId}/status`);
        const camera = cameraResult.data;
        const status = statusResult.data;

        const modal = `
            <div id="cameraModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
                <div class="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-screen overflow-y-auto">
                    <div class="p-6">
                        <div class="flex justify-between items-center mb-4">
                            <h2 class="text-2xl font-bold">${camera.name}</h2>
                            <button onclick="closeModal()" class="text-gray-500 hover:text-gray-700">
                                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                                </svg>
                            </button>
                        </div>

                        <div class="grid grid-cols-2 gap-4 mb-4">
                            <div>
                                <p class="text-sm text-gray-500">Host</p>
                                <p class="font-semibold">${camera.host}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-500">Model</p>
                                <p class="font-semibold">${camera.model || 'Unknown'}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-500">Firmware</p>
                                <p class="font-semibold">${camera.firmware_ver || 'Unknown'}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-500">Status</p>
                                <p class="font-semibold ${status.online ? 'text-green-600' : 'text-red-600'}">
                                    ${status.online ? 'Online' : 'Offline'}
                                </p>
                            </div>
                        </div>

                        <div class="mb-4">
                            <h3 class="font-semibold mb-2">Live Stream</h3>
                            <div id="streamContainer" class="bg-gray-200 rounded aspect-video flex items-center justify-center">
                                <button onclick="startStream('${cameraId}')" class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded">
                                    Start Stream
                                </button>
                            </div>
                        </div>

                        <div class="grid grid-cols-2 gap-2">
                            <button onclick="controlPTZ('${cameraId}', 'up')" class="bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded">
                                PTZ Up
                            </button>
                            <button onclick="controlPTZ('${cameraId}', 'down')" class="bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded">
                                PTZ Down
                            </button>
                            <button onclick="controlPTZ('${cameraId}', 'left')" class="bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded">
                                PTZ Left
                            </button>
                            <button onclick="controlPTZ('${cameraId}', 'right')" class="bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded">
                                PTZ Right
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;

        document.body.insertAdjacentHTML('beforeend', modal);
    } catch (error) {
        alert('Failed to load camera details: ' + error.message);
    }
}

function closeModal() {
    const modal = document.getElementById('cameraModal');
    if (modal) modal.remove();
}

// Start Stream
function startStream(cameraId) {
    const streamContainer = document.getElementById('streamContainer');

    // Create video element for HLS stream
    const videoHTML = `
        <video id="liveStream" class="w-full h-full rounded" controls autoplay muted>
            <source src="${API_BASE}/cameras/${cameraId}/stream/hls/playlist.m3u8" type="application/x-mpegURL">
            Your browser does not support HLS streaming.
        </video>
        <p class="text-xs text-gray-500 mt-2">
            Stream URL: ${API_BASE}/cameras/${cameraId}/stream/hls/playlist.m3u8
        </p>
    `;

    streamContainer.innerHTML = videoHTML;

    // Try to load HLS using native support or hls.js
    const video = document.getElementById('liveStream');
    const streamUrl = `${API_BASE}/cameras/${cameraId}/stream/hls/playlist.m3u8`;

    if (video.canPlayType('application/vnd.apple.mpegurl')) {
        // Native HLS support (Safari)
        video.src = streamUrl;
    } else if (typeof Hls !== 'undefined') {
        // Use hls.js for other browsers
        const hls = new Hls();
        hls.loadSource(streamUrl);
        hls.attachMedia(video);
    } else {
        streamContainer.innerHTML = `
            <div class="text-center p-4">
                <p class="text-red-600 mb-2">HLS streaming not supported in this browser</p>
                <p class="text-sm text-gray-600">Try using Safari or install hls.js</p>
                <a href="${API_BASE}/cameras/${cameraId}/stream/flv" target="_blank"
                   class="inline-block mt-2 bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded">
                    Open FLV Stream
                </a>
            </div>
        `;
    }
}

// Get Snapshot
async function getSnapshot(cameraId) {
    try {
        const response = await fetch(`${API_BASE}/cameras/${cameraId}/snapshot`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });

        if (!response.ok) throw new Error('Failed to get snapshot');

        const blob = await response.blob();
        const url = URL.createObjectURL(blob);

        const modal = `
            <div id="snapshotModal" class="fixed inset-0 bg-gray-600 bg-opacity-75 flex items-center justify-center z-50" onclick="this.remove()">
                <div class="max-w-4xl max-h-screen p-4">
                    <img src="${url}" class="max-w-full max-h-full rounded-lg shadow-xl" />
                </div>
            </div>
        `;

        document.body.insertAdjacentHTML('beforeend', modal);
    } catch (error) {
        alert('Failed to get snapshot: ' + error.message);
    }
}

// PTZ Control
async function controlPTZ(cameraId, direction) {
    try {
        await apiRequest(`/cameras/${cameraId}/ptz/move`, {
            method: 'POST',
            body: JSON.stringify({ operation: direction, speed: 32 })
        });
    } catch (error) {
        alert('PTZ control failed: ' + error.message);
    }
}

// Events
async function loadEvents() {
    try {
        const result = await apiRequest('/events?limit=10');
        const events = result.data.events || [];

        if (events.length === 0) {
            document.getElementById('events').innerHTML = `
                <p class="text-gray-600">No events yet.</p>
            `;
            return;
        }

        document.getElementById('events').innerHTML = events.map(event => `
            <div class="border-l-4 ${getEventColor(event.type)} pl-4 py-2">
                <div class="flex justify-between items-start">
                    <div>
                        <p class="font-semibold">${event.type.replace(/_/g, ' ').toUpperCase()}</p>
                        <p class="text-sm text-gray-600">${new Date(event.timestamp).toLocaleString()}</p>
                    </div>
                    ${!event.acknowledged ? `
                        <button onclick="acknowledgeEvent('${event.id}')" class="text-sm text-blue-600 hover:text-blue-800">
                            Acknowledge
                        </button>
                    ` : ''}
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load events:', error);
    }
}

function getEventColor(type) {
    const colors = {
        'motion_detected': 'border-yellow-500',
        'ai_person': 'border-orange-500',
        'ai_vehicle': 'border-blue-500',
        'ai_pet': 'border-green-500'
    };
    return colors[type] || 'border-gray-500';
}

async function acknowledgeEvent(eventId) {
    try {
        await apiRequest(`/events/${eventId}/acknowledge`, { method: 'PUT' });
        loadEvents();
    } catch (error) {
        alert('Failed to acknowledge event: ' + error.message);
    }
}

// WebSocket for Real-time Events
function connectEventWebSocket() {
    if (!authToken) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}${API_BASE}/ws/events?token=${authToken}`;

    eventWebSocket = new WebSocket(wsUrl);

    eventWebSocket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log('New event:', data);
        loadEvents(); // Reload events list

        // Show notification
        showNotification(data);
    };

    eventWebSocket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };

    eventWebSocket.onclose = () => {
        console.log('WebSocket closed, reconnecting in 5s...');
        setTimeout(connectEventWebSocket, 5000);
    };
}

function showNotification(event) {
    // Simple notification (could be enhanced with browser notifications)
    const notification = `
        <div class="fixed top-4 right-4 bg-white shadow-lg rounded-lg p-4 max-w-sm animate-fade-in z-50">
            <p class="font-semibold">${event.type.replace(/_/g, ' ').toUpperCase()}</p>
            <p class="text-sm text-gray-600">${new Date(event.timestamp).toLocaleString()}</p>
        </div>
    `;

    const div = document.createElement('div');
    div.innerHTML = notification;
    document.body.appendChild(div);

    setTimeout(() => div.remove(), 5000);
}

// Initialize Application
async function init() {
    if (!authToken) {
        showLogin();
        return;
    }

    await loadSystemStatus();
    await loadCameras();
    await loadEvents();
    connectEventWebSocket();

    // Refresh data periodically
    setInterval(loadSystemStatus, 30000);
    setInterval(loadCameras, 60000);
}

// Start the application
init();

