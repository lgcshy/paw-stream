# Implementation Tasks: Web UI API Integration and WebRTC Player

## 1. Project Setup and Dependencies
- [ ] 1.1 Install Pinia (if not already installed)
- [ ] 1.2 Install any required WebRTC libraries
- [ ] 1.3 Create `.env.development` with API and MediaMTX URLs
- [ ] 1.4 Create `.env.production` template
- [ ] 1.5 Update `.gitignore` to exclude `.env.local`

## 2. TypeScript Type Definitions
- [ ] 2.1 Create `src/types/api.ts` for API request/response types
  - LoginRequest, LoginResponse
  - RegisterRequest, RegisterResponse
  - UserInfo
  - DeviceInfo, PathInfo
  - ApiError
- [ ] 2.2 Update `src/types/stream.ts` with real data structure
- [ ] 2.3 Add type definitions for WebRTC player

## 3. API Client Layer
- [ ] 3.1 Create `src/api/client.ts` - Base HTTP client
  - Configure base URL from env
  - Auto-attach JWT token
  - Handle 401 errors (auto logout)
  - Unified error handling
  - TypeScript request/response typing
- [ ] 3.2 Create `src/api/auth.ts` - Authentication API
  - login(username, password)
  - register(username, password, nickname) - optional
  - getCurrentUser()
- [ ] 3.3 Create `src/api/device.ts` - Device API (optional for Phase 4)
  - listDevices()
  - getDevice(id)
- [ ] 3.4 Create `src/api/path.ts` - Path API
  - listPaths()

## 4. State Management (Pinia)
- [ ] 4.1 Setup Pinia in `src/main.ts` (if not already)
- [ ] 4.2 Create `src/stores/auth.ts` - Auth store
  - State: token, user, isAuthenticated
  - Actions: login, logout, loadToken, checkAuth
  - Persist token to localStorage
- [ ] 4.3 Create `src/stores/device.ts` - Device store
  - State: paths, loading, error
  - Actions: fetchPaths, refreshPaths
  - Getters: enabledPaths

## 5. Router Guards
- [ ] 5.1 Update `src/router/index.ts`
  - Add global beforeEach guard
  - Check auth status from auth store
  - Redirect unauthenticated users to /login
  - Redirect authenticated users away from /login
  - Preserve original URL for post-login redirect

## 6. Login View Implementation
- [ ] 6.1 Update `src/views/LoginView.vue`
  - Import and use auth store
  - Call auth.login() on form submit
  - Handle loading state
  - Handle API errors (invalid credentials, network error)
  - Show appropriate toast messages
  - Redirect to /streams on success

## 7. Stream List View Implementation
- [ ] 7.1 Update `src/views/StreamListView.vue`
  - Import and use device store
  - Fetch paths on component mount
  - Display real path data (name, location, publish_path)
  - Handle loading state
  - Handle empty state
  - Handle API errors
  - Navigate to player with correct path parameter

## 8. WebRTC Player Integration
- [ ] 8.1 Research MediaMTX WebRTC client integration
  - Check MediaMTX documentation
  - Determine if library needed or native WebRTC API sufficient
- [ ] 8.2 Create `src/utils/webrtc.ts` - WebRTC helper
  - createWebRTCConnection(path, token)
  - Handle ICE candidates
  - Handle connection state changes
  - Handle errors and reconnection
- [ ] 8.3 Update `src/views/StreamPlayerView.vue`
  - Get stream path from route params
  - Get JWT token from auth store
  - Initialize WebRTC connection on mount
  - Display video element
  - Handle connection states (connecting, connected, failed)
  - Implement reconnect logic
  - Cleanup on unmount
  - Add basic controls (play/pause, reconnect)

## 9. Environment Configuration
- [ ] 9.1 Define environment variables
  - VITE_API_BASE_URL (e.g., http://localhost:3000)
  - VITE_MEDIAMTX_WEBRTC_URL (e.g., http://localhost:8889)
- [ ] 9.2 Create `.env.development`
- [ ] 9.3 Create `.env.production` template
- [ ] 9.4 Update code to use import.meta.env

## 10. Error Handling and UX
- [ ] 10.1 Implement global error handler
- [ ] 10.2 Add loading indicators
  - Login button loading state
  - Stream list loading skeleton
  - Player loading state
- [ ] 10.3 Add error toasts
  - Network errors
  - API errors
  - WebRTC connection errors
- [ ] 10.4 Add retry mechanisms
  - API request retry
  - WebRTC reconnection

## 11. Token Management
- [ ] 11.1 Implement token persistence
  - Save to localStorage on login
  - Load from localStorage on app init
- [ ] 11.2 Implement token expiry handling
  - Detect 401 errors
  - Auto logout on token expiry
  - Clear localStorage
- [ ] 11.3 Implement logout functionality
  - Add logout button/menu item
  - Clear auth store
  - Clear localStorage
  - Redirect to login

## 12. UI/UX Polish
- [ ] 12.1 Update Layout component
  - Add user info display
  - Add logout button
- [ ] 12.2 Improve stream list UI
  - Add refresh button
  - Add pull-to-refresh (optional)
  - Better empty state
- [ ] 12.3 Improve player UI
  - Fullscreen support (optional)
  - Better error messages
  - Connection status indicator

## 13. Testing and Validation
- [ ] 13.1 Test login flow
  - Valid credentials
  - Invalid credentials
  - Network error
- [ ] 13.2 Test stream list
  - Load paths successfully
  - Empty state
  - Navigate to player
- [ ] 13.3 Test WebRTC player
  - Connect and play stream
  - Handle connection errors
  - Reconnect functionality
- [ ] 13.4 Test authentication flow
  - Login persistence across refresh
  - Route guards
  - Auto logout on token expiry
- [ ] 13.5 Test on mobile devices
  - Responsive layout
  - Touch interactions
  - Video playback

## 14. Documentation
- [ ] 14.1 Update `web/README.md`
  - Environment setup instructions
  - API configuration
  - MediaMTX configuration
- [ ] 14.2 Document WebRTC integration
  - Connection flow
  - Troubleshooting
- [ ] 14.3 Add inline code comments
  - Complex WebRTC logic
  - API error handling

## 15. Build and Deployment
- [ ] 15.1 Test production build
  - npm run build
  - Verify env variables
  - Test built app
- [ ] 15.2 Update deployment docs
  - Environment variable configuration
  - CORS configuration notes
  - MediaMTX URL configuration

## 16. Integration Testing
- [ ] 16.1 End-to-end test
  - Start API server
  - Start MediaMTX
  - Start web UI
  - Complete user journey: login → view streams → play video
- [ ] 16.2 Test with real device
  - Device publishing stream
  - User viewing stream
  - Authentication working end-to-end
