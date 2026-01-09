// Setup Wizard JavaScript

// Global state
const setupState = {
    currentStep: 1,
    maxStep: 4,  // Now only 4 steps
    apiURL: '',
    token: '',
    selectedDevice: null,
    devices: []
};

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    console.log('[Setup] Wizard initialized');
    
    // Load saved API URL if exists
    const savedAPIURL = localStorage.getItem('pawstream_api_url');
    if (savedAPIURL) {
        document.getElementById('api-url').value = savedAPIURL;
    }
    
    // Setup device form listeners
    setupDeviceFormListeners();
});

// Step navigation
function goToStep(step) {
    // Hide all steps
    document.querySelectorAll('.step-content').forEach(el => {
        el.classList.remove('active');
    });
    
    // Show target step
    document.getElementById(`step-${step}`).classList.add('active');
    
    // Update step indicators
    document.querySelectorAll('.step-item').forEach(el => {
        const stepNum = parseInt(el.dataset.step);
        el.classList.remove('active', 'completed');
        if (stepNum === step) {
            el.classList.add('active');
        } else if (stepNum < step) {
            el.classList.add('completed');
        }
    });
    
    setupState.currentStep = step;
    console.log(`[Setup] Moved to step ${step}`);
}

function nextStep() {
    goToStep(setupState.currentStep + 1);
}

function prevStep() {
    goToStep(setupState.currentStep - 1);
}

// Utility functions
function showAlert(stepNum, type, message) {
    const alertDiv = document.getElementById(`step${stepNum}-alert`);
    alertDiv.innerHTML = `<div class="alert alert-${type}">${message}</div>`;
}

function clearAlert(stepNum) {
    const alertDiv = document.getElementById(`step${stepNum}-alert`);
    alertDiv.innerHTML = '';
}

function setButtonLoading(buttonId, loading, text = '') {
    const btn = document.getElementById(buttonId);
    if (!btn) {
        console.warn(`[Setup] Button with id "${buttonId}" not found`);
        return;
    }
    if (loading) {
        btn.disabled = true;
        btn.innerHTML = `<span class="spinner"></span> ${text || '处理中...'}`;
    } else {
        btn.disabled = false;
        btn.textContent = text || '继续';
    }
}

// Step 1: Validate API server
async function validateServer() {
    const apiURL = document.getElementById('api-url').value.trim();
    
    if (!apiURL) {
        showAlert(1, 'error', '请输入 API 服务器地址');
        return;
    }
    
    clearAlert(1);
    setButtonLoading('btn-validate', true, '验证中...');
    
    try {
        const response = await fetch('/api/validate-server', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url: apiURL })
        });
        
        const data = await response.json();
        
        if (data.valid) {
            setupState.apiURL = apiURL;
            localStorage.setItem('pawstream_api_url', apiURL);
            showAlert(1, 'success', `✅ ${data.message}${data.version ? ` (版本: ${data.version})` : ''}`);
            
            setTimeout(() => {
                nextStep();
                clearAlert(1);
            }, 1000);
        } else {
            showAlert(1, 'error', `❌ ${data.message || '无法连接到 API 服务器'}`);
        }
    } catch (error) {
        console.error('[Setup] Validation error:', error);
        showAlert(1, 'error', `❌ 连接失败: ${error.message}`);
    } finally {
        setButtonLoading('btn-validate', false, '验证连接');
    }
}

// Step 2: User login
async function login() {
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;
    
    if (!username || !password) {
        showAlert(2, 'error', '请输入用户名和密码');
        return;
    }
    
    clearAlert(2);
    setButtonLoading('btn-login', true, '登录中...');
    
    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        
        const data = await response.json();
        
        if (response.ok && data.token) {
            setupState.token = data.token;
            showAlert(2, 'success', '✅ 登录成功');
            
            setTimeout(async () => {
                nextStep();
                clearAlert(2);
                await loadDevices();
            }, 500);
        } else {
            showAlert(2, 'error', `❌ ${data.message || '登录失败'}`);
        }
    } catch (error) {
        console.error('[Setup] Login error:', error);
        showAlert(2, 'error', `❌ 登录失败: ${error.message}`);
    } finally {
        setButtonLoading('btn-login', false, '登录');
    }
}

// Step 3: Setup device form listeners to deselect existing device when creating new one
function setupDeviceFormListeners() {
    const deviceNameInput = document.getElementById('device-name');
    const deviceLocationInput = document.getElementById('device-location');
    
    if (!deviceNameInput || !deviceLocationInput) return;
    
    const clearSelection = () => {
        // Deselect all devices
        document.querySelectorAll('.device-item').forEach(el => {
            el.classList.remove('selected');
        });
        setupState.selectedDevice = null;
        // Hide secret input
        const secretContainer = document.getElementById('device-secret-container');
        if (secretContainer) {
            secretContainer.classList.add('hidden');
        }
        document.getElementById('device-secret').value = '';
    };
    
    deviceNameInput.addEventListener('focus', clearSelection);
    deviceLocationInput.addEventListener('focus', clearSelection);
}

// Step 3: Load devices
async function loadDevices() {
    try {
        const response = await fetch('/api/devices', {
            headers: {
                'Authorization': `Bearer ${setupState.token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load devices');
        }
        
        const data = await response.json();
        // API returns array directly or {devices: [...]}
        setupState.devices = Array.isArray(data) ? data : (data.devices || []);
        
        if (setupState.devices.length > 0) {
            renderDeviceList();
            document.getElementById('device-list-container').classList.remove('hidden');
        }
        
        console.log(`[Setup] Loaded ${setupState.devices.length} devices`);
    } catch (error) {
        console.error('[Setup] Failed to load devices:', error);
        showAlert(3, 'info', '无法加载现有设备，您可以创建新设备');
    }
}

function renderDeviceList() {
    const listEl = document.getElementById('device-list');
    listEl.innerHTML = '';
    
    setupState.devices.forEach(device => {
        const item = document.createElement('div');
        item.className = 'device-item';
        item.innerHTML = `
            <div class="device-name">${device.name || device.id}</div>
            <div class="device-info">
                ${device.location ? `📍 ${device.location}<br>` : ''}
                ID: ${device.id}
            </div>
        `;
        
        item.addEventListener('click', () => {
            // Deselect all
            document.querySelectorAll('.device-item').forEach(el => {
                el.classList.remove('selected');
            });
            // Select this one
            item.classList.add('selected');
            setupState.selectedDevice = device;
            // Show secret input when selecting existing device
            document.getElementById('device-secret-container').classList.remove('hidden');
            console.log('[Setup] Selected device:', device.id);
        });
        
        listEl.appendChild(item);
    });
}

async function selectOrCreateDevice() {
    const deviceName = document.getElementById('device-name').value.trim();
    const deviceLocation = document.getElementById('device-location').value.trim();
    
    // If user selected existing device
    if (setupState.selectedDevice) {
        const deviceSecret = document.getElementById('device-secret').value.trim();
        
        if (!deviceSecret) {
            showAlert(3, 'error', '请输入设备密钥（secret）');
            return;
        }
        
        // Add secret to the selected device
        setupState.selectedDevice.secret = deviceSecret;
        
        clearAlert(3);
        showAlert(3, 'success', `✅ 已选择设备: ${setupState.selectedDevice.name}`);
        
        setTimeout(() => {
            nextStep();
            clearAlert(3);
            detectSources();
        }, 1000);
        return;
    }
    
    // Create new device
    if (!deviceName) {
        showAlert(3, 'error', '请输入设备名称');
        return;
    }
    
    clearAlert(3);
    setButtonLoading('btn-device', true, '创建中...');
    
    try {
        const response = await fetch('/api/devices', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${setupState.token}`
            },
            body: JSON.stringify({
                name: deviceName,
                location: deviceLocation
            })
        });
        
        const data = await response.json();
        
        if (response.ok && data.device && data.device.id) {
            setupState.selectedDevice = {
                id: data.device.id,
                name: data.device.name || deviceName,
                location: data.device.location || deviceLocation,
                secret: data.secret,
                publish_path: data.device.publish_path
            };
            
            showAlert(3, 'success', `✅ 设备创建成功: ${data.device.name}`);
            
            setTimeout(() => {
                nextStep();
                clearAlert(3);
                detectSources();
            }, 1000);
        } else {
            showAlert(3, 'error', `❌ ${data.message || '创建设备失败'}`);
        }
    } catch (error) {
        console.error('[Setup] Create device error:', error);
        showAlert(3, 'error', `❌ 创建设备失败: ${error.message}`);
    } finally {
        setButtonLoading('btn-device', false, '下一步');
    }
}

// Step 4: Detect and select input sources
async function saveAndFinish() {
    clearAlert(4);
    setButtonLoading('btn-finish', true, '保存中...');
    
    // Prepare simplified configuration with smart defaults
    const config = {
        device_id: setupState.selectedDevice.id,
        device_secret: setupState.selectedDevice.secret,
        api_url: setupState.apiURL,
        // Smart defaults - will be auto-detected by backend
        input_type: 'auto',           // Auto-detect: v4l2 > test
        input_source: 'auto',          // Auto-detect
        stream_engine: 'gstreamer',    // Default to GStreamer
        mediamtx_url: 'rtsp://localhost:8554'
    };
    
    try {
        const response = await fetch('/api/config', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config)
        });
        
        const data = await response.json();
        
        if (response.ok || data.success) {
            console.log('[Setup] Configuration saved successfully');
            
            // Update summary
            document.getElementById('summary-api-url').textContent = setupState.apiURL;
            document.getElementById('summary-device-id').textContent = setupState.selectedDevice.id;
            document.getElementById('summary-device-name').textContent = setupState.selectedDevice.name;
            document.getElementById('summary-stream-path').textContent = `dogcam/${setupState.selectedDevice.id}`;
            document.getElementById('summary-config-path').textContent = data.path || './config.yaml';
            
            showAlert(4, 'success', '✅ 配置已保存！客户端将自动检测输入源并开始推流。');
            
            // Change button to go to management interface
            setButtonLoading('btn-finish', false, '进入管理界面');
            document.getElementById('btn-finish').onclick = () => {
                window.location.href = '/';
            };
        } else {
            showAlert(4, 'error', `❌ 保存失败: ${data.message || '未知错误'}`);
            setButtonLoading('btn-finish', false, '完成并开始推流');
        }
    } catch (error) {
        console.error('[Setup] Save config error:', error);
        showAlert(4, 'error', `❌ 保存失败: ${error.message}`);
        setButtonLoading('btn-finish', false, '完成并开始推流');
    }
}


// ==================== Keyboard Shortcuts ====================

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
        const step = setupState.currentStep;
        if (step === 1) validateServer();
        else if (step === 2) login();
        else if (step === 3) selectOrCreateDevice();
        else if (step === 4) saveAndFinish();
    }
});
