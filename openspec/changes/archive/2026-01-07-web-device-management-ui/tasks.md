# Implementation Tasks: Device Management UI and User Registration

## 1. API Layer Enhancement
- [ ] 1.1 Update `src/api/device.ts`
  - Add createDevice(data)
  - Add updateDevice(id, data)
  - Add deleteDevice(id)
  - Add rotateSecret(id)
- [ ] 1.2 Update `src/api/auth.ts`
  - Add register(data)
- [ ] 1.3 Add TypeScript types for new requests/responses
  - CreateDeviceRequest, UpdateDeviceRequest
  - RotateSecretResponse

## 2. Store Enhancement
- [ ] 2.1 Update `src/stores/device.ts`
  - Add createDevice action
  - Add updateDevice action
  - Add deleteDevice action
  - Add rotateSecret action
  - Add getDeviceById getter
- [ ] 2.2 Update `src/stores/auth.ts`
  - Add register action

## 3. Reusable Components
- [ ] 3.1 Create `src/components/SecretDisplay.vue`
  - Display secret in monospace
  - Copy to clipboard button
  - One-time warning message
  - Success feedback on copy
- [ ] 3.2 Create `src/components/ConfirmDialog.vue`
  - Title, message props
  - Confirm/cancel buttons
  - Customizable button text
  - Emit confirm/cancel events
- [ ] 3.3 Update `src/components/Layout.vue`
  - Add bottom navigation bar
  - Icons for: streams, devices, profile
  - Highlight active section
  - Hide on sub-pages (use route meta)

## 4. User Registration
- [ ] 4.1 Create `src/views/RegisterView.vue`
  - Username, password, nickname fields
  - Password strength indicator
  - Form validation
  - Call auth.register()
  - Auto-login after registration
  - Navigate to /streams
- [ ] 4.2 Update `src/views/LoginView.vue`
  - Add "注册账号" link to /register

## 5. Device Management Pages
- [ ] 5.1 Create `src/views/DeviceListView.vue` (or rename/refactor StreamListView)
  - List all user devices (enabled + disabled)
  - Show device status, name, location
  - "新增设备" button
  - Click device → navigate to detail
  - Pull to refresh
- [ ] 5.2 Create `src/views/DeviceFormView.vue`
  - Form for create/edit device
  - Name, location fields
  - Submit → createDevice or updateDevice
  - Show SecretDisplay after creation
  - Validation
- [ ] 5.3 Create `src/views/DeviceDetailView.vue`
  - Display all device info
  - Enable/disable switch
  - Edit button → /devices/:id/edit
  - Rotate secret button
  - Delete button
  - Play stream button (if enabled)
- [ ] 5.4 Update `src/views/StreamListView.vue`
  - Focus on playback (enabled devices only)
  - Remove device management actions
  - Keep pull to refresh

## 6. User Profile
- [ ] 6.1 Create `src/views/ProfileView.vue`
  - Display user info (username, nickname, created_at)
  - Device count statistics
  - Logout button
  - Logout confirmation (optional)

## 7. Router Updates
- [ ] 7.1 Add new routes to `src/router/index.ts`
  - /register (no auth required)
  - /devices (auth required) - device list
  - /devices/new (auth required) - create device
  - /devices/:id (auth required) - device detail
  - /devices/:id/edit (auth required) - edit device
  - /profile (auth required) - user profile
- [ ] 7.2 Add route meta for bottom nav visibility
  - showBottomNav: true for main pages
  - showBottomNav: false for sub-pages

## 8. Device CRUD Implementation
- [ ] 8.1 Implement create device flow
  - DeviceFormView with empty form
  - Call device.createDevice()
  - Show SecretDisplay with result
  - Provide "创建另一个" and "返回列表" options
- [ ] 8.2 Implement edit device flow
  - DeviceFormView with pre-filled data
  - Call device.updateDevice()
  - Navigate back to detail on success
- [ ] 8.3 Implement delete device flow
  - Show ConfirmDialog
  - Call device.deleteDevice()
  - Navigate to device list on success
- [ ] 8.4 Implement rotate secret flow
  - Show ConfirmDialog with warning
  - Call device.rotateSecret()
  - Show SecretDisplay with new secret

## 9. User Experience Enhancements
- [ ] 9.1 Add loading states
  - Form submission loading
  - Device list loading
  - Profile loading
- [ ] 9.2 Add error handling
  - API errors with toast
  - Form validation errors
  - Network errors
- [ ] 9.3 Add success feedback
  - Toast on successful operations
  - Visual confirmation
- [ ] 9.4 Add empty states
  - No devices message
  - "创建第一个设备" CTA

## 10. Navigation Flow
- [ ] 10.1 Implement bottom navigation
  - Tab bar with icons
  - Active state highlighting
  - Route navigation on tap
- [ ] 10.2 Update Layout logic
  - Show/hide bottom nav based on route meta
  - Ensure proper spacing for content

## 11. Secret Management
- [ ] 11.1 Implement SecretDisplay component
  - Monospace display
  - Copy button with clipboard API
  - Success toast on copy
  - Warning message
- [ ] 11.2 Handle secret in create flow
  - Display immediately after creation
  - Warn user to save it
  - Don't show again after navigation
- [ ] 11.3 Handle secret in rotate flow
  - Similar to create flow
  - Warn about old secret invalidation

## 12. Confirmation Dialogs
- [ ] 12.1 Implement ConfirmDialog component
  - Vant Dialog wrapper
  - Props: title, message, confirmText, cancelText
  - Emit: confirm, cancel
- [ ] 12.2 Use ConfirmDialog for:
  - Device deletion
  - Secret rotation
  - Logout (optional)

## 13. Device List Enhancements
- [ ] 13.1 Improve device list UI
  - Show device status badge
  - Show last updated time
  - Swipe actions (edit/delete) - optional
- [ ] 13.2 Add device filtering (optional)
  - Show all / enabled only / disabled only

## 14. Profile Page Features
- [ ] 14.1 Display user statistics
  - Total devices
  - Enabled devices
  - Account age
- [ ] 14.2 Implement logout
  - Clear auth store
  - Clear localStorage
  - Navigate to /login

## 15. Testing
- [ ] 15.1 Test registration flow
  - Valid registration
  - Duplicate username
  - Validation errors
- [ ] 15.2 Test device CRUD
  - Create device
  - Edit device
  - Delete device
  - Rotate secret
- [ ] 15.3 Test navigation
  - Bottom nav works
  - Route guards work
  - Back navigation works
- [ ] 15.4 Test on mobile
  - Touch interactions
  - Responsive layout
  - Bottom nav positioning

## 16. TypeScript and Build
- [ ] 16.1 Ensure all new code is type-safe
- [ ] 16.2 Run type-check
- [ ] 16.3 Test production build
- [ ] 16.4 Fix any linting issues

## 17. Documentation
- [ ] 17.1 Update `web/README.md`
  - Document new pages
  - Update feature list
- [ ] 17.2 Update `web/TEST_PLAN.md`
  - Add device management test scenarios
  - Add registration test scenarios
- [ ] 17.3 Create user guide (optional)
  - How to register
  - How to add devices
  - How to manage devices

## 18. UI/UX Polish
- [ ] 18.1 Consistent styling
  - Match existing design
  - Proper spacing and alignment
- [ ] 18.2 Icons and visual feedback
  - Appropriate icons for actions
  - Loading spinners
  - Success/error states
- [ ] 18.3 Accessibility
  - Proper labels
  - Touch target sizes
  - Color contrast

## 19. Edge Cases
- [ ] 19.1 Handle empty states
  - No devices
  - First-time user
- [ ] 19.2 Handle errors gracefully
  - Network errors
  - API errors
  - Validation errors
- [ ] 19.3 Handle concurrent operations
  - Multiple tabs
  - Race conditions

## 20. Final Validation
- [ ] 20.1 Complete user journey test
  - Register → Login → Create Device → Edit → Delete
- [ ] 20.2 All TypeScript checks pass
- [ ] 20.3 Production build succeeds
- [ ] 20.4 Mobile responsive verified
- [ ] 20.5 All success criteria met
