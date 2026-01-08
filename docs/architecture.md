# Paw Stream Architecture

[树莓派/香橙派 + USB Camera] x 4~8
        |
   (RTSP / SRT / WebRTC 推流)
        |
      MediaMTX
        |
 ┌───────────────┐
 │  Web Server   │  (Vite + Vue3 + Vant)
 │  Auth / API   │
 └───────────────┘
        |
   Mobile Browser
