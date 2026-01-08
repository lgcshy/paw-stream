# API Server Specification - Phase 3 Implementation

## ADDED Requirements

### Requirement: User Registration API
The API server SHALL provide a registration endpoint for creating new user accounts.

#### Scenario: Successful registration
- **WHEN** `POST /api/register` is called with valid username and password
- **THEN** a new user account is created with hashed password
- **AND** returns 201 Created with user info (excluding password)
- **AND** username is unique in the system

#### Scenario: Duplicate username
- **WHEN** registration is attempted with an existing username
- **THEN** returns 409 Conflict
- **AND** includes error message indicating username is taken

#### Scenario: Invalid input
- **WHEN** username or password is missing or invalid
- **THEN** returns 400 Bad Request
- **AND** includes validation error details

---

### Requirement: User Login API
The API server SHALL provide a login endpoint that authenticates users and issues JWT tokens.

#### Scenario: Successful login
- **WHEN** `POST /api/login` is called with valid credentials
- **THEN** returns 200 OK with JWT token
- **AND** token includes user_id and username claims
- **AND** token is valid for configured expiry time (default 24h)

#### Scenario: Invalid credentials
- **WHEN** login is attempted with wrong password
- **THEN** returns 401 Unauthorized
- **AND** does not reveal whether username exists

#### Scenario: Disabled user
- **WHEN** a disabled user attempts to login
- **THEN** returns 403 Forbidden
- **AND** includes message that account is disabled

---

### Requirement: Current User Info API
The API server SHALL provide an endpoint to retrieve current authenticated user information.

#### Scenario: Get current user
- **WHEN** `GET /api/me` is called with valid JWT token
- **THEN** returns 200 OK with user profile
- **AND** response includes id, username, nickname, created_at
- **AND** password_hash is never exposed

#### Scenario: Unauthenticated request
- **WHEN** `/api/me` is called without valid token
- **THEN** returns 401 Unauthorized

---

### Requirement: Device Creation API
The API server SHALL provide an endpoint to register new streaming devices.

#### Scenario: Create device
- **WHEN** `POST /api/devices` is called with device name and location
- **THEN** returns 201 Created with device info and secret
- **AND** device is associated with authenticated user
- **AND** unique device_id and device_secret are generated
- **AND** publish_path is set to `dogcam/<device_id>`
- **AND** device_secret is only returned once in response

#### Scenario: Missing required fields
- **WHEN** device creation is attempted without name
- **THEN** returns 400 Bad Request
- **AND** includes validation error details

---

### Requirement: Device List API
The API server SHALL provide an endpoint to list user's devices.

#### Scenario: List own devices
- **WHEN** `GET /api/devices` is called by authenticated user
- **THEN** returns 200 OK with array of user's devices
- **AND** each device includes id, name, location, publish_path, disabled status
- **AND** device_secret is never exposed in list response
- **AND** results are ordered by created_at descending

#### Scenario: Empty device list
- **WHEN** user has no devices
- **THEN** returns 200 OK with empty array

---

### Requirement: Device Details API
The API server SHALL provide an endpoint to retrieve detailed device information.

#### Scenario: Get own device
- **WHEN** `GET /api/devices/:id` is called for user's device
- **THEN** returns 200 OK with device details
- **AND** response includes all device fields except secret

#### Scenario: Access other user's device
- **WHEN** user attempts to access another user's device
- **THEN** returns 404 Not Found
- **AND** does not reveal device existence

#### Scenario: Device not found
- **WHEN** non-existent device_id is requested
- **THEN** returns 404 Not Found

---

### Requirement: Device Update API
The API server SHALL provide an endpoint to update device information.

#### Scenario: Update device name/location
- **WHEN** `PUT /api/devices/:id` is called with updated fields
- **THEN** returns 200 OK with updated device info
- **AND** only name, location, and disabled fields can be updated
- **AND** updated_at timestamp is refreshed

#### Scenario: Enable/disable device
- **WHEN** device disabled field is updated
- **THEN** device publish authorization is affected
- **AND** disabled device cannot publish streams

#### Scenario: Update non-owned device
- **WHEN** user attempts to update another user's device
- **THEN** returns 404 Not Found

---

### Requirement: Device Deletion API
The API server SHALL provide an endpoint to delete devices.

#### Scenario: Delete own device
- **WHEN** `DELETE /api/devices/:id` is called for user's device
- **THEN** returns 204 No Content
- **AND** device is removed from database
- **AND** associated publish_path becomes unavailable

#### Scenario: Delete non-owned device
- **WHEN** user attempts to delete another user's device
- **THEN** returns 404 Not Found

---

### Requirement: Device Secret Rotation API
The API server SHALL provide an endpoint to rotate device secrets for security.

#### Scenario: Rotate secret
- **WHEN** `POST /api/devices/:id/rotate-secret` is called
- **THEN** returns 200 OK with new device_secret
- **AND** old secret is invalidated immediately
- **AND** secret_version is incremented
- **AND** new secret is only returned once

#### Scenario: Rotate non-owned device secret
- **WHEN** user attempts to rotate another user's device secret
- **THEN** returns 404 Not Found

---

### Requirement: Path List API
The API server SHALL provide an endpoint to list accessible stream paths.

#### Scenario: List accessible paths
- **WHEN** `GET /api/paths` is called by authenticated user
- **THEN** returns 200 OK with array of accessible paths
- **AND** each entry includes publish_path, device info (name, location)
- **AND** only enabled devices are included
- **AND** results are suitable for stream selection UI

---

### Requirement: MediaMTX Publish Authorization
The API server SHALL validate device publish requests from MediaMTX callbacks.

#### Scenario: Valid device publish
- **WHEN** MediaMTX calls `/mediamtx/auth` with action=publish and valid device_secret
- **THEN** returns 200 OK
- **AND** logs successful authorization

#### Scenario: Invalid device secret
- **WHEN** publish is attempted with wrong secret
- **THEN** returns 403 Forbidden
- **AND** logs authorization failure with reason

#### Scenario: Disabled device publish
- **WHEN** disabled device attempts to publish
- **THEN** returns 403 Forbidden
- **AND** includes reason in response

---

### Requirement: MediaMTX Read Authorization
The API server SHALL validate user playback requests from MediaMTX callbacks.

#### Scenario: Valid user playback
- **WHEN** MediaMTX calls `/mediamtx/auth` with action=read and valid user_token
- **THEN** returns 200 OK if user owns device for that path
- **AND** logs successful authorization

#### Scenario: Invalid user token
- **WHEN** playback is attempted with invalid token
- **THEN** returns 403 Forbidden
- **AND** logs authorization failure

#### Scenario: Unauthorized path access
- **WHEN** user attempts to view path they don't own
- **THEN** returns 403 Forbidden
- **AND** includes reason in response

---

## MODIFIED Requirements

### Requirement: HTTP API Handlers (Stubs)
The API server SHALL provide complete HTTP handlers for business API endpoints.

#### Scenario: Auth endpoints implemented
- **WHEN** `POST /api/register` or `POST /api/login` is requested
- **THEN** the handler executes full authentication logic
- **AND** returns appropriate success or error responses
- **AND** no longer returns 501 Not Implemented

#### Scenario: Device endpoints implemented
- **WHEN** device management endpoints are requested
- **THEN** handlers execute full CRUD logic
- **AND** enforce user ownership permissions
- **AND** no longer returns 501 Not Implemented
