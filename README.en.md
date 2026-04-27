# PawStream

Self-hosted live video streaming for home pet monitoring. Edge devices push streams to [MediaMTX](https://github.com/bluenviron/mediamtx), users watch via WebRTC in a mobile-first web UI.

```
[Edge Device + Camera] ──RTSP──▶ MediaMTX ──WebRTC──▶ Mobile Browser
                                    ▲
                               API Server
```

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Edge Client | Go, FFmpeg / GStreamer |
| API Server | Go, Fiber, SQLite (pure Go) |
| Media Server | MediaMTX (Docker) |
| Frontend | Vue 3, Vite, Vant |

## Project Structure

```
client/edge/       Edge client (video capture + push)
server/api/        API server (auth, device management)
server/mediamtx/   MediaMTX configuration
web/               Vue 3 frontend
```

## Quick Start

```bash
# 1. Start MediaMTX (replace <IP> with your LAN IP)
docker run --rm -d \
  -e MTX_AUTHMETHOD=http \
  -e MTX_AUTHHTTPADDRESS=http://<IP>:3000/mediamtx/auth \
  -e MTX_WEBRTCADDITIONALHOSTS=<IP> \
  -p 8554:8554 -p 1935:1935 -p 8888:8888 -p 8889:8889 \
  -p 8890:8890/udp -p 8189:8189/udp \
  --name mediamtx bluenviron/mediamtx:latest-ffmpeg

# 2. Start API server
cd server/api && cp config.example.yaml config.yaml
go build -o api ./cmd/api && ./api

# 3. Start frontend
cd web && npm install && npm run dev
```

See [CLAUDE.md](CLAUDE.md) for detailed configuration.

## License

[MIT](LICENSE)
