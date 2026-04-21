# PawStream

Self-hosted live video streaming system for home pet monitoring.

Multiple edge devices (Raspberry Pi / Orange Pi + USB cameras) stream video to a central [MediaMTX](https://github.com/bluenviron/mediamtx) server. Users watch real-time footage through a mobile-first web UI via WebRTC.

## Architecture

```
[Edge Device + Camera] ──RTSP/RTMP──▶ MediaMTX ──WebRTC/HLS──▶ Mobile Browser
                                         ▲
                                    API Server
                                  (auth · devices)
```

- **Edge Client** (Go) — captures video, pushes to MediaMTX with device credentials
- **API Server** (Go) — control plane: user auth, device management, MediaMTX auth callback
- **MediaMTX** — handles all media transport and protocol translation
- **Web UI** (Vue 3) — mobile-first stream viewer with WebRTC playback

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Edge Client | Go 1.23, FFmpeg / GStreamer |
| API Server | Go 1.24, Fiber v2, SQLite (pure Go) |
| Media Server | MediaMTX (Docker) |
| Frontend | Vue 3, Vite 7, Vant, Pinia |
| Playback | WebRTC (WHEP), HLS fallback |

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js 18+
- Docker
- FFmpeg

### 1. Start MediaMTX

```bash
# Replace <YOUR_IP> with your LAN IP
docker run --rm -d \
  -e MTX_AUTHMETHOD=http \
  -e MTX_AUTHHTTPADDRESS=http://<YOUR_IP>:3000/mediamtx/auth \
  -e MTX_WEBRTCADDITIONALHOSTS=<YOUR_IP> \
  -p 8554:8554 -p 1935:1935 -p 8888:8888 -p 8889:8889 \
  -p 8890:8890/udp -p 8189:8189/udp \
  --name mediamtx bluenviron/mediamtx:latest-ffmpeg
```

### 2. Start API Server

```bash
cd server/api
cp config.example.yaml config.yaml  # Edit mediamtx URLs with your IP
go build -o api ./cmd/api && ./api
```

### 3. Create User & Device

```bash
# Register
curl -X POST http://localhost:3000/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Login
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

# Create device (save the returned id and secret)
curl -X POST http://localhost:3000/api/devices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"my-cam","location":"living-room"}'
```

### 4. Push a Test Stream

```bash
# Replace DEVICE_ID, DEVICE_SECRET, PUBLISH_PATH with values from previous step
ffmpeg -re -f lavfi -i "testsrc=size=1280x720:rate=30" \
  -c:v libx264 -preset ultrafast -tune zerolatency -b:v 2000k -g 30 -keyint_min 30 \
  -f rtsp -rtsp_transport tcp \
  "rtsp://DEVICE_ID:DEVICE_SECRET@localhost:8554/PUBLISH_PATH"
```

### 5. Start Web UI

```bash
cd web
npm install && npm run dev
```

Open http://localhost:5173, log in, and play your stream.

## Project Structure

```
client/edge/       Go edge client (video capture + push)
server/api/        Go API server (auth, device management)
server/mediamtx/   MediaMTX configuration
web/               Vue 3 frontend
docs/              Architecture & design docs
openspec/          Change proposals & specs
```

## Key Design Decisions

- **Media and control separation** — API server never touches video data; MediaMTX handles all media transport
- **No real hardware required for development** — all components work with FFmpeg test patterns and simulated sources
- **Edge devices are stateless** — configuration-driven, disposable
- **Pure Go SQLite** — `modernc.org/sqlite`, no CGO dependency
- **WebRTC auth via query parameter** — JWT passed as `?jwt=TOKEN` on WHEP URLs to avoid CORS preflight issues with `Authorization` header

## License

[MIT](LICENSE)
