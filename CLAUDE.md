# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

mobilecli is a universal command-line tool for managing iOS and Android devices, simulators, emulators, and apps. It provides both a CLI interface and an HTTP/WebSocket API using JSON-RPC 2.0.

**Module**: `github.com/mobile-next/mobilecli`
**Language**: Go 1.25.6
**CLI Framework**: Cobra

## Build Commands

```bash
make build          # Build the binary (CGO_ENABLED=0, stripped)
make test           # Run unit tests with race detection
make test-cover     # Run tests with coverage, generates HTML report
make lint           # Run golangci-lint
make fmt            # Format code with go fmt and goimports
```

For running a single test:
```bash
go test -v -run TestName ./path/to/package
```

## Architecture

The codebase follows a layered architecture with clear separation of concerns:

### CLI Layer (`cli/`)
Cobra command definitions that parse flags and delegate to the commands layer. Each file maps to a CLI command group (e.g., `apps.go`, `device.go`, `io.go`).

### Commands Layer (`commands/`)
Business logic implementation. This is where the actual work happens - CLI flags are processed here and device operations are orchestrated.

### Device Abstraction (`devices/`)
Platform-specific implementations behind a common interface:

- **`ControllableDevice`** interface (`common.go`): The core abstraction that all device types implement. Provides methods for screenshot, tap, swipe, app management, etc.
- **`android.go`**: Android device/emulator communication via ADB
- **`ios.go`**: iOS real device communication via go-ios library
- **`simulator.go`**: iOS simulator communication via Xcode tools
- **`wda/`**: WebDriverAgent client for iOS (screenshots, gestures, buttons, etc.)

### Server (`server/`)
HTTP/WebSocket API server implementing JSON-RPC 2.0:
- **`dispatch.go`**: Method registry mapping RPC method names to handlers
- **`server.go`**: HTTP server implementation
- **`websocket.go`**: WebSocket support for persistent connections

### Supporting Packages
- **`daemon/`**: Background process support (platform-specific)
- **`types/`**: Shared type definitions (screen elements, performance metrics)
- **`utils/`**: Common utilities (logging, file handling, downloads)
- **`assets/`**: Embedded assets (agent binaries for iOS devices)

## Key Patterns

### Device Interface
All devices implement `ControllableDevice` interface. When adding new device functionality:
1. Add method to the interface in `devices/common.go`
2. Implement in `android.go`, `ios.go`, and `simulator.go`
3. Add command in `commands/` layer
4. Add CLI flag handling in `cli/` layer
5. Add RPC handler and register in `server/dispatch.go`

### RPC Method Naming
Methods follow the pattern: `category.action` (e.g., `device.screenshot`, `device.apps.launch`)

### Verbose Logging
Use `utils.Verbose()` for debug output that respects the `--verbose` flag.

### Shutdown Hooks
Use `devices.ShutdownHook` for tracking resources that need cleanup on exit.

## Platform Notes

- **Android**: Requires `adb` in PATH
- **iOS Real Devices**: Requires WebDriverAgent installed on device
- **iOS Simulator**: Requires Xcode with iOS runtimes installed

## E2E Testing

E2E tests are in `test/` directory using Node.js/Mocha. They require iOS simulators with specific iOS runtimes. See `docs/TESTING.md` for details.
