# Project Context

## Purpose

PawStream is a self-hosted live video streaming system for home pet monitoring.

The goal of the project is to enable multiple edge devices (Raspberry Pi / Orange Pi + USB cameras) to stream live video to a central server, where users can securely view real-time footage through a mobile-friendly Web UI.

Key goals:
- Low-latency live streaming
- Simple and reliable architecture
- Full control over data and privacy
- See-first, scale-later development suitable for a single developer

Non-goals:
- Commercial SaaS features
- Native mobile applications (for now)
- Cloud-managed or proprietary camera ecosystems

---

## Tech Stack

### Edge (Device / Client)
- Go (planned primary language)
- GStreamer / FFmpeg (development and simulation)
- Linux (ARM / x86)

### Server
- MediaMTX (media server, external dependency)
- Go (API server for control plane)
- RTSP / SRT (video ingest protocols)
- WebRTC / HLS (playback protocols)
- Token-based authentication
- SQLite / Postgres (database options for metadata)

### Frontend
- Vue 3
- Vite 7
- Vant (mobile-first UI library)
- WebRTC (browser-native playback)
- Mobile-first responsive design

### Repository Structure
- Monorepo layout
- Separate domains: `edge/`, `server/`, `web/`
- Documentation in `docs/`
- OpenSpec workflow for change proposals and specifications

---

## Project Conventions

### Code Style

- Prefer clarity over brevity
- Explicit configuration over implicit behavior
- Avoid revealing “magic” abstractions
- Follow standard formatting tools:
  - Go: `gofmt`
  - Frontend: Prettier defaults
- Naming:
  - Use descriptive names over abbreviations
  - Device- and stream-related concepts must be explicit

---

### Architecture Patterns

- Edge devices are stateless and disposable
- Server acts as a control plane, not a media processor
- MediaMTX is treated as an external dependency
- Media and control logic are strictly separated
- Video input sources are abstracted (camera, file, test source)

High-level pattern:
- RTSP/SRT in → MediaMTX → WebRTC/HLS out
- API server manages authentication, device metadata, and UI access

---

### Testing Strategy

Testing is a first-class concern during development.

Development is expected to occur **without physical camera hardware**.

Testing layers:
1. **Simulated video sources**
   - FFmpeg streaming from video files
   - GStreamer test sources
2. **Virtual cameras**
   - v4l2loopback for /dev/video* simulation
3. **Real hardware**
   - Used only for performance and stability validation

All components must support simulated inputs.
No component may assume the existence of real USB cameras.

Manual integration testing is acceptable for early stages.
Automated tests are added only when they provide clear value.

---

### Git Workflow

- Single main branch (`main`)
- Small, incremental commits
- Commit messages should describe intent, not implementation details
- Prefer working end-to-end features over partial abstractions
- Use OpenSpec workflow for significant changes:
  - Create change proposals in `openspec/changes/`
  - Validate with `openspec validate --strict` before implementation
  - Archive to `openspec/changes/archive/` after deployment

---

### Development Phases

The project follows a phased development approach:

1. **Phase 0**: Foundation and project setup
2. **Phase 1**: Media pipeline validation (RTSP → MediaMTX → WebRTC)
3. **Phase 2**: Edge client MVP (Go-based streaming agent)
4. **Phase 3**: Server API (control plane, stream registry, auth)
5. **Phase 4**: Web UI MVP (Vue 3 mobile-friendly interface)
6. **Phase 5**: Hardening and production readiness

Future phases may include recording, AI detection, and notifications.

---

## Domain Context

### Core Concepts

- **Device**: Logical streaming endpoint (not tied to physical hardware)
  - Can be a Raspberry Pi, Orange Pi, or development machine
  - Identified by unique device ID
  - May support multiple camera/stream sources

- **Stream**: Logical video feed from a device
  - Has unique stream name/path
  - Maps to MediaMTX path (e.g., `/stream/device-id/camera-1`)
  - Independent of transport protocol

- **Video Source**: Input to a stream
  - Physical USB camera (production)
  - Video file loop (development/testing)
  - Synthetic test pattern (GStreamer test source)
  - Virtual camera (v4l2loopback)

### Component Boundaries

- **MediaMTX**: Handles all media transport and protocol translation
  - RTSP/SRT ingestion
  - WebRTC/HLS distribution
  - No custom code required

- **API Server**: Control plane only
  - Stream registry and metadata
  - Authentication and authorization
  - Health and status monitoring
  - **Does NOT** manipulate video data

- **Edge Agent**: Thin streaming client
  - Reads camera configuration
  - Manages streaming process (GStreamer/FFmpeg)
  - Reports health to API server
  - Stateless and disposable

### Terminology Standards

- Use "device" (not "camera", "board", "client")
- Use "stream" (not "feed", "channel", "source")
- Use "edge agent" or "agent" (not "client", "daemon")
- Use "API server" (not "backend", "control server")
- Use "web UI" or "frontend" (not "dashboard", "portal")

---

## Important Constraints

- Single developer project
- Limited hardware availability during development
- Must be deployable on low-cost servers (e.g., affordable VPS)
- Must remain understandable after long periods of inactivity
- Security should be simple but explicit (no hidden defaults)
- Performance requirements: support 4-8 concurrent camera streams
- Low-latency requirement: suitable for real-time pet monitoring
- Self-hosted: complete data privacy and control

---

## Security Considerations

- **Authentication**: Simple token-based auth (Phase 3+)
- **Stream access control**: Protected WebRTC/HLS endpoints
- **Data privacy**: All video data remains on user's infrastructure
- **No cloud dependencies**: Zero third-party video processing
- **Explicit configuration**: No auto-discovery or magic defaults that could expose streams
- **HTTPS/TLS**: Required for production deployment (WebRTC requirement)

---

## External Dependencies

### Core Dependencies
- **MediaMTX**: RTSP/SRT/WebRTC media server (runs as separate service)
- **GStreamer / FFmpeg**: Media pipeline tools for development and testing
- **WebRTC**: Browser-native streaming (no additional libraries required)

### System Dependencies
- **Linux video subsystem**: v4l2 for camera access
- **v4l2loopback**: Virtual camera devices for testing without hardware
- **Docker**: Optional, for MediaMTX deployment

### Development Tools
- **Go**: `gofmt` for code formatting
- **Node.js/npm**: Frontend development (Vite, Vue 3)
- **Prettier**: Frontend code formatting
- **OpenSpec CLI**: Change proposal and specification management

---

## Development Environment

### Target Platforms
- **Edge devices**: Raspberry Pi / Orange Pi (ARM), also x86 for development
- **Server**: Linux server (low-cost VPS or local)
- **Client**: Mobile browsers (iOS Safari, Android Chrome)

### Typical Development Setup
1. Local MediaMTX instance (Docker or native)
2. Simulated video sources (FFmpeg with test video files, v4l2loopback)
3. API server running locally
4. Frontend dev server (Vite HMR)
5. No physical cameras required for most development work

### Key Scripts & Commands
- `server/mediamtx/run_mtx_docker.sh`: Start MediaMTX in Docker
- `server/mediamtx/publish.sh`: Test RTSP publishing with FFmpeg
- OpenSpec commands for change management (see `openspec/AGENTS.md`)
