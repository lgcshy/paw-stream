// Setup Wizard JavaScript

// Global state
const setupState = {
    currentStep: 1,
    apiURL: '',
    token: '',
    selectedDevice: null,
    selectedSource: null,
    selectedEngine: 'ffmpeg',
    selectedPreset: null,
    devices: [],
    sources: [],
    engines: [],
    presets: [],
    encoders: {}
};

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    console.log('[Setup] Wizard initialized');
    
    // Load saved API URL if exists
    const savedAPIURL = localStorage.getItem('pawstream_api_url');
    if (savedAPIURL) {
        document.getElementById('api-url').value = savedAPIURL;
    }
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
        showAlert(3, 'info', '⚠️ 使用现有设备将需要您提供设备密钥（secret）。是否继续？');
        
        // For existing device, we need to prompt for secret
        // For simplicity, we'll skip this and require creating a new device
        showAlert(3, 'error', '暂不支持使用现有设备，请创建新设备');
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
async function detectSources() {
    clearAlert(4);
    
    try {
        const response = await fetch('/api/input-sources');
        const data = await response.json();
        
        setupState.sources = data.sources || data || [];
        renderSourceList();
        
        console.log(`[Setup] Detected ${setupState.sources.length} sources`);
        showAlert(4, 'success', `✅ 检测到 ${setupState.sources.length} 个输入源`);
        
        // Show custom source option
        document.getElementById('custom-source-container').classList.remove('hidden');
    } catch (error) {
        console.error('[Setup] Source detection error:', error);
        showAlert(4, 'error', `检测失败: ${error.message}`);
        
        // Still show custom source option
        document.getElementById('custom-source-container').classList.remove('hidden');
    }
}

function renderSourceList() {
    const listEl = document.getElementById('source-list');
    listEl.innerHTML = '';
    
    setupState.sources.forEach(source => {
        const card = document.createElement('div');
        card.className = 'source-card';
        card.innerHTML = `
            <span class="source-type ${source.type}">${source.type.toUpperCase()}</span>
            <div class="source-name">${source.name || source.type}</div>
            <div class="source-device">${source.device || source.source}</div>
            ${source.description ? `<div class="source-description">${source.description}</div>` : ''}
        `;
        
        card.addEventListener('click', () => {
            // Deselect all
            document.querySelectorAll('.source-card').forEach(el => {
                el.classList.remove('selected');
            });
            // Select this one
            card.classList.add('selected');
            setupState.selectedSource = source;
            console.log('[Setup] Selected source:', source);
        });
        
        listEl.appendChild(card);
    });
}

async function saveAndFinish() {
    // Check if source is selected or custom source is provided
    let inputType, inputSource;
    
    if (setupState.selectedSource) {
        inputType = setupState.selectedSource.type;
        inputSource = setupState.selectedSource.device || setupState.selectedSource.source;
    } else {
        // Check custom source
        inputType = document.getElementById('custom-source-type').value;
        inputSource = document.getElementById('custom-source').value.trim();
        
        if (!inputSource) {
            showAlert(4, 'error', '请选择一个输入源或输入自定义源');
            return;
        }
    }
    
    clearAlert(4);
    setButtonLoading('btn-finish', true, '保存中...');
    
    // Prepare configuration
    const config = {
        device_id: setupState.selectedDevice.id,
        device_secret: setupState.selectedDevice.secret,
        api_url: setupState.apiURL,
        input_type: inputType,
        input_source: inputSource,
        mediamtx_url: 'rtsp://localhost:8554',
        stream_engine: setupState.selectedEngine,
        stream_preset: setupState.selectedPreset || ''
    };
    
    // Add engine-specific configuration
    if (setupState.selectedEngine === 'ffmpeg') {
        config.ffmpeg_preset = document.getElementById('ffmpeg-preset').value;
        config.ffmpeg_tune = document.getElementById('ffmpeg-tune').value;
        config.ffmpeg_hwaccel = document.getElementById('ffmpeg-hwaccel').value;
    } else if (setupState.selectedEngine === 'gstreamer') {
        config.gstreamer_latency_ms = parseInt(document.getElementById('gst-latency').value);
        config.gstreamer_use_hardware = document.getElementById('gst-hardware').checked;
    }
    
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
            document.getElementById('summary-input-type').textContent = inputType;
            document.getElementById('summary-input-source').textContent = inputSource;
            document.getElementById('summary-engine').textContent = setupState.selectedEngine;
            document.getElementById('summary-preset').textContent = setupState.selectedPreset || '自定义配置';
            document.getElementById('summary-config-path').textContent = data.path || './config.yaml';
            
            setTimeout(() => {
                nextStep();
            }, 500);
        } else {
            showAlert(4, 'error', `❌ 保存失败: ${data.message || '未知错误'}`);
        }
    } catch (error) {
        console.error('[Setup] Save config error:', error);
        showAlert(4, 'error', `❌ 保存失败: ${error.message}`);
    } finally {
        setButtonLoading('btn-finish', false, '完成配置');
    }
}

// ==================== Engine and Preset Functions ====================

// Load available engines
async function loadEngines() {
    try {
        const response = await fetch('/api/engines/available');
        const data = await response.json();
        setupState.engines = data.engines || [];
        
        displayEngines();
    } catch (error) {
        console.error('[Setup] Load engines error:', error);
        showAlert(5, 'error', `❌ 加载引擎列表失败: ${error.message}`);
    }
}

// Load available presets
async function loadPresets() {
    try {
        const response = await fetch('/api/presets');
        const data = await response.json();
        setupState.presets = data.presets || [];
        
        displayPresets();
    } catch (error) {
        console.error('[Setup] Load presets error:', error);
        showAlert(5, 'error', `❌ 加载预设列表失败: ${error.message}`);
    }
}

// Display engines
function displayEngines() {
    const container = document.getElementById('engines-list');
    container.innerHTML = '';
    
    setupState.engines.forEach(engine => {
        const card = document.createElement('div');
        card.className = `engine-card ${engine.available ? '' : 'unavailable'} ${engine.name === setupState.selectedEngine ? 'selected' : ''}`;
        
        if (engine.available) {
            card.onclick = () => selectEngine(engine.name);
        }
        
        card.innerHTML = `
            <div class="engine-header">
                <div class="engine-name">${engine.displayName}</div>
                <div class="engine-badge ${engine.available ? 'badge-available' : 'badge-unavailable'}">
                    ${engine.available ? '✓ 可用' : '✗ 不可用'}
                </div>
            </div>
            <div class="engine-description">${engine.description}</div>
            <div class="engine-features" style="margin-top: 0.5rem;">
                ${engine.features.map(f => `<span style="font-size: 0.875rem; color: var(--text-secondary);">• ${f}</span>`).join('<br>')}
            </div>
            ${!engine.available && engine.installCommand ? `
                <div style="margin-top: 0.5rem; padding: 0.5rem; background: var(--warning-bg); border-radius: 4px; font-size: 0.875rem;">
                    <strong>安装命令:</strong><br>
                    <code style="display: block; margin-top: 0.25rem;">${engine.installCommand}</code>
                </div>
            ` : ''}
        `;
        
        container.appendChild(card);
    });
}

// Display presets
function displayPresets() {
    const container = document.getElementById('presets-list');
    container.innerHTML = '';
    
    // Add "No preset" option
    const noPresetCard = document.createElement('div');
    noPresetCard.className = `preset-card ${!setupState.selectedPreset ? 'selected' : ''}`;
    noPresetCard.onclick = () => selectPreset(null);
    noPresetCard.innerHTML = `
        <div class="preset-header">
            <div class="preset-name">⚙️ 自定义配置</div>
        </div>
        <div class="preset-description">不使用预设，手动配置所有参数</div>
    `;
    container.appendChild(noPresetCard);
    
    setupState.presets.forEach(preset => {
        const card = document.createElement('div');
        card.className = `preset-card ${preset.id === setupState.selectedPreset ? 'selected' : ''}`;
        card.onclick = () => selectPreset(preset.id);
        
        card.innerHTML = `
            <div class="preset-header">
                <div style="display: flex; align-items: center;">
                    <span class="preset-icon">${preset.icon}</span>
                    <div>
                        <div class="preset-name">${preset.name}</div>
                        <div style="font-size: 0.875rem; color: var(--text-secondary);">${preset.description}</div>
                    </div>
                </div>
                ${preset.id === 'low-latency' ? '<div class="engine-badge badge-recommended">推荐</div>' : ''}
            </div>
            <div class="preset-meta">
                <div class="preset-meta-item">📡 引擎: ${preset.engine}</div>
                <div class="preset-meta-item">⏱️ 延迟: ${preset.latency}</div>
                <div class="preset-meta-item">🎨 质量: ${preset.quality}</div>
                <div class="preset-meta-item">💻 资源: ${preset.resource}</div>
            </div>
            <div style="margin-top: 0.5rem; font-size: 0.875rem; color: var(--text-secondary);">
                适用场景: ${preset.scenario}
            </div>
        `;
        
        container.appendChild(card);
    });
}

// Select engine
function selectEngine(engineName) {
    setupState.selectedEngine = engineName;
    displayEngines();
    updateEngineConfig();
}

// Select preset
function selectPreset(presetId) {
    setupState.selectedPreset = presetId;
    displayPresets();
    
    // If preset selected, update engine selection
    if (presetId) {
        const preset = setupState.presets.find(p => p.id === presetId);
        if (preset) {
            setupState.selectedEngine = preset.engine;
            displayEngines();
        }
    }
}

// Update engine configuration panel visibility
function updateEngineConfig() {
    const ffmpegConfig = document.getElementById('ffmpeg-config');
    const gstreamerConfig = document.getElementById('gstreamer-config');
    
    if (setupState.selectedEngine === 'ffmpeg') {
        ffmpegConfig.classList.remove('hidden');
        gstreamerConfig.classList.add('hidden');
    } else if (setupState.selectedEngine === 'gstreamer') {
        ffmpegConfig.classList.add('hidden');
        gstreamerConfig.classList.remove('hidden');
    }
}

// Override nextStep for step 4 to load engines and presets
const originalNextStep = nextStep;
nextStep = function() {
    if (setupState.currentStep === 4) {
        // Load engines and presets when entering step 5
        loadEngines();
        loadPresets();
    }
    originalNextStep();
};

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
        const step = setupState.currentStep;
        if (step === 1) validateServer();
        else if (step === 2) login();
        else if (step === 3) selectOrCreateDevice();
        else if (step === 4) nextStep(); // Changed from saveAndFinish
        else if (step === 5) saveAndFinish();
    }
});
