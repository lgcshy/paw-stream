# web-ui Specification

## Purpose

The Web UI capability provides a mobile-first, TypeScript-based frontend interface for PawStream. Built with Vue 3, Vite 7, and Vant 4, it enables users to authenticate, view live camera streams, and monitor their pets through a responsive web application accessible from mobile browsers.

This specification defines the project structure, development workflows, component architecture, and user interface requirements for the PawStream web frontend.
## Requirements
### Requirement: Project Foundation
The web UI SHALL be built with Vue 3, TypeScript, Vite 7, and Vant 4 (Vue 3 compatible) following mobile-first design principles and type-safe development practices.

#### Scenario: Development server startup
- **WHEN** developer runs `npm run dev` in the web directory
- **THEN** the Vite development server starts successfully on the configured port
- **AND** hot module replacement is enabled for live updates
- **AND** TypeScript compilation is enabled in development mode

#### Scenario: Type checking
- **WHEN** developer runs `npm run type-check`
- **THEN** TypeScript compiler validates all type definitions
- **AND** no type errors are reported
- **AND** strict mode type checking is enforced

#### Scenario: Mobile viewport rendering
- **WHEN** the application is viewed in a mobile browser or mobile viewport simulation
- **THEN** the UI renders correctly at mobile screen sizes (320px - 768px width)
- **AND** Vant 4 components display with proper mobile styling
- **AND** all TypeScript-defined component props work correctly

#### Scenario: Production build
- **WHEN** developer runs `npm run build`
- **THEN** the project builds successfully without TypeScript errors
- **AND** optimized static assets are generated in the dist directory
- **AND** type-checked and transpiled JavaScript is produced

---

### Requirement: Project Structure
The web project SHALL follow a clear directory structure separating views, components, routing, stores, and TypeScript type definitions.

#### Scenario: Code organization
- **WHEN** developers add new features
- **THEN** page components are placed in `src/views/` with TypeScript script setup
- **AND** reusable components are placed in `src/components/` with TypeScript
- **AND** routing configuration is in `src/router/` as TypeScript files
- **AND** state management stores are in `src/stores/` as TypeScript files
- **AND** shared type definitions are in `src/types/`

#### Scenario: TypeScript configuration
- **WHEN** importing modules using path aliases
- **THEN** the `@/` alias resolves to `src/` directory
- **AND** TypeScript compiler recognizes the path mappings
- **AND** IDE provides proper autocomplete and type inference

#### Scenario: Asset management
- **WHEN** static assets are needed
- **THEN** they are placed in `src/assets/` or `public/` based on processing needs
- **AND** Vite handles asset optimization automatically
- **AND** TypeScript declarations for asset imports are provided

---

### Requirement: Code Formatting
The web project SHALL use Prettier with default configuration for consistent code formatting across TypeScript and Vue files.

#### Scenario: Automated formatting
- **WHEN** developer runs Prettier on the codebase
- **THEN** all Vue, TypeScript, JavaScript, CSS files are formatted consistently
- **AND** no formatting conflicts occur between team members
- **AND** TypeScript syntax is preserved correctly

---

### Requirement: Navigation Structure
The web UI SHALL provide basic routing between placeholder pages for login, stream list, and stream player.

#### Scenario: Page navigation
- **WHEN** user navigates between pages using Vue Router
- **THEN** the URL updates correctly
- **AND** the appropriate view component renders
- **AND** navigation transitions are smooth on mobile devices

#### Scenario: Initial route
- **WHEN** user first visits the application
- **THEN** they are directed to the login page if not authenticated
- **OR** they are directed to the stream list if authenticated (placeholder logic)

---

### Requirement: Mobile Layout
The web UI SHALL include a base layout component with mobile-optimized navigation bar using Vant 4 components with proper TypeScript typing.

#### Scenario: Layout rendering
- **WHEN** any page is displayed
- **THEN** it includes the base layout with navigation bar
- **AND** the layout adapts to the screen size
- **AND** Vant 4 NavBar component is used for consistent mobile UI
- **AND** all component props are type-checked by TypeScript

#### Scenario: Touch interactions
- **WHEN** user interacts with UI elements on a touch device
- **THEN** all buttons and interactive elements have appropriate touch target sizes (minimum 44x44px)
- **AND** touch feedback is provided via Vant 4 component animations
- **AND** TypeScript ensures type-safe event handlers

---

### Requirement: WebRTC Player Placeholder
The web UI SHALL include a TypeScript-based placeholder component for the WebRTC video player to be implemented in Phase 4.

#### Scenario: Player component structure
- **WHEN** the StreamPlayer view is accessed
- **THEN** it displays a placeholder indicating where WebRTC player will be integrated
- **AND** the component structure is ready for WebRTC SDK integration
- **AND** the component accepts stream ID as a typed prop (string)
- **AND** TypeScript interfaces define the expected player API

#### Scenario: Type-safe props
- **WHEN** StreamPlayer component is used in parent components
- **THEN** TypeScript validates the stream ID prop type
- **AND** IDE provides autocomplete for component props
- **AND** compilation fails if incorrect prop types are passed

---

### Requirement: Development Documentation
The web project SHALL include a README with setup instructions, available scripts, TypeScript guidelines, and development conventions.

#### Scenario: New developer onboarding
- **WHEN** a new developer reads web/README.md
- **THEN** they understand how to install dependencies
- **AND** they know how to start the development server
- **AND** they understand how to run TypeScript type checking
- **AND** they understand the project structure and TypeScript conventions
- **AND** they know how to use Vant 4 components with TypeScript
- **AND** they know where to find OpenSpec project conventions

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

### Requirement: Theme System
The web UI SHALL provide a theme system supporting light mode, dark mode, and automatic mode based on system preferences.

#### Scenario: Theme initialization
- **WHEN** the application loads
- **THEN** the theme preference is loaded from localStorage
- **AND** if no preference is stored, the system theme preference is detected
- **AND** the appropriate theme is applied to the UI

#### Scenario: System theme detection
- **WHEN** the theme mode is set to "auto"
- **THEN** the application detects the system theme using `prefers-color-scheme`
- **AND** applies dark mode if the system prefers dark
- **AND** applies light mode if the system prefers light

#### Scenario: System theme change
- **WHEN** the user changes their system theme preference
- **AND** the application theme mode is set to "auto"
- **THEN** the application automatically switches to match the new system theme
- **AND** no page reload is required

---

### Requirement: Theme State Management
The web UI SHALL use a Pinia store to manage theme state and persistence.

#### Scenario: Theme store initialization
- **WHEN** the application initializes
- **THEN** a theme store is created with state: `mode` (light/dark/auto) and `effectiveTheme` (light/dark)
- **AND** the store loads saved preferences from localStorage
- **AND** the store provides actions to set and toggle theme

#### Scenario: Theme persistence
- **WHEN** the user changes the theme mode
- **THEN** the preference is immediately saved to localStorage
- **AND** the preference persists across browser sessions
- **AND** the preference is applied on next application load

---

### Requirement: Theme Switching UI
The web UI SHALL provide a theme switching interface in the user profile page.

#### Scenario: Access theme settings
- **WHEN** the user navigates to the profile page
- **THEN** a "主题设置" cell is visible in the settings section
- **AND** the cell displays the current theme mode

#### Scenario: Select theme mode
- **WHEN** the user taps on the theme settings cell
- **THEN** a selection dialog opens with three options: "跟随系统", "浅色模式", "深色模式"
- **AND** the current selection is highlighted
- **AND** selecting an option immediately applies the theme
- **AND** the dialog closes after selection

#### Scenario: Theme labels
- **WHEN** displaying theme modes to the user
- **THEN** "auto" mode is labeled "跟随系统"
- **AND** "light" mode is labeled "浅色模式"
- **AND** "dark" mode is labeled "深色模式"

---

### Requirement: Dark Theme CSS Variables
The web UI SHALL define a complete set of CSS variables for dark theme styling.

#### Scenario: CSS variable definition
- **WHEN** dark theme is applied
- **THEN** the root element has class "dark-theme"
- **AND** CSS variables are defined for:
  - Background colors (primary, secondary, card)
  - Text colors (primary, secondary, disabled)
  - Border colors
  - Component-specific colors (navbar, tabbar, buttons)
  - Shadow colors
- **AND** all colors ensure sufficient contrast for readability

#### Scenario: Component styling
- **WHEN** a component uses semantic color variables
- **THEN** the component automatically adapts to the active theme
- **AND** no theme-specific logic is required in component code

---

### Requirement: Dark Theme Visual Consistency
The web UI SHALL ensure visual consistency and readability in dark mode across all pages and components.

#### Scenario: Page adaptation
- **WHEN** dark theme is active
- **THEN** all pages (Login, Register, Streams, Devices, Profile) display with appropriate dark colors
- **AND** text is clearly readable against dark backgrounds
- **AND** interactive elements are clearly visible
- **AND** gradients and animations are adapted for dark theme

#### Scenario: Vant component integration
- **WHEN** dark theme is active
- **THEN** Vant 4 components (NavBar, Tabbar, Cell, Button, etc.) display with dark theme styling
- **AND** Vant CSS variables are configured for dark theme
- **AND** component colors remain consistent with the application theme

#### Scenario: Special elements
- **WHEN** dark theme is active
- **THEN** video player controls remain visible
- **AND** gradient backgrounds (login, register, profile headers) are adapted to darker tones
- **AND** floating animation elements blend naturally with the dark background

---

### Requirement: Theme Switching Performance
The web UI SHALL provide smooth and immediate theme switching without visual glitches.

#### Scenario: Instant theme application
- **WHEN** the user changes the theme mode
- **THEN** the new theme is applied immediately without page reload
- **AND** no white flash or layout shift occurs during the switch

#### Scenario: Initial load optimization
- **WHEN** the application loads with dark theme preference
- **THEN** the correct theme is applied before rendering
- **AND** no flash of light theme occurs during initialization

---

### Requirement: Accessibility and Contrast
The web UI dark theme SHALL meet accessibility standards for color contrast.

#### Scenario: Text contrast
- **WHEN** dark theme is active
- **THEN** primary text has a contrast ratio of at least 4.5:1 against the background
- **AND** secondary text has a contrast ratio of at least 3:1 against the background

#### Scenario: Interactive element visibility
- **WHEN** dark theme is active
- **THEN** buttons, links, and interactive elements are clearly distinguishable
- **AND** focus indicators are visible
- **AND** disabled states are clearly differentiated from enabled states

