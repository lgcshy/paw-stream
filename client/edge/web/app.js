/**
 * PawStream Edge Client — Unified App
 */

// ---- State ----
let authToken = '';
let selectedDeviceId = '';
let selectedDeviceSecret = '';
let sseClient = null;
let sysInfoTimer = null;
let startTime = null;

// ---- Setup Flow ----

async function connectServer() {
    const btn = document.getElementById('btn-connect');
    const raw = document.getElementById('api-url').value.trim();
    if (!raw) return showAlert('请输入服务器地址', 'error');

    const url = raw.startsWith('http') ? raw : 'http://' + raw;

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner"></span> 连接中...';

    try {
        const res = await fetch('/api/validate-server', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url })
        });
        const data = await res.json();

        if (!data.valid) throw new Error(data.message || '无法连接');

        showAlert('服务器连接成功', 'success');
        document.getElementById('step-login').style.display = '';
        document.getElementById('step-login').scrollIntoView({ behavior: 'smooth' });
    } catch (e) {
        showAlert('连接失败: ' + e.message, 'error');
    } finally {
        btn.disabled = false;
        btn.textContent = '连接';
    }
}

async function doLogin() {
    const btn = document.getElementById('btn-login');
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;

    if (!username || !password) return showAlert('请输入用户名和密码', 'error');

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner"></span> 登录中...';

    try {
        const res = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err.message || '登录失败');
        }

        const data = await res.json();
        authToken = data.token || '';

        showAlert('登录成功', 'success');
        document.getElementById('step-device').style.display = '';
        document.getElementById('step-device').scrollIntoView({ behavior: 'smooth' });
        loadDevices();
    } catch (e) {
        showAlert('登录失败: ' + e.message, 'error');
    } finally {
        btn.disabled = false;
        btn.textContent = '登录';
    }
}

async function loadDevices() {
    try {
        const res = await fetch('/api/devices', {
            headers: { 'Authorization': 'Bearer ' + authToken }
        });
        if (!res.ok) return;
        const data = await res.json();
        renderDevices(data.devices || data || []);
    } catch (_) { /* non-critical */ }
}

function renderDevices(devices) {
    const grid = document.getElementById('device-list');
    if (!devices.length) {
        grid.innerHTML = '<p style="color:var(--text-dim);font-size:0.9rem;">暂无设备，请创建新设备</p>';
        return;
    }

    grid.innerHTML = devices.map(d => `
        <div class="device-card" data-id="${d.id}" onclick="selectDevice('${d.id}','${d.name || ''}')">
            <div class="name">${escapeHtml(d.name || d.id)}</div>
            <div class="meta">${escapeHtml(d.location || '')} · ${d.enabled ? '已启用' : '未启用'}</div>
        </div>
    `).join('');
}

function selectDevice(id, name) {
    selectedDeviceId = id;
    // Clear new-device fields
    document.getElementById('device-name').value = '';
    document.getElementById('device-location').value = '';

    // Highlight selected card
    document.querySelectorAll('.device-card').forEach(c => c.classList.remove('selected'));
    const card = document.querySelector(`.device-card[data-id="${id}"]`);
    if (card) card.classList.add('selected');

    // Show secret input
    document.getElementById('device-secret-field').style.display = '';
    document.getElementById('device-secret').focus();
}

async function startStreaming() {
    const btn = document.getElementById('btn-start');
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner"></span> 配置中...';

    try {
        let deviceId = selectedDeviceId;
        let deviceSecret = document.getElementById('device-secret').value.trim();
        const newName = document.getElementById('device-name').value.trim();

        // Create new device if name is provided and no existing device selected
        if (newName && !selectedDeviceId) {
            const location = document.getElementById('device-location').value.trim();
            const res = await fetch('/api/devices', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + authToken
                },
                body: JSON.stringify({ name: newName, location })
            });

            if (!res.ok) {
                const err = await res.json().catch(() => ({}));
                throw new Error(err.message || '创建设备失败');
            }

            const data = await res.json();
            deviceId = data.id || data.device?.id;
            deviceSecret = data.secret || data.device?.secret || '';
        }

        if (!deviceId) throw new Error('请选择已有设备或输入新设备名称');
        if (!deviceSecret) throw new Error('请输入设备密钥');

        // Get API URL
        const raw = document.getElementById('api-url').value.trim();
        const apiUrl = raw.startsWith('http') ? raw : 'http://' + raw;

        // Quick setup — saves config with smart defaults & triggers streaming
        const setupRes = await fetch('/api/quick-setup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                api_url: apiUrl,
                device_id: deviceId,
                device_secret: deviceSecret
            })
        });

        if (!setupRes.ok) {
            const err = await setupRes.json().catch(() => ({}));
            throw new Error(err.message || '配置保存失败');
        }

        showDashboard(newName || deviceId, deviceId);
    } catch (e) {
        showAlert(e.message, 'error');
        btn.disabled = false;
        btn.textContent = '开始推流';
    }
}

// ---- Dashboard ----

function showDashboard(deviceName, deviceId) {
    document.getElementById('setup').classList.remove('active');
    document.getElementById('dashboard').classList.add('active');
    document.getElementById('settings-btn').style.display = '';
    document.getElementById('dash-device').textContent = deviceName || 'PawStream';

    // Populate settings fields
    const raw = document.getElementById('api-url').value.trim();
    document.getElementById('cfg-api-url').value = raw.startsWith('http') ? raw : 'http://' + raw;
    document.getElementById('cfg-device-id').value = deviceId || '';

    startTime = Date.now();

    // Start SSE
    initSSE();
    // Start system info polling
    loadSystemInfo();
    sysInfoTimer = setInterval(loadSystemInfo, 5000);
}

function initSSE() {
    if (sseClient) sseClient.disconnect();

    sseClient = new SSEClient('/api/events');

    sseClient.on('connected', () => {
        setBadge(true);
    });

    sseClient.on('status', (status) => {
        updateDashStatus(status);
    });

    sseClient.on('log', (entry) => {
        addLogEntry(entry);
    });

    sseClient.on('config', () => {
        showAlert('配置已更新', 'info');
    });

    sseClient.on('error', () => {
        setBadge(false);
    });

    sseClient.connect();
}

function setBadge(online) {
    const badge = document.getElementById('dash-badge');
    const text = document.getElementById('dash-badge-text');
    badge.className = 'status-badge ' + (online ? 'online' : 'offline');
    text.textContent = online ? '在线' : '离线';
}

function updateDashStatus(status) {
    setBadge(true);

    const uptime = status.uptime || formatUptime(Date.now() - startTime);
    document.getElementById('dash-uptime').textContent = uptime;
    document.getElementById('dash-frames').textContent = status.frames_pushed || '--';
    document.getElementById('dash-stream').textContent = status.stream_status || status.client_status || '--';
}

async function loadSystemInfo() {
    try {
        const res = await fetch('/api/system/info');
        if (!res.ok) return;
        const info = await res.json();

        const cpu = (info.cpu?.usage_total || 0).toFixed(1);
        const mem = (info.memory?.used_percent || 0).toFixed(1);
        const disk = (info.disk?.used_percent || 0).toFixed(1);

        document.getElementById('res-cpu').textContent = cpu + '%';
        document.getElementById('res-mem').textContent = mem + '%';
        document.getElementById('res-disk').textContent = disk + '%';
        document.getElementById('res-cpu-bar').style.width = cpu + '%';
        document.getElementById('res-mem-bar').style.width = mem + '%';
        document.getElementById('res-disk-bar').style.width = disk + '%';
    } catch (_) { /* silent */ }
}

// ---- Logs ----

function addLogEntry(entry) {
    const box = document.getElementById('log-box');
    const time = new Date(entry.timestamp).toLocaleTimeString('zh-CN', { hour12: false });
    const lvl = (entry.level || 'info').toLowerCase();

    const el = document.createElement('div');
    el.className = 'entry';
    el.innerHTML = `<span class="time">${time}</span><span class="lvl lvl-${lvl}">${entry.level || 'INFO'}</span><span class="msg">${escapeHtml(entry.message || '')}</span>`;
    box.appendChild(el);

    // Keep max 500 entries
    while (box.children.length > 500) box.removeChild(box.firstChild);

    if (document.getElementById('auto-scroll').checked) {
        box.scrollTop = box.scrollHeight;
    }
}

function clearLogs() {
    const box = document.getElementById('log-box');
    box.innerHTML = '<div class="entry"><span class="time">--:--:--</span><span class="lvl lvl-info">INFO</span><span class="msg">日志已清空</span></div>';
}

// ---- Settings Drawer ----

function openSettings() {
    document.getElementById('settings-drawer').classList.add('open');
    document.getElementById('settings-overlay').classList.add('show');
}

function closeSettings() {
    document.getElementById('settings-drawer').classList.remove('open');
    document.getElementById('settings-overlay').classList.remove('show');
}

async function saveAdvancedConfig() {
    const config = {
        api_url: document.getElementById('cfg-api-url').value.trim(),
        device_id: document.getElementById('cfg-device-id').value.trim(),
        mediamtx_url: document.getElementById('cfg-mediamtx').value.trim(),
        input_type: document.getElementById('cfg-input-type').value
    };

    try {
        const res = await fetch('/api/config', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config)
        });

        if (!res.ok) throw new Error('保存失败');
        showAlert('设置已保存', 'success');
        closeSettings();
    } catch (e) {
        showAlert('保存失败: ' + e.message, 'error');
    }
}

// ---- Alerts ----

function showAlert(msg, type) {
    const el = document.getElementById('alert');
    el.textContent = msg;
    el.className = 'alert show alert-' + (type || 'info');

    clearTimeout(showAlert._timer);
    showAlert._timer = setTimeout(() => {
        el.classList.remove('show');
    }, 4000);
}

// ---- Utilities ----

function escapeHtml(text) {
    const d = document.createElement('div');
    d.textContent = text;
    return d.innerHTML;
}

function formatUptime(ms) {
    const s = Math.floor(ms / 1000);
    const h = Math.floor(s / 3600);
    const m = Math.floor((s % 3600) / 60);
    if (h > 0) return h + '小时 ' + m + '分';
    if (m > 0) return m + '分 ' + (s % 60) + '秒';
    return s + '秒';
}

// ---- Init: check if already configured ----

document.addEventListener('DOMContentLoaded', async () => {
    try {
        const res = await fetch('/api/config');
        if (!res.ok) return;
        const cfg = await res.json();

        // If device is already configured, go straight to dashboard
        if (cfg.device?.id && cfg.device?.secret) {
            showDashboard(cfg.device.id, cfg.device.id);
        }
    } catch (_) {
        // Fresh start — stay on setup
    }
});
