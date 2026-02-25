---
name: ux-review
description: Audits UI components, generates a screen inventory, and captures Android screenshots into /UXReview. Triggered by "UX review" or "screen documentation" requests.
---

# UX Review Skill

## When to use
- User asks for a "UX review" or "screen documentation"
- User wants an inventory of all UI screens and their states
- User needs screenshots of the running app for design review

## Platform constraint
This skill targets **Windows, Linux, and Android only**. Do not scan for or reference iOS-specific UI files (SwiftUI, Storyboards, XIBs).

## Instructions

### 1. Analyze the codebase for UI screens
- Recursively scan source files for UI definitions:
  - **Android native:** XML layouts (`res/layout/*.xml`), Activity/Fragment classes (`*Activity.java`, `*Fragment.java`, `*Activity.kt`, `*Fragment.kt`)
  - **Jetpack Compose:** Files containing `@Composable` annotations
  - **React Native / Expo:** `.tsx` / `.jsx` files under `screens/`, `views/`, or `pages/` directories
  - **Flutter:** `.dart` files under `lib/screens/` or `lib/pages/`
- For each screen, identify:
  - Screen name and file path
  - Key UI components (buttons, inputs, lists, images, navigation elements)
  - Possible states (loading, empty, error, populated, authenticated/unauthenticated)

### 2. Create the UXReview output directory
```bash
mkdir -p UXReview/screenshots
```

### 3. Generate UX_INVENTORY.md
Create `UXReview/UX_INVENTORY.md` with the following structure:

```markdown
# UX Screen Inventory
Generated: [date]

## Screen: [ScreenName]
- **File:** `path/to/file`
- **Components:** [list of key UI components]
- **States:** [list of possible states]
- **Screenshot:** ![ScreenName](screenshots/screen_name.png)
- **Notes:** [any observations about UX patterns, accessibility, or issues]
```

Repeat for every discovered screen.

### 4. Capture screenshots (requires connected Android device or emulator)
Use `mobilecli` (this project's own CLI) or `adb` to capture screenshots:

```bash
# Using mobilecli (preferred)
mobilecli screenshot --device <device-id> --output UXReview/screenshots/<screen_name>.png

# Fallback: using adb directly
adb -s <device-id> exec-out screencap -p > UXReview/screenshots/<screen_name>.png
```

- If no device is connected, note this in the inventory and skip screenshot capture.
- Navigate the app to each screen before capturing (use `mobilecli io tap` or `adb shell input`).

### 5. Link screenshots in the inventory
Update each screen entry in `UX_INVENTORY.md` with relative image paths:
```markdown
- **Screenshot:** ![ScreenName](screenshots/screen_name.png)
```

### 6. Final output
All deliverables should reside in the `UXReview/` directory:
```
UXReview/
├── UX_INVENTORY.md
└── screenshots/
    ├── screen_home.png
    ├── screen_settings.png
    └── ...
```

## Important notes
- Skip any iOS-only files (`.swift`, `.storyboard`, `.xib`, `Podfile`)
- If the project uses a cross-platform framework (React Native, Flutter, Compose Multiplatform), document all screens but only capture on Android
- Use PNG format for screenshots (lossless quality for review)
- If `mobilecli` binary is not built yet, fall back to `adb` commands
