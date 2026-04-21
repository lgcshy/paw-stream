# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PawStream is a self-hosted live video streaming system for home pet monitoring. Edge devices (Raspberry Pi/Orange Pi + USB cameras) stream video to a central MediaMTX server, with users viewing streams through a mobile-first Vue 3 web UI. The system is a single-developer project emphasizing simplicity, privacy, and self-hosting.

## Repository Structure (Monorepo)

- **`client/edge/`** — Go edge client that captures video and pushes to MediaMTX (module: `github.com/lgc/pawstream/edge-client`, Go 1.23)
- **`server/api/`** — Go API server (control plane only, no media processing) (module: `github.com/lgc/pawstream/api`, Go 1.24)
- **`server/mediamtx/`** — MediaMTX configuration files
- **`web/`** — Vue 3 + Vite 7 + Vant mobile-first frontend
- **`openspec/`** — Change proposals and specs (see `openspec/AGENTS.md` for workflow)
- **`docs/`** — Architecture docs and design decisions

## Build & Development Commands

### Edge Client (`client/edge/`)
```bash
cd client/edge
make build                    # Build binary to build/edge-client
make run                      # Build and run
make test                     # Run all tests
go test -v ./internal/...     # Run tests for specific package
make lint                     # golangci-lint (if installed)
make fmt                      # gofmt
make build-all                # Cross-compile for all platforms
```

### API Server (`server/api/`)
```bash
cd server/api
go build -o api ./cmd/api     # Build
go run ./cmd/api              # Run (reads config.yaml or PAWSTREAM_* env vars)
go test -v ./...              # Run all tests
go test -v ./internal/domain/device/...  # Run specific package tests
```

Config is loaded via Viper: `config.yaml` in current directory, or `PAWSTREAM_*` environment variables (e.g., `PAWSTREAM_SERVER_PORT`, `PAWSTREAM_JWT_SECRET`, `PAWSTREAM_DB_PATH`).

### Web Frontend (`web/`)
```bash
cd web
npm install
npm run dev                   # Vite dev server with HMR
npm run build                 # Type-check + production build
npm run type-check            # vue-tsc --noEmit
npm run format                # Prettier
```

### Docker (Full Stack)
```bash
cd server/api/deployments
docker-compose up             # Starts API server + MediaMTX
```

## Architecture

### Data Flow
```
Edge Device → (RTSP/RTMP) → MediaMTX → (WebRTC/HLS) → Web UI
                                ↑
                    API Server (auth callback)
```

### Auth Model
- **Publish (device→MediaMTX)**: Device authenticates with device ID + secret via RTSP/RTMP credentials. MediaMTX calls `POST /mediamtx/auth` on the API server for every publish/read action.
- **Read (user→MediaMTX)**: Users authenticate with JWT tokens passed via WHEP URL query parameter (`?jwt=TOKEN`). The `Authorization: Bearer` header does NOT work for cross-origin WebRTC/WHEP requests due to CORS preflight restrictions. The API server extracts the JWT from the MediaMTX auth callback's `query` field.
- **API access**: JWT-based auth via `Authorization` header. Middleware at `server/api/internal/transport/http/middleware/auth.go`.

### Edge Client Internals
- **Input sources** (`client/edge/internal/capture/`): Abstracted via `InputSource` interface — v4l2, RTSP, file, test pattern. All must support simulated inputs (no real camera required).
- **Stream engines** (`client/edge/internal/stream/`): Factory pattern selects FFmpeg or GStreamer engine. GStreamer is preferred but falls back to FFmpeg if unavailable. Hardware encoding auto-detected (NVENC, VAAPI, QSV, VideoToolbox).
- **Stream manager**: Handles lifecycle, reconnection logic, and error monitoring.
- **Config hot-reload**: File watcher triggers config reload via fsnotify.
- **Web UI**: Optional embedded web server for device setup and status monitoring (SSE-based).

### API Server Internals
- **Layered architecture**: `transport/http/handlers` → `domain/{user,device,acl}` services → `store/sqlite` repositories
- **Database**: SQLite via `modernc.org/sqlite` (pure Go, no CGO). Migrations in `server/api/migrations/`.
- **HTTP framework**: Fiber v2 with middleware stack (recovery, request ID, logger, CORS).

### Frontend
- Vue 3 + Pinia stores + Vue Router
- Vant component library (mobile-first)
- API client layer in `web/src/api/`

### Local Dev Stack
Start services in order (MediaMTX must use env-based config — the bundled `mediamtx.yml` has fields incompatible with latest image):
```bash
# 1. MediaMTX (Docker, env-based config, replace IP with your LAN IP)
docker run --rm -d \
  -e MTX_AUTHMETHOD=http \
  -e MTX_AUTHHTTPADDRESS=http://<YOUR_IP>:3000/mediamtx/auth \
  -e MTX_WEBRTCADDITIONALHOSTS=<YOUR_IP> \
  -p 8554:8554 -p 1935:1935 -p 8888:8888 -p 8889:8889 \
  -p 8890:8890/udp -p 8189:8189/udp \
  --name mediamtx bluenviron/mediamtx:latest-ffmpeg

# 2. API Server
cd server/api && go build -o api ./cmd/api && ./api

# 3. Test stream (FFmpeg, use -g 30 for 1s keyframe interval to avoid WebRTC black screen)
ffmpeg -re -f lavfi -i "testsrc=size=1280x720:rate=30" \
  -c:v libx264 -preset ultrafast -tune zerolatency -b:v 2000k -g 30 -keyint_min 30 \
  -f rtsp -rtsp_transport tcp "rtsp://DEVICE_ID:SECRET@localhost:8554/PUBLISH_PATH"

# 4. Web Frontend
cd web && npm run dev
```

## Key Conventions

- **Terminology**: Use "device" (not camera/client), "stream" (not feed/channel), "edge agent" (not daemon), "API server" (not backend)
- **Code style**: `gofmt` for Go, Prettier for frontend. Prefer clarity over brevity. Explicit configuration over implicit behavior.
- **No real hardware assumed**: All components must work with simulated video sources (test patterns, file loops, v4l2loopback)
- **Media and control separation**: API server never touches video data; MediaMTX handles all media transport
- **Edge devices are stateless**: Configuration-driven, disposable

## OpenSpec Workflow

For planning, proposals, or architecture changes, consult `openspec/AGENTS.md`. Change proposals go in `openspec/changes/` and must be validated with `openspec validate --strict` before implementation.
