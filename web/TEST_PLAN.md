# Phase 4 Test Plan - Web UI Integration

## Test Environment Setup

### 1. Start API Server
```bash
cd server/api
./bin/api
# API should be running on http://localhost:3000
```

### 2. Start MediaMTX (if testing actual streaming)
```bash
cd server/mediamtx
docker-compose up -d
# MediaMTX should be running on http://localhost:8889
```

### 3. Start Web UI
```bash
cd web
npm run dev
# Web UI should be running on http://localhost:5173
```

## Test Scenarios

### ✅ Scenario 1: User Authentication
1. Open browser to http://localhost:5173
2. Should redirect to `/login`
3. Enter test credentials:
   - Username: `testuser` (or create via API)
   - Password: `test123`
4. Click "登录"
5. Should see success toast
6. Should redirect to `/streams`

**Expected**:
- Login API called: `POST http://localhost:3000/api/login`
- JWT token saved to localStorage
- User redirected to streams page

### ✅ Scenario 2: Token Persistence
1. After successful login, refresh the page (F5)
2. Should remain logged in
3. Should not redirect to login page

**Expected**:
- Token loaded from localStorage
- `/api/me` called to validate token
- User stays authenticated

### ✅ Scenario 3: Stream List Display
1. After login, on `/streams` page
2. Should see loading indicator
3. Should fetch and display devices

**Expected**:
- API called: `GET http://localhost:3000/api/paths`
- Devices displayed with name, location
- "在线" tag shown for enabled devices

### ✅ Scenario 4: Empty Stream List
1. If user has no devices
2. Should show empty state message
3. Message: "暂无设备,请先在 API 中创建设备"

**Expected**:
- Empty state displayed
- No error messages

### ✅ Scenario 5: Navigate to Player
1. Click on a stream in the list
2. Should navigate to `/stream/:path`
3. Player page should load

**Expected**:
- URL encodes publish_path correctly
- Player component mounts

### ✅ Scenario 6: WebRTC Stream Playback (Requires MediaMTX + Active Stream)
1. Navigate to a stream with active publishing
2. Should see "连接中..." loading indicator
3. WebRTC connection should establish
4. Video should start playing

**Expected**:
- WebRTC WHEP request: `POST http://localhost:8889/{path}/whep`
- Authorization header includes JWT token
- Connection state changes: new → connecting → connected
- Video element displays stream

### ✅ Scenario 7: Stream Connection Error
1. Navigate to a non-existent or inactive stream
2. Should see error message
3. "重试" and "返回" buttons available

**Expected**:
- Error overlay displayed
- User can retry or go back
- No app crash

### ✅ Scenario 8: Logout (Manual - via API or Token Expiry)
1. Clear localStorage manually or wait for token expiry
2. Try to navigate to `/streams`
3. Should redirect to `/login`

**Expected**:
- 401 from API triggers logout
- Redirect to login page
- Token cleared from storage

### ✅ Scenario 9: Route Protection
1. Without logging in, try to access `/streams` directly
2. Should redirect to `/login`
3. After login, should redirect back to original URL

**Expected**:
- Route guard blocks access
- Redirect with query param: `/login?redirect=/streams`
- Post-login redirect to original destination

### ✅ Scenario 10: Pull to Refresh
1. On streams list page
2. Pull down to refresh
3. Should reload device list

**Expected**:
- Loading indicator shown
- `/api/paths` called again
- List updates
- Success toast shown

## Manual Testing Checklist

- [ ] User can login with valid credentials
- [ ] Invalid credentials show error
- [ ] Token persists across page refresh
- [ ] Unauthenticated users redirected to login
- [ ] Stream list displays real API data
- [ ] Empty state shown when no devices
- [ ] Click stream navigates to player
- [ ] WebRTC player attempts connection
- [ ] Error handling works (network errors, invalid streams)
- [ ] Logout clears state and redirects
- [ ] Mobile responsive design works
- [ ] All TypeScript types are correct
- [ ] Production build succeeds
- [ ] No console errors in normal flow

## API Integration Test

Create a test user and device via API:

```bash
# 1. Register user
curl -X POST http://localhost:3000/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123","nickname":"Test User"}'

# 2. Login
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}' | jq -r '.token')

# 3. Create device
curl -X POST http://localhost:3000/api/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"测试摄像头","location":"测试位置"}'

# 4. List paths
curl http://localhost:3000/api/paths \
  -H "Authorization: Bearer $TOKEN"
```

Then test in Web UI:
1. Login with `testuser` / `test123`
2. Should see "测试摄像头" in stream list
3. Click to open player (will show error if no active stream)

## Known Limitations (Phase 4)

- No user registration UI (use API directly)
- No device management UI (use API directly)
- No stream recording/playback
- Connection state not real-time (requires manual refresh)
- No advanced player controls (volume, fullscreen via native controls only)

## Success Criteria

All checkboxes above should be checked ✅

Phase 4 is complete when:
- Users can login via Web UI
- Token persistence works
- Stream list shows real API data
- WebRTC player connects to MediaMTX (when stream available)
- All error cases handled gracefully
- Mobile responsive
- TypeScript strict mode passes
- Production build succeeds
