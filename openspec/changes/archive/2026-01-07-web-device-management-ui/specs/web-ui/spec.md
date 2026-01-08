# Web UI Specification - Device Management and User Registration

## ADDED Requirements

### Requirement: User Registration Page
The web UI SHALL provide a registration page for new users to create accounts.

#### Scenario: Successful registration
- **WHEN** user fills in username, password, and nickname on registration page
- **THEN** the `/api/register` endpoint is called
- **AND** the user account is created
- **AND** the user is automatically logged in
- **AND** redirected to the stream list page

#### Scenario: Username already exists
- **WHEN** user tries to register with an existing username
- **THEN** a 409 error is returned
- **AND** an error message is displayed: "用户名已被占用"

#### Scenario: Password strength validation
- **WHEN** user enters a password shorter than 6 characters
- **THEN** a validation error is shown before submission
- **AND** the form cannot be submitted

---

### Requirement: Device Creation UI
The web UI SHALL provide a form for users to create new devices.

#### Scenario: Create device successfully
- **WHEN** user fills in device name and location
- **THEN** the `/api/devices` endpoint is called
- **AND** a new device is created
- **AND** the device_secret is displayed once
- **AND** a copy button is provided for the secret
- **AND** user is warned that the secret will not be shown again

#### Scenario: Copy device secret
- **WHEN** user clicks the copy button next to device_secret
- **THEN** the secret is copied to clipboard
- **AND** a success toast is shown: "已复制到剪贴板"

#### Scenario: Navigate after creation
- **WHEN** device is created successfully
- **THEN** user can choose to create another device or return to list
- **AND** the new device appears in the device list

---

### Requirement: Device Editing UI
The web UI SHALL allow users to edit device information.

#### Scenario: Edit device name and location
- **WHEN** user navigates to device edit page
- **THEN** current device information is pre-filled
- **AND** user can modify name and location
- **AND** changes are saved via `PUT /api/devices/:id`

#### Scenario: Enable/disable device
- **WHEN** user toggles the device enabled/disabled switch
- **THEN** the device status is updated immediately
- **AND** disabled devices do not appear in stream list

---

### Requirement: Device Deletion UI
The web UI SHALL provide a way to delete devices with confirmation.

#### Scenario: Delete device with confirmation
- **WHEN** user clicks delete button on device detail page
- **THEN** a confirmation dialog is shown
- **AND** the dialog warns about permanent deletion
- **AND** user must confirm to proceed

#### Scenario: Confirm deletion
- **WHEN** user confirms deletion
- **THEN** `DELETE /api/devices/:id` is called
- **AND** the device is removed
- **AND** user is redirected to device list
- **AND** a success toast is shown

#### Scenario: Cancel deletion
- **WHEN** user cancels the deletion dialog
- **THEN** no API call is made
- **AND** the device remains unchanged

---

### Requirement: Device Secret Rotation UI
The web UI SHALL allow users to rotate device secrets for security.

#### Scenario: Rotate secret
- **WHEN** user clicks "轮换 Secret" button on device detail page
- **THEN** a confirmation dialog is shown
- **AND** the dialog warns that old secret will be invalidated
- **AND** user must confirm to proceed

#### Scenario: Confirm secret rotation
- **WHEN** user confirms secret rotation
- **THEN** `POST /api/devices/:id/rotate-secret` is called
- **AND** the new secret is displayed once
- **AND** a copy button is provided
- **AND** user is warned to update device configuration

---

### Requirement: Device Detail Page
The web UI SHALL provide a detailed view of device information.

#### Scenario: View device details
- **WHEN** user navigates to device detail page
- **THEN** all device information is displayed (name, location, publish_path, status, timestamps)
- **AND** enable/disable switch is available
- **AND** edit button is available
- **AND** rotate secret button is available
- **AND** delete button is available

#### Scenario: Quick actions from detail page
- **WHEN** user is on device detail page
- **THEN** user can quickly navigate to play stream
- **AND** user can edit device
- **AND** user can rotate secret
- **AND** user can delete device

---

### Requirement: User Profile Page
The web UI SHALL provide a user profile page showing account information.

#### Scenario: View profile
- **WHEN** user navigates to profile page
- **THEN** user information is displayed (username, nickname, registration date)
- **AND** device count statistics are shown
- **AND** logout button is available

#### Scenario: Logout from profile
- **WHEN** user clicks logout button
- **THEN** a confirmation dialog is shown (optional)
- **AND** user is logged out
- **AND** redirected to login page

---

### Requirement: Bottom Navigation
The web UI SHALL provide a bottom navigation bar for easy access to main sections.

#### Scenario: Navigate between sections
- **WHEN** user is on any main page
- **THEN** a bottom navigation bar is visible
- **AND** it shows icons for: 流列表, 设备管理, 个人中心
- **AND** the current section is highlighted
- **AND** tapping an icon navigates to that section

---

## ADDED Requirements (Continued)

### Requirement: Login Page with Registration Link
The login view SHALL include a link to the registration page.

#### Scenario: Navigate to registration
- **WHEN** user is on login page
- **THEN** a "注册账号" link is visible
- **AND** clicking it navigates to `/register`

---

### Requirement: Stream List Reorganization
The stream list view SHALL be reorganized to focus on viewing streams, with device management moved to a separate section.

#### Scenario: Focused stream list
- **WHEN** user navigates to stream list
- **THEN** only enabled devices with active streams are shown
- **AND** the list is optimized for quick access to playback
- **AND** device management actions are in the separate devices section

---

### Requirement: Bottom Navigation Bar
The layout component SHALL include a bottom navigation bar for main sections.

#### Scenario: Bottom navigation visibility
- **WHEN** user is on a main page (streams, devices, profile)
- **THEN** the bottom navigation bar is visible
- **AND** it provides quick access to all main sections

#### Scenario: Hide navigation on sub-pages
- **WHEN** user is on a sub-page (player, device form, device detail)
- **THEN** the bottom navigation bar is hidden
- **AND** a back button is provided instead

---

## ADDED Components

### Requirement: Secret Display Component
The web UI SHALL provide a reusable component for displaying and copying secrets.

#### Scenario: Display secret securely
- **WHEN** a secret is displayed
- **THEN** it is shown in a monospace font
- **AND** a copy button is next to it
- **AND** a warning about one-time display is shown

#### Scenario: Copy secret
- **WHEN** user clicks copy button
- **THEN** the secret is copied to clipboard
- **AND** a success feedback is shown

---

### Requirement: Confirm Dialog Component
The web UI SHALL provide a reusable confirmation dialog component.

#### Scenario: Show confirmation
- **WHEN** a destructive action is initiated
- **THEN** a confirmation dialog is shown
- **AND** it explains the action and consequences
- **AND** it provides confirm and cancel buttons

#### Scenario: Confirm action
- **WHEN** user clicks confirm
- **THEN** the action is executed
- **AND** the dialog closes

#### Scenario: Cancel action
- **WHEN** user clicks cancel or backdrop
- **THEN** no action is taken
- **AND** the dialog closes
