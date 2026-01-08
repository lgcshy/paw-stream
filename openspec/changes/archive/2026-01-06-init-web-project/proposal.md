# Change: Initialize Web Frontend Project

## Why

The PawStream system requires a mobile-friendly web UI to allow users to view live camera streams from their pets. Currently, the `web/` directory is empty, and we need to scaffold a complete Vue 3 + Vite project with the necessary dependencies and project structure.

This is a foundational change that enables Phase 4 of the development plan (Web UI MVP).

## What Changes

- Initialize a new Vue 3 + TypeScript project using Vite 7
- Install and configure Vant 4 UI component library (Vue 3 compatible) for mobile-first design
- Set up TypeScript configuration with strict type checking
- Set up project structure with routing and state management preparation
- Configure Prettier for code formatting
- Add basic layout structure and navigation scaffolding
- Create development server configuration
- Add WebRTC player component foundation (placeholder for actual implementation)

**No breaking changes** - this is a new capability.

## Impact

### Affected Specs
- **NEW**: `web-ui` - Mobile-friendly user interface capability

### Affected Code
- `web/` directory - will be populated with Vue 3 + TypeScript project files
- New files:
  - `web/package.json` - dependency management
  - `web/tsconfig.json` - TypeScript configuration
  - `web/vite.config.ts` - build configuration (TypeScript)
  - `web/src/` - source code directory (TypeScript)
  - `web/public/` - static assets
  - `web/env.d.ts` - TypeScript environment declarations

### Dependencies Added
- Vue 3 (framework)
- TypeScript (type-safe development)
- Vite 7 (build tool with TypeScript support)
- Vant 4 (mobile UI library, Vue 3 compatible)
- Vue Router 4 (routing with TypeScript support)
- Pinia (state management with TypeScript, optional for later)
- Prettier (code formatting)

## Success Criteria

- `npm run dev` starts development server successfully
- Basic layout renders in mobile viewport
- Project follows conventions defined in `openspec/project.md`
- All linting and formatting tools work correctly
