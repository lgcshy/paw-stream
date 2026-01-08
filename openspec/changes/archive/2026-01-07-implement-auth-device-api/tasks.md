# Implementation Tasks: User Authentication and Device Management APIs

## 1. User Authentication Implementation
- [ ] 1.1 Implement Register handler (`POST /api/register`)
  - Input validation (username, password)
  - Password strength check
  - Duplicate username check
  - User creation via UserService
  - Return user info (exclude password)
- [ ] 1.2 Implement Login handler (`POST /api/login`)
  - Credentials validation
  - UserService.Login call
  - JWT token generation
  - Return token in response
- [ ] 1.3 Implement GetMe handler (`GET /api/me`)
  - Extract user_id from JWT context
  - UserService.GetByID call
  - Return user profile

## 2. Device Management Implementation
- [ ] 2.1 Implement Create Device handler (`POST /api/devices`)
  - Extract user_id from JWT
  - Validate input (name required)
  - DeviceService.Create call
  - Return device + secret (one-time)
- [ ] 2.2 Implement List Devices handler (`GET /api/devices`)
  - Extract user_id from JWT
  - DeviceService.ListByOwner call
  - Filter out secrets
  - Return device array
- [ ] 2.3 Implement Get Device handler (`GET /api/devices/:id`)
  - Extract user_id from JWT
  - DeviceService.GetByID call
  - Verify ownership
  - Return 404 if not owned or not found
- [ ] 2.4 Implement Update Device handler (`PUT /api/devices/:id`)
  - Extract user_id from JWT
  - DeviceService.GetByID and verify ownership
  - Parse update fields (name, location, disabled)
  - DeviceService.Update call
  - Return updated device
- [ ] 2.5 Implement Delete Device handler (`DELETE /api/devices/:id`)
  - Extract user_id from JWT
  - DeviceService.GetByID and verify ownership
  - DeviceService.Delete call (or add Delete to service)
  - Return 204 No Content
- [ ] 2.6 Implement Rotate Secret handler (`POST /api/devices/:id/rotate-secret`)
  - Extract user_id from JWT
  - Verify ownership
  - DeviceService.RotateSecret call
  - Return new secret (one-time)

## 3. Path Query Implementation
- [ ] 3.1 Create path_handler.go
- [ ] 3.2 Implement List Paths handler (`GET /api/paths`)
  - Extract user_id from JWT
  - DeviceService.ListByOwner call
  - Filter enabled devices only
  - Map to path response format (path + device info)
  - Return path array

## 4. MediaMTX Integration Enhancement
- [ ] 4.1 Complete Publish authorization in MediaMTXHandler.Auth
  - Already implemented, verify logic
  - Add detailed logging
  - Test with actual device secret
- [ ] 4.2 Complete Read/Playback authorization in MediaMTXHandler.Auth
  - Already implemented, verify logic
  - Add detailed logging
  - Test with actual user token

## 5. Request/Response Types
- [ ] 5.1 Define AuthHandler request/response types
  - RegisterRequest (username, password, nickname)
  - LoginRequest (username, password)
  - LoginResponse (token, user)
- [ ] 5.2 Define DeviceHandler request/response types
  - CreateDeviceRequest (name, location)
  - CreateDeviceResponse (device + secret)
  - UpdateDeviceRequest (name, location, disabled)
  - DeviceResponse (device without secret)
- [ ] 5.3 Define PathHandler response types
  - PathInfo (publish_path, device_name, device_location, device_id)

## 6. Error Handling
- [ ] 6.1 Add validation error helper
  - ValidationError struct
  - Field-level error messages
- [ ] 6.2 Add ownership check helper
  - Consistent 404 for unauthorized access
- [ ] 6.3 Improve error responses
  - Consistent error format
  - Request ID in all errors

## 7. Route Registration
- [ ] 7.1 Update routes.go to wire new handlers
  - Register auth endpoints (no auth required)
  - Register device endpoints (auth required)
  - Register path endpoint (auth required)
- [ ] 7.2 Verify middleware order
  - Recovery → RequestID → Logger → CORS → routes

## 8. Service Layer Enhancements (if needed)
- [ ] 8.1 Review UserService
  - Ensure Register, Login, GetByID are sufficient
- [ ] 8.2 Review DeviceService
  - Ensure all CRUD operations exist
  - Add Delete method if missing
- [ ] 8.3 Review ACLService
  - Verify publish/read authorization logic
  - Add device disabled check

## 9. Testing
- [ ] 9.1 Manual API testing with curl/Postman
  - Test user registration flow
  - Test user login flow
  - Test device creation flow
  - Test device CRUD operations
  - Test path listing
- [ ] 9.2 Test MediaMTX integration
  - Test publish with device secret
  - Test playback with user token
  - Test authorization failures
- [ ] 9.3 Test edge cases
  - Duplicate username
  - Invalid credentials
  - Unauthorized access
  - Disabled device/user
- [ ] 9.4 Write integration tests (optional for Phase 3)
  - auth_test.go
  - device_test.go
  - mediamtx_test.go

## 10. Documentation
- [ ] 10.1 Update server/api/README.md
  - Document all API endpoints
  - Add request/response examples
  - Add authentication instructions
- [ ] 10.2 Create API usage examples
  - curl commands for each endpoint
  - Postman collection (optional)
- [ ] 10.3 Update VERIFICATION_REPORT.md
  - Add Phase 3 test results

## 11. End-to-End Flow Validation
- [ ] 11.1 Test complete user journey
  - Register → Login → Create Device → List Devices
- [ ] 11.2 Test complete streaming flow
  - Create device → Get secret → Publish (simulated)
  - Login → Get token → Read/playback (simulated)
- [ ] 11.3 Test MediaMTX callback flow
  - Configure MediaMTX to use auth callback
  - Test actual publish authorization
  - Test actual playback authorization

## 12. Code Quality
- [ ] 12.1 Run go fmt ./...
- [ ] 12.2 Run go vet ./...
- [ ] 12.3 Check for code duplication
- [ ] 12.4 Add code comments for complex logic
- [ ] 12.5 Review error handling completeness

## 13. Performance Considerations
- [ ] 13.1 Review database query efficiency
- [ ] 13.2 Consider adding database indexes if needed
- [ ] 13.3 Test concurrent request handling

## 14. Security Review
- [ ] 14.1 Verify password is never logged
- [ ] 14.2 Verify device_secret is only returned once
- [ ] 14.3 Verify JWT tokens expire correctly
- [ ] 14.4 Verify ownership checks on all operations
- [ ] 14.5 Review CORS configuration for production

## 15. Final Validation
- [ ] 15.1 All handlers return appropriate status codes
- [ ] 15.2 All endpoints have proper authentication
- [ ] 15.3 No 501 Not Implemented responses remain
- [ ] 15.4 Integration with MediaMTX works end-to-end
- [ ] 15.5 Documentation is accurate and complete
