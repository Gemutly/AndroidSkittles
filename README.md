# AndroidSkittles

A command-line tool for managing Android devices and emulators. Forked from [mobilecli](https://github.com/mobile-next/mobilecli) with a focus on Android-only support.

## Features

- **Device Management**: List, manage, and interact with connected Android devices and emulators
- **Emulator Control**: Boot and shutdown emulators programmatically
- **Screenshot Capture**: Take screenshots with PNG/JPEG format options
- **Screen Streaming**: Stream MJPEG/H.264 video directly from device
- **Device Control**: Reboot, tap, swipe, long-press, and hardware button input
- **App Management**: Launch, terminate, install, uninstall, and list apps
- **HTTP/WebSocket API**: JSON-RPC 2.0 interface for remote control

## Platform Support

| Platform | Supported |
|----------|:---------:|
| Android Real Device | Yes |
| Android Emulator | Yes |
| iOS Real Device | No (removed) |
| iOS Simulator | No (removed) |

## Installation

### Prerequisites

- **Android SDK** with `adb` in PATH
- **Go 1.25+** (for building from source)

### Build from Source

```bash
git clone https://github.com/Gemutly/AndroidSkittles.git
cd AndroidSkittles
make build
```

### Install Android Platform Tools

```bash
# macOS
brew install --cask android-platform-tools

# Windows (with Chocolatey)
choco install adb

# Linux
sudo apt install android-tools-adb
```

## Usage

All commands output JSON for easy parsing.

### List Devices

```bash
# List online devices
androidskittles devices

# Include offline emulators
androidskittles devices --include-offline
```

Example output:
```json
[
  {
    "id": "emulator-5554",
    "name": "Pixel 6",
    "platform": "android",
    "type": "emulator",
    "state": "online"
  }
]
```

### Screenshots

```bash
# PNG screenshot (default)
androidskittles screenshot --device <device-id>

# JPEG with quality
androidskittles screenshot --device <device-id> --format jpeg --quality 80

# Save to file
androidskittles screenshot --device <device-id> --output screen.png

# Output to stdout (for piping)
androidskittles screenshot --device <device-id> --output -
```

### Screen Streaming

```bash
androidskittles screencapture --device <device-id> --format mjpeg | ffplay -
```

### Device Control

```bash
# Boot emulator
androidskittles device boot --device <device-id>

# Shutdown emulator
androidskittles device shutdown --device <device-id>

# Reboot device
androidskittles device reboot --device <device-id>

# Tap at coordinates
androidskittles io tap --device <device-id> 100,200

# Long press
androidskittles io longpress --device <device-id> 100,200 --duration 2000

# Swipe
androidskittles io swipe --device <device-id> 100,500 100,200

# Hardware buttons
androidskittles io button --device <device-id> HOME
androidskittles io button --device <device-id> BACK
androidskittles io button --device <device-id> VOLUME_UP

# Send text
androidskittles io text --device <device-id> 'hello world'
```

### Supported Buttons

| Button | Description |
|--------|-------------|
| `HOME` | Home button |
| `BACK` | Back button |
| `POWER` | Power button |
| `VOLUME_UP` | Volume up |
| `VOLUME_DOWN` | Volume down |
| `DPAD_UP/DOWN/LEFT/RIGHT/CENTER` | D-pad controls |

### App Management

```bash
# List installed apps
androidskittles apps list --device <device-id>

# Get foreground app
androidskittles apps foreground --device <device-id>

# Launch app
androidskittles apps launch com.example.app --device <device-id>

# Terminate app
androidskittles apps terminate com.example.app --device <device-id>

# Install APK
androidskittles apps install /path/to/app.apk --device <device-id>

# Uninstall app
androidskittles apps uninstall com.example.app --device <device-id>
```

## HTTP API

Start the server for HTTP/WebSocket access:

```bash
androidskittles server start --port 12000
```

### JSON-RPC 2.0 Examples

```bash
# List devices
curl http://localhost:12000/rpc -X POST -d '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "devices.list",
  "params": {}
}'

# Take screenshot
curl http://localhost:12000/rpc -X POST -d '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "device.screenshot",
  "params": {"deviceId": "emulator-5554"}
}'

# Tap screen
curl http://localhost:12000/rpc -X POST -d '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "device.io.tap",
  "params": {"deviceId": "emulator-5554", "x": 100, "y": 200}
}'
```

### Available RPC Methods

| Method | Description |
|--------|-------------|
| `devices.list` | List connected devices |
| `device.screenshot` | Capture screenshot |
| `device.info` | Get device info |
| `device.boot` | Boot emulator |
| `device.shutdown` | Shutdown emulator |
| `device.reboot` | Reboot device |
| `device.io.tap` | Tap coordinates |
| `device.io.longpress` | Long press |
| `device.io.swipe` | Swipe gesture |
| `device.io.button` | Press button |
| `device.io.text` | Send text |
| `device.apps.list` | List apps |
| `device.apps.launch` | Launch app |
| `device.apps.terminate` | Terminate app |
| `device.apps.install` | Install app |
| `device.apps.uninstall` | Uninstall app |

## WebSocket Support

Connect via WebSocket for persistent connections:

```bash
wscat -c ws://localhost:12000/ws

> {"jsonrpc":"2.0","id":1,"method":"devices.list","params":{}}
< {"jsonrpc":"2.0","id":1,"result":[...]}
```

## Development

### Build & Test

```bash
make build       # Build binary
make test        # Run unit tests
make test-cover  # Run with coverage
make lint        # Run linter
make fmt         # Format code
```

### Run Single Test

```bash
go test -v -run TestName ./path/to/package
```

### Project Structure

```
androidskittles/
├── cli/           # Cobra CLI commands
├── commands/      # Business logic
├── devices/       # Device abstraction layer
├── server/        # HTTP/WebSocket server
├── daemon/        # Background process support
├── types/         # Type definitions
├── utils/         # Utilities
└── main.go        # Entry point
```

## License

AGPL v3.0 - See [LICENSE](LICENSE) for details.

## Acknowledgments

Forked from [mobilecli](https://github.com/mobile-next/mobilecli) by Mobile Next.
