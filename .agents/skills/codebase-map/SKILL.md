---
name: codebase-map
description: Generates a comprehensive, LLM-readable codebase map with feature index, architecture overview, and navigation guide. Creates CODEBASE_MAP.md that other skills can reference. Triggered by "map the codebase", "document architecture", or "create codebase overview" requests.
---

# Codebase Map Skill

## When to use
- User asks to "map the codebase", "document architecture", or "create codebase overview"
- Before starting work on a new feature in an unfamiliar codebase
- When other skills need a comprehensive understanding of the project structure
- To create a reusable reference document for future LLM context

## Output
Generates `CodebaseMap/CODEBASE_MAP.md` — a structured, navigable document optimized for LLM consumption.

## Instructions

### 1. Analyze project structure and technology stack

**Detect project type:**
- Check for `go.mod` → Go project
- Check for `package.json` → Node.js/TypeScript project
- Check for `build.gradle` / `build.gradle.kts` → Android/Kotlin/Java project
- Check for `Cargo.toml` → Rust project
- Check for `requirements.txt` / `pyproject.toml` → Python project
- Check for `pubspec.yaml` → Flutter/Dart project

**Identify frameworks and libraries:**
- Parse dependency files (e.g., `go.mod`, `package.json`, `build.gradle`)
- Note major frameworks (e.g., Cobra, React Native, Jetpack Compose, Flask, Express)
- Flag UI/CLI frameworks, database drivers, HTTP clients, testing frameworks

**Map directory structure:**
```bash
# Get high-level directory tree (2-3 levels deep)
tree -L 3 -d -I 'node_modules|.git|vendor|build|dist'
# Windows fallback:
Get-ChildItem -Recurse -Directory -Depth 2 | Select-Object FullName
```

Categorize directories:
- **Source code:** `src/`, `lib/`, `app/`, `cmd/`, `internal/`, `pkg/`
- **Tests:** `test/`, `tests/`, `__tests__/`, `*_test.go`
- **Configuration:** `config/`, `.github/`, `configs/`
- **Documentation:** `docs/`, `doc/`, `README.md`, `CHANGELOG.md`
- **Assets:** `assets/`, `static/`, `public/`, `res/`
- **Build/Deploy:** `build/`, `scripts/`, `Makefile`, `Dockerfile`

### 2. Create the feature index

Scan the codebase for key entities and create a searchable index:

**Functions/Commands (CLI projects):**
```bash
# Go CLI example
grep -rn "func.*Command()" --include="*.go" | head -50
grep -rn "cobra.Command" --include="*.go"

# Extract command names and descriptions from Cobra/Click/argparse definitions
```

**API Endpoints (server projects):**
```bash
# Go HTTP handlers
grep -rn "http.HandleFunc\|router\.\(GET\|POST\|PUT\|DELETE\)" --include="*.go"

# Express/Node.js routes
grep -rn "app.\(get\|post\|put\|delete\)\|router.\(get\|post\)" --include="*.js" --include="*.ts"
```

**UI Components (mobile/web projects):**
```bash
# React/React Native components
grep -rn "export.*function\|export.*const.*=.*=>" --include="*.tsx" --include="*.jsx"

# Jetpack Compose
grep -rn "@Composable" --include="*.kt"

# Flutter widgets
grep -rn "class.*extends StatelessWidget\|class.*extends StatefulWidget" --include="*.dart"
```

**State Management:**
```bash
# Redux/MobX stores
find . -name "*store*" -o -name "*reducer*" -o -name "*slice*"

# Go context/state
grep -rn "type.*struct {" --include="*.go" | grep -i "state\|context\|manager"
```

**Configuration/Parameters:**
```bash
# Config structs
grep -rn "type.*Config struct" --include="*.go"
grep -rn "interface Config" --include="*.ts"

# Environment variables
grep -rn "os.Getenv\|process.env\|System.getenv" --include="*.go" --include="*.ts" --include="*.java"
```

**Assets (icons, images):**
```bash
# Find all image/icon files
find . -type f \( -name "*.png" -o -name "*.jpg" -o -name "*.svg" -o -name "*.ico" \) | head -100

# Android drawables/mipmaps
find . -path "*/res/drawable*" -o -path "*/res/mipmap*"
```

### 3. Build the CODEBASE_MAP.md document

Create `CodebaseMap/CODEBASE_MAP.md` with the following structure:

```markdown
# Codebase Map — [Project Name]
Generated: [date]
Version: [git commit hash or version number]

> **Purpose:** This document provides a comprehensive, LLM-readable overview of the codebase. Use the feature index to jump directly to relevant sections without reading the entire document.

---

## Quick Navigation
- [Technology Stack](#technology-stack)
- [Architecture Overview](#architecture-overview)
- [Feature Index](#feature-index)
- [Directory Structure](#directory-structure)
- [Key Files Reference](#key-files-reference)
- [Entry Points](#entry-points)
- [External Dependencies](#external-dependencies)
- [Configuration & Environment](#configuration--environment)
- [Testing Strategy](#testing-strategy)
- [Build & Deployment](#build--deployment)

---

## Technology Stack
- **Language:** [e.g., Go 1.25, TypeScript 5.x, Kotlin 2.0]
- **Primary Framework:** [e.g., Cobra CLI, React Native, Express]
- **UI/Rendering:** [e.g., Jetpack Compose, React, terminal output]
- **Data Layer:** [e.g., SQLite, PostgreSQL, in-memory cache]
- **Testing:** [e.g., Go testing, Jest, JUnit]
- **Build System:** [e.g., go build, Gradle, npm scripts]
- **Package Manager:** [e.g., go mod, npm, Maven]

---

## Architecture Overview
[2-3 paragraph high-level description of the system architecture]

**Key architectural patterns:**
- [e.g., CLI with subcommands, Client-server daemon, MVC, MVVM, Clean Architecture]
- [e.g., Device abstraction layer for iOS/Android]
- [e.g., WebSocket streaming for real-time data]

**Data flow:**
1. [Step 1: e.g., User invokes CLI command]
2. [Step 2: e.g., Command handler calls device manager]
3. [Step 3: e.g., Device manager communicates via adb/go-ios]
4. [Step 4: e.g., Response formatted and returned to user]

---

## Feature Index
> **Navigation tip:** Search this section for the feature you need, then jump to the referenced file.

### Commands (CLI)
| Command | Description | File | Function/Handler |
|---------|-------------|------|------------------|
| `devices` | List connected devices | `commands/devices.go` | `DevicesCmd()` |
| `screenshot` | Capture device screenshot | `commands/screenshot.go` | `ScreenshotCmd()`, `TakeScreenshot()` |
| `apps list` | List installed apps | `commands/apps.go` | `AppsListCmd()` |
| [etc.] | | | |

### API Endpoints (if applicable)
| Endpoint | Method | Description | Handler File | Function |
|----------|--------|-------------|--------------|----------|
| `/api/devices` | GET | Get device list | `server/handlers.go` | `GetDevices()` |
| [etc.] | | | | |

### UI Components (if applicable)
| Component | Purpose | File | Props/Parameters |
|-----------|---------|------|------------------|
| `DeviceCard` | Display device info | `components/DeviceCard.tsx` | `device: Device` |
| [etc.] | | | |

### State/Data Models
| Type/Class | Purpose | File | Key Fields |
|------------|---------|------|------------|
| `Device` | Represents mobile device | `types/device.go` | `ID`, `Name`, `Platform`, `State` |
| `Screenshot` | Screenshot metadata | `types/screenshot.go` | `Format`, `Quality`, `Data` |
| [etc.] | | | |

### Configuration Parameters
| Parameter | Type | Source | Description |
|-----------|------|--------|-------------|
| `MOBILECLI_DAEMON_PORT` | env var | `daemon/config.go` | Daemon server port |
| `defaultOutputPath` | const | `commands/screenshot.go` | Default screenshot path |
| [etc.] | | | |

### Assets (Icons/Images)
| Asset | Path | Used In | Purpose |
|-------|------|---------|---------|
| `ic_launcher.png` | `assets/icons/` | Android app icon | Main launcher icon |
| `logo.svg` | `assets/` | README | Project logo |
| [etc.] | | | |

---

## Directory Structure
```
[project-root]/
├── .agents/              # AI agent skills and workflows
│   └── skills/          # Skill definitions (ux-review, adb-sleuth, etc.)
├── cli/                 # CLI framework initialization
├── commands/            # Command implementations (devices, screenshot, etc.)
├── daemon/              # Background daemon for persistent connections
├── devices/             # Device abstraction layer (iOS/Android)
├── server/              # WebSocket/HTTP server for streaming
├── types/               # Shared data types and interfaces
├── utils/               # Utility functions (logging, formatting, etc.)
├── assets/              # Static assets (icons, images)
├── doc/                 # Documentation
├── docs/                # Additional docs (possibly generated)
├── test/                # Test files and fixtures
└── main.go              # Application entry point
```

**Key subdirectories:**
- `commands/` — Each file typically implements one CLI subcommand
- `devices/` — Platform-specific device management (iOS vs Android)
- `daemon/` — Long-running process for maintaining device connections
- `types/` — Shared structs and interfaces used across the codebase

---

## Key Files Reference
| File | Purpose | Jump to for... |
|------|---------|----------------|
| `main.go` | Application entry point | Understanding startup flow, signal handling |
| `cli/root.go` | Root command setup | Adding new commands, global flags |
| `commands/devices.go` | Device listing | Device enumeration logic |
| `commands/screenshot.go` | Screenshot capture | Image capture and format conversion |
| `devices/manager.go` | Device manager | Device lifecycle, connection pooling |
| `daemon/server.go` | Daemon server | WebSocket/HTTP server setup |
| `types/device.go` | Device type definition | Device model structure |
| `go.mod` | Go dependencies | Understanding external libraries |
| `Makefile` | Build automation | Build, test, install commands |

---

## Entry Points
**CLI invocation:**
- `main.go` → `cli.Execute()` → Cobra command tree

**Daemon mode:**
- `main.go` checks `daemon.IsChild()` → forks daemon process → `daemon/server.go`

**Command execution flow:**
1. User runs `mobilecli <command> [args]`
2. `cli/root.go` routes to `commands/<command>.go`
3. Command handler calls `devices.GetDeviceManager()`
4. Device-specific logic in `devices/ios.go` or `devices/android.go`
5. Result formatted and returned to stdout

---

## External Dependencies
> See `go.mod` / `package.json` for full list

**Critical dependencies:**
| Package | Purpose | Used In |
|---------|---------|---------|
| `github.com/spf13/cobra` | CLI framework | All command files |
| `github.com/danielpaulus/go-ios` | iOS device communication | `devices/ios.go` |
| `github.com/gorilla/websocket` | WebSocket server | `server/websocket.go` |
| [etc.] | | |

---

## Configuration & Environment
**Environment variables:**
- `MOBILECLI_DAEMON_PORT` — Daemon server port (default: 9898)
- `MOBILECLI_LOG_LEVEL` — Logging verbosity (debug, info, warn, error)
- `ADB_PATH` — Custom adb binary path (if not in PATH)

**Config files:**
- None (CLI uses flags and environment variables)

---

## Testing Strategy
**Test structure:**
- Unit tests: `*_test.go` files alongside source
- Integration tests: `test/` directory
- Mocks: Generated via `go generate` or manual mocks in `test/mocks/`

**Run tests:**
```bash
go test ./...
make test
```

**Coverage:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Build & Deployment
**Build from source:**
```bash
make build
# Output: ./mobilecli binary
```

**Install globally:**
```bash
make install
# Installs to $GOPATH/bin
```

**Release process:**
- Tag version: `git tag v1.2.3`
- GitHub Actions builds binaries for Linux, macOS, Windows
- Binaries published to GitHub Releases
- npm package published to `@mobilenext/mobilecli`

---

## Navigation Tips for LLMs
When working with this codebase:
1. **To add a new CLI command:** Start with `commands/` directory, reference `commands/screenshot.go` as a template
2. **To modify device communication:** Check `devices/` for platform-specific logic
3. **To understand data structures:** See `types/` directory
4. **To add a new feature:** Update Feature Index above after implementation
5. **To debug daemon issues:** Check `daemon/server.go` and logging in `utils/log.go`

**Other skills can reference this document by:**
- Reading `CodebaseMap/CODEBASE_MAP.md` before starting work
- Searching the Feature Index for specific functions/components
- Using file paths from Key Files Reference to locate relevant code
```

### 4. Generate supplementary index files (optional)

If the codebase is very large (>100 files), create supplementary indexes:

**`CodebaseMap/FUNCTION_INDEX.md`:**
```markdown
# Function Index
[Alphabetical list of all exported functions with file:line references]

## A
- `AddDevice(device Device)` — `devices/manager.go:45`
- `AppsListCmd()` — `commands/apps.go:23`

## B
- `BootDevice(deviceID string)` — `commands/boot.go:18`
[etc.]
```

**`CodebaseMap/TYPE_INDEX.md`:**
```markdown
# Type Index
[Alphabetical list of all structs/interfaces/types]

- `Device` — `types/device.go:12` — Represents a mobile device
- `Screenshot` — `types/screenshot.go:8` — Screenshot metadata
[etc.]
```

### 5. Integrate with existing skills

Add a note at the top of the main `CODEBASE_MAP.md`:

```markdown
## For Other Skills
If you are an AI skill (e.g., `design-bridge`, `adb-sleuth`, `ux-review`) and need to understand this codebase:
1. Read the [Architecture Overview](#architecture-overview) first (2-3 paragraphs)
2. Search the [Feature Index](#feature-index) for the specific feature you're working on
3. Jump to the referenced file using the [Key Files Reference](#key-files-reference)
4. Do NOT read this entire document linearly — use the navigation links
```

### 6. Final output
All deliverables reside in the `CodebaseMap/` directory:
```
CodebaseMap/
├── CODEBASE_MAP.md (main document)
├── FUNCTION_INDEX.md (optional, for large codebases)
└── TYPE_INDEX.md (optional, for large codebases)
```

## Important notes
- **Optimize for LLM reading:** Use clear section headings, tables, and navigation links
- **Keep it current:** This document should be regenerated when major architectural changes occur
- **Feature Index is critical:** This is the primary way other skills will find relevant code
- **Avoid dumping code:** Reference file paths and line numbers, but don't paste entire functions
- **Navigation over completeness:** Better to have a well-organized 80% map than a complete but unnavigable 100% dump
- **Other skills should check for this document first:** If `CodebaseMap/CODEBASE_MAP.md` exists, skills like `design-bridge` and `adb-sleuth` should read the Feature Index before doing their own scanning
- **Update instructions:** Add a note at the top: "Regenerate this document by running: `@oz map the codebase`"
- **Git ignore generated indexes:** Add `CodebaseMap/` to `.gitignore` if it's purely derived (or commit it as documentation)
- **Windows compatibility:** Use PowerShell equivalents for all bash commands (e.g., `Get-ChildItem` instead of `find`)
