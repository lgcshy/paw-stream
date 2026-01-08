# PawStream API 使用指南

## 基础信息

- **Base URL**: `http://localhost:3000`
- **认证方式**: JWT Bearer Token
- **Content-Type**: `application/json`

## 认证流程

### 1. 用户注册

```bash
POST /api/register
```

**请求体:**
```json
{
  "username": "testuser",
  "password": "test123",
  "nickname": "Test User"  // 可选,默认为 username
}
```

**响应 (201 Created):**
```json
{
  "id": "4c44ec82-e2a9-40cd-b23b-9e4c8d3ae49f",
  "username": "testuser",
  "nickname": "Test User",
  "disabled": false,
  "created_at": "2026-01-07T09:53:35+08:00",
  "updated_at": "2026-01-07T09:53:35+08:00"
}
```

### 2. 用户登录

```bash
POST /api/login
```

**请求体:**
```json
{
  "username": "testuser",
  "password": "test123"
}
```

**响应 (200 OK):**
```json
{
  "token": "eyJhbGci...",
  "user": {
    "id": "4c44ec82-e2a9-40cd-b23b-9e4c8d3ae49f",
    "username": "testuser",
    "nickname": "Test User",
    "disabled": false,
    "created_at": "2026-01-07T09:53:35+08:00",
    "updated_at": "2026-01-07T09:53:35+08:00"
  }
}
```

### 3. 获取当前用户信息

```bash
GET /api/me
Authorization: Bearer <token>
```

**响应 (200 OK):**
```json
{
  "id": "4c44ec82-e2a9-40cd-b23b-9e4c8d3ae49f",
  "username": "testuser",
  "nickname": "Test User",
  "disabled": false,
  "created_at": "2026-01-07T09:53:35+08:00",
  "updated_at": "2026-01-07T09:53:35+08:00"
}
```

## 设备管理

所有设备管理接口都需要认证 (Bearer Token)。

### 1. 创建设备

```bash
POST /api/devices
Authorization: Bearer <token>
```

**请求体:**
```json
{
  "name": "家里的狗狗摄像头",
  "location": "客厅"  // 可选
}
```

**响应 (201 Created):**
```json
{
  "device": {
    "id": "a67cc67b-5ed4-4f45-8de2-bef7fae71aca",
    "name": "家里的狗狗摄像头",
    "location": "客厅",
    "publish_path": "dogcam/a67cc67b-5ed4-4f45-8de2-bef7fae71aca",
    "disabled": false,
    "created_at": "2026-01-07T09:53:47+08:00",
    "updated_at": "2026-01-07T09:53:47+08:00"
  },
  "secret": "CGG-79J5R4OJiMYf_DM1iZnJrZ_UuZDwKvvpuxvti38="
}
```

⚠️ **重要**: `secret` 只返回一次!请妥善保存用于设备推流。

### 2. 列出设备

```bash
GET /api/devices
Authorization: Bearer <token>
```

**响应 (200 OK):**
```json
[
  {
    "id": "a67cc67b-5ed4-4f45-8de2-bef7fae71aca",
    "name": "家里的狗狗摄像头",
    "location": "客厅",
    "publish_path": "dogcam/a67cc67b-5ed4-4f45-8de2-bef7fae71aca",
    "disabled": false,
    "created_at": "2026-01-07T09:53:47+08:00",
    "updated_at": "2026-01-07T09:53:47+08:00"
  }
]
```

### 3. 获取设备详情

```bash
GET /api/devices/:id
Authorization: Bearer <token>
```

**响应 (200 OK):** 同上单个设备对象

### 4. 更新设备

```bash
PUT /api/devices/:id
Authorization: Bearer <token>
```

**请求体 (所有字段可选):**
```json
{
  "name": "新名称",
  "location": "新位置",
  "disabled": false
}
```

**响应 (200 OK):** 更新后的设备对象

### 5. 删除设备

```bash
DELETE /api/devices/:id
Authorization: Bearer <token>
```

**响应:** 204 No Content

### 6. 轮换设备 Secret

```bash
POST /api/devices/:id/rotate-secret
Authorization: Bearer <token>
```

**响应 (200 OK):**
```json
{
  "secret": "new-secret-here",
  "secret_version": 2
}
```

⚠️ **重要**: 轮换后旧 secret 立即失效!

## 路径查询

### 列出可访问的流路径

```bash
GET /api/paths
Authorization: Bearer <token>
```

**响应 (200 OK):**
```json
[
  {
    "publish_path": "dogcam/a67cc67b-5ed4-4f45-8de2-bef7fae71aca",
    "device_id": "a67cc67b-5ed4-4f45-8de2-bef7fae71aca",
    "device_name": "家里的狗狗摄像头",
    "device_location": "客厅"
  }
]
```

## MediaMTX 鉴权回调

这是内部接口,由 MediaMTX 调用,不需要手动调用。

### 发布鉴权

```bash
POST /mediamtx/auth
```

**请求体:**
```json
{
  "action": "publish",
  "path": "dogcam/a67cc67b-5ed4-4f45-8de2-bef7fae71aca",
  "password": "<device_secret>",
  "protocol": "rtsp",
  "ip": "127.0.0.1"
}
```

**响应:**
- `200 OK` - 允许发布
- `403 Forbidden` - 拒绝发布

### 读取/播放鉴权

```bash
POST /mediamtx/auth
```

**请求体:**
```json
{
  "action": "read",
  "path": "dogcam/a67cc67b-5ed4-4f45-8de2-bef7fae71aca",
  "token": "<user_jwt_token>",
  "protocol": "webrtc",
  "ip": "127.0.0.1"
}
```

**响应:**
- `200 OK` - 允许读取
- `403 Forbidden` - 拒绝读取

## 错误响应格式

所有错误响应遵循统一格式:

```json
{
  "error": "error_code",
  "message": "Human readable error message",
  "request_id": "uuid"
}
```

常见错误码:
- `bad_request` (400) - 请求格式错误
- `unauthorized` (401) - 未认证或 token 无效
- `forbidden` (403) - 无权限
- `device_not_found` (404) - 设备不存在
- `duplicate_username` (409) - 用户名已存在
- `internal_error` (500) - 服务器内部错误

## 完整示例: curl

### 用户注册和登录

```bash
# 1. 注册用户
curl -X POST http://localhost:3000/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo123","nickname":"Demo User"}'

# 2. 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo123"}' | jq -r '.token')

echo "Token: $TOKEN"
```

### 设备管理

```bash
# 3. 创建设备
DEVICE_RESP=$(curl -s -X POST http://localhost:3000/api/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"我的宠物摄像头","location":"卧室"}')

echo "$DEVICE_RESP" | jq .

# 保存 device_id 和 secret
DEVICE_ID=$(echo "$DEVICE_RESP" | jq -r '.device.id')
DEVICE_SECRET=$(echo "$DEVICE_RESP" | jq -r '.secret')

echo "Device ID: $DEVICE_ID"
echo "Device Secret: $DEVICE_SECRET"

# 4. 列出设备
curl -s http://localhost:3000/api/devices \
  -H "Authorization: Bearer $TOKEN" | jq .

# 5. 获取设备详情
curl -s http://localhost:3000/api/devices/$DEVICE_ID \
  -H "Authorization: Bearer $TOKEN" | jq .

# 6. 更新设备
curl -s -X PUT http://localhost:3000/api/devices/$DEVICE_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"location":"客厅"}' | jq .

# 7. 列出可访问路径
curl -s http://localhost:3000/api/paths \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## 完整示例: JavaScript/Fetch

```javascript
// 基础 URL
const BASE_URL = 'http://localhost:3000';

// 1. 注册用户
async function register(username, password, nickname) {
  const response = await fetch(`${BASE_URL}/api/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password, nickname })
  });
  return response.json();
}

// 2. 登录
async function login(username, password) {
  const response = await fetch(`${BASE_URL}/api/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  const data = await response.json();
  return data.token; // 返回 token
}

// 3. 创建设备
async function createDevice(token, name, location) {
  const response = await fetch(`${BASE_URL}/api/devices`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ name, location })
  });
  return response.json();
}

// 4. 列出设备
async function listDevices(token) {
  const response = await fetch(`${BASE_URL}/api/devices`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return response.json();
}

// 5. 列出路径
async function listPaths(token) {
  const response = await fetch(`${BASE_URL}/api/paths`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return response.json();
}

// 使用示例
async function main() {
  // 注册并登录
  await register('demo', 'demo123', 'Demo User');
  const token = await login('demo', 'demo123');
  
  // 创建设备
  const { device, secret } = await createDevice(token, '我的摄像头', '客厅');
  console.log('Device ID:', device.id);
  console.log('Device Secret:', secret); // 只返回一次!
  
  // 列出所有设备
  const devices = await listDevices(token);
  console.log('My devices:', devices);
  
  // 列出可访问路径
  const paths = await listPaths(token);
  console.log('Accessible paths:', paths);
}

main();
```

## 安全注意事项

1. **生产环境必须修改 JWT Secret**
   - 通过环境变量 `PAWSTREAM_JWT_SECRET` 设置
   
2. **设备 Secret 妥善保管**
   - 仅在创建/轮换时返回一次
   - 用于设备推流认证
   
3. **HTTPS 部署**
   - 生产环境必须使用 HTTPS
   - 保护 token 和 secret 传输安全
   
4. **Token 有效期**
   - 默认 24 小时
   - 可通过配置调整

5. **CORS 配置**
   - 生产环境需配置允许的来源
   - 当前默认允许所有来源 (开发用)

## 下一步

查看完整文档:
- **README.md** - 项目概述和配置
- **VERIFICATION_REPORT.md** - 测试报告
- **部署指南** - deployments/ 目录
