# PawStream

[English](README.en.md) | 中文

自托管的家庭宠物实时视频监控系统。边缘设备推流至 [MediaMTX](https://github.com/bluenviron/mediamtx)，用户通过 Web UI 以 WebRTC 实时观看。

```
[边缘设备 + 摄像头] ──RTSP──▶ MediaMTX ──WebRTC──▶ 手机浏览器
                                  ▲
                             API 服务器
```

## 技术栈

| 组件 | 技术 |
|------|------|
| 边缘客户端 | Go, FFmpeg / GStreamer |
| API 服务器 | Go, Fiber, SQLite (纯 Go) |
| 媒体服务器 | MediaMTX (Docker) |
| 前端 | Vue 3, Vite, Vant |

## 项目结构

```
client/edge/       边缘客户端（视频采集 + 推流）
server/api/        API 服务器（认证、设备管理）
server/mediamtx/   MediaMTX 配置
web/               Vue 3 前端
```

## 快速开始

```bash
# 1. 启动 MediaMTX（替换 <IP> 为局域网 IP）
docker run --rm -d \
  -e MTX_AUTHMETHOD=http \
  -e MTX_AUTHHTTPADDRESS=http://<IP>:3000/mediamtx/auth \
  -e MTX_WEBRTCADDITIONALHOSTS=<IP> \
  -p 8554:8554 -p 1935:1935 -p 8888:8888 -p 8889:8889 \
  -p 8890:8890/udp -p 8189:8189/udp \
  --name mediamtx bluenviron/mediamtx:latest-ffmpeg

# 2. 启动 API 服务器
cd server/api && cp config.example.yaml config.yaml
go build -o api ./cmd/api && ./api

# 3. 启动前端
cd web && npm install && npm run dev
```

详细配置说明见 [CLAUDE.md](CLAUDE.md)。

## 许可证

[MIT](LICENSE)
