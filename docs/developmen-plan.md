# PawStream Development Plan

## Phase 0 – Foundation (Day 1)

Goal: Establish project skeleton and development baseline.

- Create monorepo structure
- Initialize OpenSpec
- Define coding conventions
- Add README and basic documentation

Deliverables:
- project.md
- AGENTS.md
- development-plan.md
- Repository structure committed

---

## Phase 1 – Media Pipeline Validation (Day 2–3)

Goal: Verify end-to-end video streaming.

Tasks:
- Deploy MediaMTX on server
- Configure RTSP input path
- Configure WebRTC output
- Use GStreamer CLI to push USB camera stream
- Verify playback in browser

Deliverables:
- Working RTSP → WebRTC flow
- MediaMTX config documented

---

## Phase 2 – Edge Client MVP (Day 4–6)

Goal: Basic, reliable device-side streaming agent.

Tasks:
- Define edge client configuration format (YAML/TOML)
- Implement Go-based process wrapper OR GStreamer integration
- Support:
  - Device ID
  - Stream name
  - Server endpoint
- Add basic logging and restart strategy

Deliverables:
- `pawstream-agent` MVP
- One-command startup on device

---

## Phase 3 – Server API (Day 7–9)

Goal: Control plane for streams and users.

Tasks:
- Implement API server
- Stream registry (logical, not media-level)
- Authentication (simple token-based)
- Health and status endpoints

Deliverables:
- API server running
- Streams visible via API

---

## Phase 4 – Web UI MVP (Day 10–12)

Goal: Mobile-friendly live view.

Tasks:
- Vue 3 + Vite + Vant setup
- Login page
- Stream list page
- Live view page (WebRTC player)
- Basic error handling

Deliverables:
- Mobile browser live view works
- Authentication enforced

---

## Phase 5 – Hardening & Polish (Day 13–15)

Goal: Make the system stable and usable.

Tasks:
- Stream auth protection
- Better logging
- Configuration cleanup
- Deployment docs
- Failure recovery tests

Deliverables:
- End-to-end stable demo
- Documentation sufficient for redeployment

---

## Future Phases (Optional)

- Recording and playback
- AI pet behavior detection
- Notifications
- Native app packaging
