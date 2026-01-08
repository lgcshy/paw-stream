# Web UI Specification - Phase 4 Implementation

## ADDED Requirements

### Requirement: API Client Layer
The web UI SHALL provide a type-safe HTTP client for communicating with the PawStream API server.

#### Scenario: API request with authentication
- **WHEN** an authenticated API request is made
- **THEN** the JWT token is automatically included in the Authorization header
- **AND** the request is sent to the configured API base URL
- **AND** TypeScript types ensure type safety for request and response

#### Scenario: API error handling
- **WHEN** an API request fails with 401 Unauthorized
- **THEN** the user is automatically logged out
- **AND** redirected to the login page
- **AND** a user-friendly error message is displayed

#### Scenario: Network error handling
- **WHEN** an API request fails due to network issues
- **THEN** a user-friendly error message is displayed
- **AND** the error is logged for debugging

---

### Requirement: Authentication State Management
The web UI SHALL manage user authentication state using Pinia store with persistent token storage.

#### Scenario: User login
- **WHEN** user submits valid credentials on the login page
- **THEN** the API `/api/login` endpoint is called
- **AND** the returned JWT token is stored in localStorage
- **AND** the user information is stored in Pinia auth store
- **AND** the user is redirected to the stream list page

#### Scenario: Token persistence
- **WHEN** the user refreshes the page
- **THEN** the JWT token is loaded from localStorage
- **AND** the user remains logged in
- **AND** the auth store is rehydrated with user information

#### Scenario: User logout
- **WHEN** user logs out or token expires
- **THEN** the JWT token is removed from localStorage
- **AND** the auth store is cleared
- **AND** the user is redirected to the login page

---

### Requirement: Device Data Management
The web UI SHALL fetch and manage device and stream path data from the API.

#### Scenario: Load stream paths
- **WHEN** an authenticated user navigates to the stream list page
- **THEN** the `/api/paths` endpoint is called with JWT token
- **AND** the returned paths are stored in device store
- **AND** the stream list is displayed with device names and locations

#### Scenario: Empty stream list
- **WHEN** the user has no devices
- **THEN** an empty state message is displayed
- **AND** the message suggests creating a device

---

### Requirement: Route Protection
The web UI SHALL implement route guards to protect authenticated pages.

#### Scenario: Access protected route without authentication
- **WHEN** an unauthenticated user tries to access `/streams` or `/stream/:id`
- **THEN** the user is redirected to `/login`
- **AND** the original URL is preserved for redirect after login

#### Scenario: Access login page when authenticated
- **WHEN** an authenticated user navigates to `/login`
- **THEN** the user is redirected to `/streams`

---

### Requirement: WebRTC Stream Player
The web UI SHALL integrate WebRTC playback for live camera streams from MediaMTX.

#### Scenario: Start stream playback
- **WHEN** user clicks on a stream in the stream list
- **THEN** the player page opens with the stream ID
- **AND** a WebRTC connection is established to MediaMTX
- **AND** the JWT token is included in the WebRTC connection for authentication
- **AND** the video stream is displayed in the player

#### Scenario: Stream connection failure
- **WHEN** the WebRTC connection fails to establish
- **THEN** an error message is displayed to the user
- **AND** a retry button is provided
- **AND** the error is logged for debugging

#### Scenario: Stream disconnection
- **WHEN** an active stream disconnects
- **THEN** the user is notified of the disconnection
- **AND** an automatic reconnection is attempted
- **AND** a manual reconnect button is available

---

### Requirement: Environment Configuration
The web UI SHALL support environment-specific configuration for API and MediaMTX URLs.

#### Scenario: Development environment
- **WHEN** running in development mode (`npm run dev`)
- **THEN** API requests are sent to `VITE_API_BASE_URL` from `.env.development`
- **AND** WebRTC connections use `VITE_MEDIAMTX_WEBRTC_URL` from `.env.development`

#### Scenario: Production build
- **WHEN** building for production (`npm run build`)
- **THEN** API and MediaMTX URLs are read from `.env.production`
- **AND** the configuration is embedded in the build

---

## ADDED Requirements

### Requirement: Login View Implementation
The login view SHALL authenticate users against the real API server instead of using placeholder logic.

#### Scenario: Successful login
- **WHEN** user enters valid credentials and submits the form
- **THEN** the `/api/login` API is called
- **AND** the JWT token is saved to localStorage
- **AND** the user is redirected to `/streams`
- **AND** a success toast is displayed

#### Scenario: Invalid credentials
- **WHEN** user enters invalid credentials
- **THEN** a 401 error is returned from the API
- **AND** an error toast is displayed: "用户名或密码错误"
- **AND** the user remains on the login page

#### Scenario: Network error during login
- **WHEN** the login API request fails due to network issues
- **THEN** an error toast is displayed: "网络错误,请稍后重试"
- **AND** the user can retry the login

---

### Requirement: Stream List View Implementation
The stream list view SHALL display real device data from the API instead of mock data.

#### Scenario: Display user's devices
- **WHEN** an authenticated user navigates to the stream list page
- **THEN** the `/api/paths` API is called
- **AND** each device is displayed with name, location, and publish path
- **AND** only enabled devices are shown

#### Scenario: Navigate to stream player
- **WHEN** user clicks on a stream in the list
- **THEN** the user is navigated to `/stream/:path` with the publish path
- **AND** the stream player page opens

#### Scenario: No devices available
- **WHEN** the user has no devices
- **THEN** an empty state is displayed
- **AND** the message reads: "暂无设备,请先创建设备"

---

### Requirement: Stream Player View Implementation
The stream player view SHALL play live WebRTC streams from MediaMTX instead of showing a placeholder.

#### Scenario: Load and play stream
- **WHEN** the player page loads with a stream path
- **THEN** a WebRTC connection is established to MediaMTX
- **AND** the JWT token is sent for authorization
- **AND** the video stream is rendered in a video element
- **AND** the stream plays automatically

#### Scenario: Stream not found
- **WHEN** the stream path does not exist or is not accessible
- **THEN** an error message is displayed: "视频流不存在或无权访问"
- **AND** a back button is provided to return to the stream list

#### Scenario: Player controls
- **WHEN** the stream is playing
- **THEN** basic controls are available (play/pause, reconnect)
- **AND** the stream path is displayed for reference
