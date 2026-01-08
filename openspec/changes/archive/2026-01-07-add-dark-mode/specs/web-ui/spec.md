# Web UI Specification - Dark Mode Support

## ADDED Requirements

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
