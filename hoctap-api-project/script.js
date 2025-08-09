// API Configuration
const API_BASE_URL = 'http://localhost:8080';

// DOM Elements
const apiStatus = document.getElementById('api-status');
const lastCheck = document.getElementById('last-check');
const responseTime = document.getElementById('response-time');
const checkHealthBtn = document.getElementById('check-health');
const refreshUsersBtn = document.getElementById('refresh-users');
const addUserForm = document.getElementById('add-user-form');
const usersContainer = document.getElementById('users-container');
const responseContainer = document.getElementById('response-container');
const toastContainer = document.getElementById('toast-container');

// State
let users = [];
let apiOnline = false;

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 HocTap API Dashboard initialized');
    
    // Set up event listeners
    setupEventListeners();
    
    // Initial health check
    checkApiHealth();
    
    // Load users
    loadUsers();
    
    // Auto-refresh every 30 seconds
    setInterval(checkApiHealth, 30000);
});

// Event Listeners Setup
function setupEventListeners() {
    checkHealthBtn.addEventListener('click', checkApiHealth);
    refreshUsersBtn.addEventListener('click', loadUsers);
    addUserForm.addEventListener('submit', handleAddUser);
}

// API Health Check
async function checkApiHealth() {
    const startTime = Date.now();
    
    try {
        updateApiStatus('checking', 'Checking API...');
        
        const response = await fetch(`${API_BASE_URL}/health`);
        const data = await response.json();
        
        const responseTimeMs = Date.now() - startTime;
        
        if (response.ok) {
            apiOnline = true;
            updateApiStatus('online', 'API Online');
            updateLastCheck();
            updateResponseTime(responseTimeMs);
            console.log('✅ API Health Check: Online', data);
        } else {
            throw new Error(`HTTP ${response.status}`);
        }
    } catch (error) {
        apiOnline = false;
        updateApiStatus('offline', 'API Offline');
        updateLastCheck();
        updateResponseTime(null);
        console.error('❌ API Health Check: Failed', error);
        showToast('API Error', 'Failed to connect to the API server', 'error');
    }
}

// Update API Status Display
function updateApiStatus(status, text) {
    apiStatus.className = `status-badge ${status}`;
    apiStatus.innerHTML = `<i class="fas fa-circle"></i> ${text}`;
}

// Update Last Check Time
function updateLastCheck() {
    const now = new Date();
    lastCheck.textContent = now.toLocaleTimeString();
}

// Update Response Time
function updateResponseTime(ms) {
    if (ms !== null) {
        responseTime.textContent = `${ms}ms`;
        responseTime.style.color = ms < 100 ? '#48bb78' : ms < 500 ? '#ed8936' : '#f56565';
    } else {
        responseTime.textContent = 'N/A';
        responseTime.style.color = '#718096';
    }
}

// Load Users from API
async function loadUsers() {
    try {
        showLoading(usersContainer);
        
        const response = await fetch(`${API_BASE_URL}/api/users`);
        const data = await response.json();
        
        if (response.ok) {
            users = data.data || [];
            renderUsers();
            console.log('📋 Users loaded:', users);
        } else {
            throw new Error(data.message || `HTTP ${response.status}`);
        }
    } catch (error) {
        console.error('❌ Failed to load users:', error);
        showError(usersContainer, 'Failed to load users');
        showToast('Load Error', 'Failed to load users from API', 'error');
    }
}

// Render Users in the UI
function renderUsers() {
    if (users.length === 0) {
        usersContainer.innerHTML = `
            <div class="no-users">
                <i class="fas fa-users"></i>
                <p>No users found. Add some users to get started!</p>
            </div>
        `;
        return;
    }

    usersContainer.innerHTML = users.map(user => `
        <div class="user-card fade-in">
            <div class="user-info">
                <h4><i class="fas fa-user"></i> ${escapeHtml(user.name)}</h4>
                <p><i class="fas fa-envelope"></i> ${escapeHtml(user.email)}</p>
                <p><i class="fas fa-id-badge"></i> ID: ${user.id}</p>
            </div>
            <div class="user-actions">
                <button class="btn btn-outline btn-small" onclick="getUserDetails(${user.id})">
                    <i class="fas fa-eye"></i> View
                </button>
                <button class="btn btn-secondary btn-small" onclick="editUser(${user.id})">
                    <i class="fas fa-edit"></i> Edit
                </button>
            </div>
        </div>
    `).join('');
}

// Handle Add User Form
async function handleAddUser(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const userData = {
        name: document.getElementById('user-name').value.trim(),
        email: document.getElementById('user-email').value.trim()
    };

    // Validation
    if (!userData.name || !userData.email) {
        showToast('Validation Error', 'Please fill in all fields', 'warning');
        return;
    }

    if (!isValidEmail(userData.email)) {
        showToast('Validation Error', 'Please enter a valid email address', 'warning');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/api/users`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });

        const data = await response.json();

        if (response.ok) {
            // Clear form
            event.target.reset();
            
            // Reload users
            await loadUsers();
            
            showToast('Success', `User "${userData.name}" created successfully!`, 'success');
            console.log('✅ User created:', data.data);
        } else {
            throw new Error(data.message || `HTTP ${response.status}`);
        }
    } catch (error) {
        console.error('❌ Failed to create user:', error);
        showToast('Create Error', 'Failed to create user', 'error');
    }
}

// Get User Details
async function getUserDetails(userId) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/users/${userId}`);
        const data = await response.json();
        
        if (response.ok) {
            displayResponse('GET', `/api/users/${userId}`, data, response.status);
            showToast('Success', `User details loaded for ID: ${userId}`, 'success');
        } else {
            throw new Error(data.message || `HTTP ${response.status}`);
        }
    } catch (error) {
        console.error('❌ Failed to get user details:', error);
        showToast('Error', 'Failed to get user details', 'error');
    }
}

// Edit User (placeholder function)
function editUser(userId) {
    showToast('Info', `Edit functionality for user ID ${userId} - Coming Soon!`, 'warning');
}

// Test API Endpoint
async function testEndpoint(method, path) {
    const startTime = Date.now();
    
    try {
        const response = await fetch(`${API_BASE_URL}${path}`, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            }
        });
        
        const responseTime = Date.now() - startTime;
        const data = await response.json();
        
        displayResponse(method, path, data, response.status, responseTime);
        
        const statusType = response.ok ? 'success' : 'error';
        showToast('API Test', `${method} ${path} - ${response.status}`, statusType);
        
    } catch (error) {
        console.error('❌ API Test failed:', error);
        displayError(method, path, error.message);
        showToast('API Test Failed', error.message, 'error');
    }
}

// Display API Response
function displayResponse(method, path, data, status, responseTime = null) {
    const timestamp = new Date().toISOString();
    const responseInfo = {
        timestamp,
        method,
        path,
        status,
        responseTime: responseTime ? `${responseTime}ms` : null,
        response: data
    };

    responseContainer.innerHTML = `
        <div class="response-header">
            <div class="response-meta">
                <span class="method ${method.toLowerCase()}">${method}</span>
                <span class="path">${path}</span>
                <span class="status status-${Math.floor(status/100)}xx">${status}</span>
                ${responseTime ? `<span class="timing">${responseTime}ms</span>` : ''}
            </div>
            <div class="response-time">${timestamp}</div>
        </div>
        <pre class="response-body">${JSON.stringify(responseInfo, null, 2)}</pre>
    `;
}

// Display API Error
function displayError(method, path, error) {
    const timestamp = new Date().toISOString();
    responseContainer.innerHTML = `
        <div class="response-header error">
            <div class="response-meta">
                <span class="method ${method.toLowerCase()}">${method}</span>
                <span class="path">${path}</span>
                <span class="status error">ERROR</span>
            </div>
            <div class="response-time">${timestamp}</div>
        </div>
        <pre class="response-body error">
{
  "error": "${error}",
  "timestamp": "${timestamp}",
  "method": "${method}",
  "path": "${path}"
}</pre>
    `;
}

// Show Loading State
function showLoading(container) {
    container.innerHTML = `
        <div class="loading">
            <i class="fas fa-spinner fa-spin"></i> Loading...
        </div>
    `;
}

// Show Error State
function showError(container, message) {
    container.innerHTML = `
        <div class="error-state">
            <i class="fas fa-exclamation-triangle"></i>
            <p>${message}</p>
            <button class="btn btn-secondary" onclick="loadUsers()">
                <i class="fas fa-retry"></i> Retry
            </button>
        </div>
    `;
}

// Toast Notification System
function showToast(title, message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    const icon = getToastIcon(type);
    
    toast.innerHTML = `
        <i class="fas ${icon}"></i>
        <div class="toast-content">
            <div class="toast-title">${title}</div>
            <div class="toast-message">${message}</div>
        </div>
        <button class="toast-close">
            <i class="fas fa-times"></i>
        </button>
    `;

    // Add close functionality
    const closeBtn = toast.querySelector('.toast-close');
    closeBtn.addEventListener('click', () => removeToast(toast));

    // Add to container
    toastContainer.appendChild(toast);

    // Trigger animation
    setTimeout(() => toast.classList.add('show'), 100);

    // Auto remove after 5 seconds
    setTimeout(() => removeToast(toast), 5000);
}

// Get Toast Icon
function getToastIcon(type) {
    const icons = {
        success: 'fa-check-circle',
        error: 'fa-exclamation-circle',
        warning: 'fa-exclamation-triangle',
        info: 'fa-info-circle'
    };
    return icons[type] || icons.info;
}

// Remove Toast
function removeToast(toast) {
    toast.classList.remove('show');
    setTimeout(() => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    }, 300);
}

// Utility Functions
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Add some additional CSS for response display
const additionalCSS = `
<style>
.response-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
    border-bottom: 1px solid #4a5568;
    margin-bottom: 15px;
}

.response-meta {
    display: flex;
    align-items: center;
    gap: 10px;
}

.response-meta .status {
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bold;
}

.status-2xx { background: #48bb78; color: white; }
.status-4xx { background: #ed8936; color: white; }
.status-5xx { background: #f56565; color: white; }
.status.error { background: #f56565; color: white; }

.timing {
    color: #a0aec0;
    font-size: 0.8rem;
}

.response-time {
    color: #a0aec0;
    font-size: 0.8rem;
}

.response-body {
    color: #e2e8f0;
    margin: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
}

.response-body.error {
    color: #fed7d7;
}

.no-users {
    text-align: center;
    padding: 40px;
    color: #718096;
}

.no-users i {
    font-size: 3rem;
    margin-bottom: 15px;
    opacity: 0.5;
}

.error-state {
    text-align: center;
    padding: 40px;
    color: #718096;
}

.error-state i {
    font-size: 2rem;
    color: #f56565;
    margin-bottom: 15px;
}
</style>
`;

// Inject additional CSS
document.head.insertAdjacentHTML('beforeend', additionalCSS);

// Export functions for testing (if needed)
window.HocTapAPI = {
    checkApiHealth,
    loadUsers,
    testEndpoint,
    showToast
};
