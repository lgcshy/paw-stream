/**
 * PawStream Edge Client Web UI
 */

// Global state
const state = {
    config: {},
    status: {},
    logs: [],
    sseClient: null
};

// Initialize app on DOM ready
document.addEventListener('DOMContentLoaded', () => {
    initTabs();
    initConfigForm();
    initSSE();
    loadConfig();
    loadSystemInfo();
    initAutoRefresh();
});

/**
 * Initialize tab switching
 */
function initTabs() {
    const tabBtns = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');

    tabBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const tabName = btn.dataset.tab;

            // Remove active class from all tabs
            tabBtns.forEach(b => b.classList.remove('active'));
            tabContents.forEach(c => c.classList.remove('active'));

            // Add active class to clicked tab
            btn.classList.add('active');
            document.getElementById(tabName).classList.add('active');
        });
    });
}

/**
 * Initialize configuration form
 */
function initConfigForm() {
    const form = document.getElementById('configForm');
    const inputType = document.getElementById('inputType');
    const inputSource = document.getElementById('inputSource');
    const inputSourceHelp = document.getElementById('inputSourceHelp');
    const webuiAuthEnabled = document.getElementById('webuiAuthEnabled');
    const webuiAuthGroup = document.getElementById('webuiAuthGroup');

    // Input type change handler
    inputType.addEventListener('change', () => {
        const type = inputType.value;
        switch (type) {
            case 'testsrc':
                inputSource.placeholder = '留空使用默认';
                inputSourceHelp.textContent = '使用内置测试图案';
                break;
            case 'file':
                inputSource.placeholder = '/path/to/video.mp4';
                inputSourceHelp.textContent = '视频文件的完整路径';
                break;
            case 'v4l2':
                inputSource.placeholder = '/dev/video0';
                inputSourceHelp.textContent = 'V4L2 设备路径';
                break;
            case 'rtsp':
                inputSource.placeholder = 'rtsp://camera-ip:554/stream';
                inputSourceHelp.textContent = 'RTSP 流地址';
                break;
        }
    });

    // Web UI Auth toggle
    webuiAuthEnabled.addEventListener('change', () => {
        webuiAuthGroup.style.display = webuiAuthEnabled.checked ? 'block' : 'none';
    });

    // Password toggle
    document.querySelectorAll('.toggle-password').forEach(btn => {
        btn.addEventListener('click', () => {
            const input = btn.previousElementSibling;
            input.type = input.type === 'password' ? 'text' : 'password';
            btn.textContent = input.type === 'password' ? '👁️' : '🙈';
        });
    });

    // Validate API server
    document.getElementById('validateBtn').addEventListener('click', async () => {
        const apiUrl = document.getElementById('apiUrl').value;
        if (!apiUrl) {
            showNotification('请输入 API 服务器地址', 'warning');
            return;
        }

        try {
            const response = await fetch('/api/validate-server', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url: apiUrl })
            });

            const result = await response.json();
            if (result.valid) {
                showNotification('API 服务器连接成功！', 'success');
            } else {
                showNotification(`连接失败: ${result.message}`, 'error');
            }
        } catch (err) {
            showNotification(`验证失败: ${err.message}`, 'error');
        }
    });

    // Form submission
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveConfig();
    });

    // Reload button
    document.getElementById('reloadBtn').addEventListener('click', () => {
        loadConfig();
    });

    // Clear logs
    document.getElementById('clearLogsBtn').addEventListener('click', () => {
        state.logs = [];
        updateLogDisplay();
    });
}

/**
 * Initialize SSE connection
 */
function initSSE() {
    state.sseClient = new SSEClient('/api/events');

    state.sseClient.on('connected', () => {
        updateConnectionStatus(true);
    });

    state.sseClient.on('status', (status) => {
        state.status = status;
        updateStatusDisplay();
    });

    state.sseClient.on('log', (log) => {
        addLog(log);
    });

    state.sseClient.on('config', () => {
        showNotification('配置已更新', 'info');
        loadConfig();
    });

    state.sseClient.on('error', () => {
        updateConnectionStatus(false);
    });

    state.sseClient.connect();
}

/**
 * Load configuration from server
 */
async function loadConfig() {
    try {
        const response = await fetch('/api/config');
        if (!response.ok) {
            throw new Error('Failed to load config');
        }

        const config = await response.json();
        state.config = config;
        populateConfigForm(config);
    } catch (err) {
        console.error('Failed to load config:', err);
        showNotification('加载配置失败', 'error');
    }
}

/**
 * Save configuration to server
 */
async function saveConfig() {
    const form = document.getElementById('configForm');
    const formData = new FormData(form);
    
    // Convert flat form data to nested object
    const config = {};
    for (const [key, value] of formData.entries()) {
        const keys = key.split('.');
        let current = config;
        
        for (let i = 0; i < keys.length - 1; i++) {
            if (!current[keys[i]]) {
                current[keys[i]] = {};
            }
            current = current[keys[i]];
        }
        
        current[keys[keys.length - 1]] = value;
    }

    // Add checkbox values
    config.webui = config.webui || {};
    config.webui.auth = config.webui.auth || {};
    config.webui.auth.enabled = document.getElementById('webuiAuthEnabled').checked;

    try {
        const response = await fetch('/api/config', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config)
        });

        if (!response.ok) {
            throw new Error('Failed to save config');
        }

        const result = await response.json();
        showNotification('配置保存成功！将自动重载...', 'success');
        state.config = config;
    } catch (err) {
        console.error('Failed to save config:', err);
        showNotification('保存配置失败', 'error');
    }
}

/**
 * Populate form with config data
 */
function populateConfigForm(config) {
    // API settings
    setFormValue('apiUrl', config.api?.url);
    
    // Device settings
    setFormValue('deviceId', config.device?.id);
    setFormValue('deviceSecret', config.device?.secret);
    
    // Input settings
    setFormValue('inputType', config.input?.type);
    setFormValue('inputSource', config.input?.source);
    
    // MediaMTX settings
    setFormValue('mediamtxUrl', config.mediamtx?.url);
    
    // Web UI settings
    setFormValue('webuiPort', config.webui?.port);
    
    const authEnabled = config.webui?.auth?.enabled || false;
    document.getElementById('webuiAuthEnabled').checked = authEnabled;
    document.getElementById('webuiAuthGroup').style.display = authEnabled ? 'block' : 'none';
    
    setFormValue('webuiAuthUser', config.webui?.auth?.username);
    setFormValue('webuiAuthPass', config.webui?.auth?.password);

    // Trigger input type change to update help text
    document.getElementById('inputType').dispatchEvent(new Event('change'));
}

/**
 * Set form field value
 */
function setFormValue(id, value) {
    const element = document.getElementById(id);
    if (element && value !== undefined && value !== null) {
        element.value = value;
    }
}

/**
 * Load system information
 */
async function loadSystemInfo() {
    try {
        const response = await fetch('/api/system/info');
        if (!response.ok) {
            throw new Error('Failed to load system info');
        }

        const info = await response.json();
        updateSystemInfo(info);
    } catch (err) {
        console.error('Failed to load system info:', err);
    }
}

/**
 * Update system information display
 */
function updateSystemInfo(info) {
    // Host info
    setText('sysHostname', info.host?.hostname || '-');
    setText('sysOS', `${info.host?.platform || '-'} ${info.host?.platform_version || ''}`);
    setText('sysUptime', info.host?.uptime_formatted || '-');

    // CPU
    const cpuUsage = info.cpu?.usage_total || 0;
    setText('cpuUsage', `${cpuUsage.toFixed(1)}%`);
    setProgressBar('cpuProgress', cpuUsage);

    // Memory
    const memPercent = info.memory?.used_percent || 0;
    const memUsed = formatBytes(info.memory?.used || 0);
    const memTotal = formatBytes(info.memory?.total || 0);
    setText('memUsage', `${memUsed} / ${memTotal} (${memPercent.toFixed(1)}%)`);
    setProgressBar('memProgress', memPercent);

    // Disk
    const diskPercent = info.disk?.used_percent || 0;
    const diskUsed = formatBytes(info.disk?.used || 0);
    const diskTotal = formatBytes(info.disk?.total || 0);
    setText('diskUsage', `${diskUsed} / ${diskTotal} (${diskPercent.toFixed(1)}%)`);
    setProgressBar('diskProgress', diskPercent);
}

/**
 * Update status display
 */
function updateStatusDisplay() {
    const status = state.status;
    
    setText('clientStatus', status.client_status || '-');
    setText('streamStatus', status.stream_status || '-');
    setText('uptime', status.uptime || '-');
    setText('framesPushed', status.frames_pushed || '0');
}

/**
 * Update connection status indicator
 */
function updateConnectionStatus(connected) {
    const statusDot = document.getElementById('statusDot');
    const statusText = document.getElementById('statusText');

    if (connected) {
        statusDot.classList.add('connected');
        statusDot.classList.remove('error');
        statusText.textContent = '已连接';
    } else {
        statusDot.classList.remove('connected');
        statusDot.classList.add('error');
        statusText.textContent = '连接断开';
    }
}

/**
 * Add log entry
 */
function addLog(log) {
    state.logs.push(log);
    
    // Keep only last 500 logs
    if (state.logs.length > 500) {
        state.logs = state.logs.slice(-500);
    }

    updateLogDisplay();
}

/**
 * Update log display
 */
function updateLogDisplay() {
    const container = document.getElementById('logContainer');
    const autoScroll = document.getElementById('autoScrollCheck').checked;

    if (state.logs.length === 0) {
        container.innerHTML = `
            <div class="log-entry">
                <span class="log-time">--:--:--</span>
                <span class="log-level level-info">INFO</span>
                <span class="log-message">暂无日志</span>
            </div>
        `;
        return;
    }

    container.innerHTML = state.logs.map(log => {
        const time = new Date(log.timestamp).toLocaleTimeString('zh-CN', { hour12: false });
        const levelClass = `level-${log.level.toLowerCase()}`;
        
        return `
            <div class="log-entry">
                <span class="log-time">${time}</span>
                <span class="log-level ${levelClass}">${log.level}</span>
                <span class="log-message">${escapeHtml(log.message)}</span>
            </div>
        `;
    }).join('');

    if (autoScroll) {
        container.scrollTop = container.scrollHeight;
    }
}

/**
 * Initialize auto-refresh for system info
 */
function initAutoRefresh() {
    // Refresh system info every 5 seconds
    setInterval(() => {
        loadSystemInfo();
    }, 5000);
}

/**
 * Helper: Set text content
 */
function setText(id, text) {
    const element = document.getElementById(id);
    if (element) {
        element.textContent = text;
    }
}

/**
 * Helper: Set progress bar
 */
function setProgressBar(id, percent) {
    const element = document.getElementById(id);
    if (element) {
        element.style.width = `${Math.min(percent, 100)}%`;
        
        // Change color based on usage
        if (percent > 90) {
            element.style.background = 'var(--danger-color)';
        } else if (percent > 70) {
            element.style.background = 'var(--warning-color)';
        } else {
            element.style.background = 'var(--primary-color)';
        }
    }
}

/**
 * Helper: Format bytes
 */
function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * Helper: Escape HTML
 */
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

/**
 * Show notification (simple alert for now)
 */
function showNotification(message, type = 'info') {
    // TODO: Implement better notification system
    console.log(`[${type.toUpperCase()}]`, message);
    alert(message);
}
