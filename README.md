# PawStream

[English](README.en.md) | 中文

自托管的家庭宠物实时视频监控系统。

多台边缘设备（树莓派 / 香橙派 + USB 摄像头）将视频推流至中心 [MediaMTX](https://github.com/bluenviron/mediamtx) 媒体服务器，用户通过移动端优先的 Web UI 以 WebRTC 实时观看画面。

## 系统架构

```
[边缘设备 + 摄像头] ──RTSP/RTMP──▶ MediaMTX ──WebRTC/HLS──▶ 手机浏览器
                                       ▲
                                  API 服务器
                                (认证 · 设备管理)
```

- **边缘客户端** (Go) — 采集视频，使用设备凭证推流到 MediaMTX
- **API 服务器** (Go) — 控制面：用户认证、设备管理、MediaMTX 认证回调
- **MediaMTX** — 处理所有媒体传输和协议转换
- **Web UI** (Vue 3) — 移动端优先的流媒体播放界面

## 技术栈

| 层级 | 技术 |
|------|------|
| 边缘客户端 | Go 1.23, FFmpeg / GStreamer |
| API 服务器 | Go 1.24, Fiber v2, SQLite (纯 Go, 无 CGO) |
| 媒体服务器 | MediaMTX (Docker) |
| 前端 | Vue 3, Vite 7, Vant, Pinia |
| 播放协议 | WebRTC (WHEP), HLS 备选 |

## 快速开始

### 环境要求

- Go 1.23+
- Node.js 18+
- Docker
- FFmpeg

### 1. 启动 MediaMTX

```bash
# 将 <YOUR_IP> 替换为本机局域网 IP
docker run --rm -d \
  -e MTX_AUTHMETHOD=http \
  -e MTX_AUTHHTTPADDRESS=http://<YOUR_IP>:3000/mediamtx/auth \
  -e MTX_WEBRTCADDITIONALHOSTS=<YOUR_IP> \
  -p 8554:8554 -p 1935:1935 -p 8888:8888 -p 8889:8889 \
  -p 8890:8890/udp -p 8189:8189/udp \
  --name mediamtx bluenviron/mediamtx:latest-ffmpeg
```

### 2. 启动 API 服务器

```bash
cd server/api
cp config.example.yaml config.yaml  # 编辑 mediamtx 相关 URL 为本机 IP
go build -o api ./cmd/api && ./api
```

### 3. 创建用户和设备

```bash
# 注册用户
curl -X POST http://localhost:3000/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 登录获取 Token
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

# 创建设备（记录返回的 id 和 secret）
curl -X POST http://localhost:3000/api/devices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"my-cam","location":"living-room"}'
```

### 4. 推送测试流

```bash
# 将 DEVICE_ID、DEVICE_SECRET、PUBLISH_PATH 替换为上一步返回的值
ffmpeg -re -f lavfi -i "testsrc=size=1280x720:rate=30" \
  -c:v libx264 -preset ultrafast -tune zerolatency -b:v 2000k -g 30 -keyint_min 30 \
  -f rtsp -rtsp_transport tcp \
  "rtsp://DEVICE_ID:DEVICE_SECRET@localhost:8554/PUBLISH_PATH"
```

### 5. 启动 Web 前端

```bash
cd web
npm install && npm run dev
```

打开 http://localhost:5173，登录后即可观看视频流。

## 项目结构

```
client/edge/       Go 边缘客户端（视频采集 + 推流）
server/api/        Go API 服务器（认证、设备管理）
server/mediamtx/   MediaMTX 配置文件
web/               Vue 3 前端
docs/              架构文档与设计决策
openspec/          变更提案与规格说明
```

## 关键设计决策

- **媒体与控制分离** — API 服务器不接触视频数据，MediaMTX 负责所有媒体传输
- **开发无需真实硬件** — 所有组件均支持 FFmpeg 测试画面和模拟视频源
- **边缘设备无状态** — 配置驱动，可随时替换
- **纯 Go SQLite** — 使用 `modernc.org/sqlite`，无 CGO 依赖
- **WebRTC 认证走 query 参数** — JWT 通过 `?jwt=TOKEN` 传递，避免跨域 CORS preflight 拦截 `Authorization` header

## 许可证

[MIT](LICENSE)
