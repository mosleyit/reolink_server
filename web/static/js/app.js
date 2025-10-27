// API base URL
const API_BASE = '/api/v1';

// Check system health on load
async function checkHealth() {
    try {
        const response = await fetch('/health');
        const data = await response.json();
        
        if (data.success) {
            const statusDiv = document.getElementById('status');
            statusDiv.innerHTML = `
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <p class="text-sm text-gray-500">Status</p>
                        <p class="text-lg font-semibold text-green-600">${data.data.status}</p>
                    </div>
                    <div>
                        <p class="text-sm text-gray-500">Version</p>
                        <p class="text-lg font-semibold">${data.data.version}</p>
                    </div>
                    <div>
                        <p class="text-sm text-gray-500">Uptime</p>
                        <p class="text-lg font-semibold">${data.data.uptime}</p>
                    </div>
                </div>
            `;
        }
    } catch (error) {
        console.error('Failed to check health:', error);
        const statusDiv = document.getElementById('status');
        statusDiv.innerHTML = '<p class="text-red-600">Failed to connect to server</p>';
    }
}

// Load cameras
async function loadCameras() {
    try {
        const response = await fetch(`${API_BASE}/cameras`);
        const data = await response.json();
        
        if (data.success && data.data && data.data.length > 0) {
            const camerasDiv = document.getElementById('cameras');
            camerasDiv.innerHTML = data.data.map(camera => `
                <div class="border rounded-lg p-4">
                    <h3 class="font-semibold text-lg mb-2">${camera.name}</h3>
                    <p class="text-sm text-gray-600">Status: <span class="font-medium ${camera.status === 'online' ? 'text-green-600' : 'text-red-600'}">${camera.status}</span></p>
                    <p class="text-sm text-gray-600">Model: ${camera.model || 'Unknown'}</p>
                    <div class="mt-4 flex space-x-2">
                        <button class="bg-blue-500 hover:bg-blue-700 text-white text-sm py-1 px-3 rounded">
                            View
                        </button>
                        <button class="bg-gray-500 hover:bg-gray-700 text-white text-sm py-1 px-3 rounded">
                            Settings
                        </button>
                    </div>
                </div>
            `).join('');
        }
    } catch (error) {
        console.error('Failed to load cameras:', error);
    }
}

// Load recent events
async function loadEvents() {
    try {
        const response = await fetch(`${API_BASE}/events?limit=10`);
        const data = await response.json();
        
        if (data.success && data.data && data.data.items && data.data.items.length > 0) {
            const eventsDiv = document.getElementById('events');
            eventsDiv.innerHTML = data.data.items.map(event => `
                <div class="border-l-4 ${getEventColor(event.type)} pl-4 py-2">
                    <p class="font-semibold">${event.camera_name} - ${event.type}</p>
                    <p class="text-sm text-gray-600">${new Date(event.timestamp).toLocaleString()}</p>
                </div>
            `).join('');
        }
    } catch (error) {
        console.error('Failed to load events:', error);
    }
}

// Get event color based on type
function getEventColor(type) {
    const colors = {
        'motion_detected': 'border-yellow-500',
        'ai_person': 'border-red-500',
        'ai_vehicle': 'border-orange-500',
        'ai_pet': 'border-blue-500',
        'camera_online': 'border-green-500',
        'camera_offline': 'border-gray-500',
    };
    return colors[type] || 'border-gray-300';
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    checkHealth();
    loadCameras();
    loadEvents();
    
    // Refresh data every 30 seconds
    setInterval(() => {
        checkHealth();
        loadCameras();
        loadEvents();
    }, 30000);
});

