---
name: play-store-screenshots
description: Automates capturing app screenshots via ADB for Play Store submission.
---

# Play Store Screenshots Skill

Use this skill when the user wants to generate visual assets for the Play Store, verify UI layouts, or automate screenshot capture.

## Prerequisites

Before using this skill, ensure:
- ImageMagick installed: `brew install imagemagick` (or `choco install imagemagick` on Windows)
- ADB in PATH (part of Android Platform Tools)
- Device/emulator connected with USB debugging enabled

Run the diagnostics script to verify all prerequisites.

## Scripts

| Script | Purpose | Usage |
|--------|---------|-------|
| `capture_view.sh` | Capture full screen or crop by View ID | `./capture_view.sh output.png [view_id]` |
| `diagnostics.sh` | Verify environment is ready | `./diagnostics.sh` |

## Quick Commands

| Task | Command |
|------|---------|
| Run diagnostics | `.claude/skills/play-store-screenshots/scripts/diagnostics.sh` |
| Full screen capture | `.claude/skills/play-store-screenshots/scripts/capture_view.sh home.png` |
| Capture specific View | `.claude/skills/play-store-screenshots/scripts/capture_view.sh grid.png id/skittle_container` |

## Workflow

1. **Verify Environment:** Run diagnostics script
2. **Verify Device:** Run `adb devices`
3. **Navigate App:** Use `adb shell input tap x y` or semantic navigation
4. **Capture:** Run appropriate capture script
5. **Validate:** Check output file exists and is valid

## Semantic Navigation (Decision Point)

When asked to navigate, DO NOT guess coordinates. Instead:

1. Run `adb shell uiautomator dump /sdcard/ui.xml && adb pull /sdcard/ui.xml`
2. Parse the XML to find the target element:
   - Search by `text` attribute (e.g., `text="Settings"`)
   - Search by `content-desc` (accessibility label)
   - Search by `resource-id`
3. Extract the `bounds` attribute: `bounds="[left,top][right,bottom]"`
4. Calculate center: `x = (left + right) / 2`, `y = (top + bottom) / 2`
5. Tap: `adb shell input tap x y`

## Pre-Capture Verification (Decision Point)

Before capturing, verify the app is in the expected state:

1. Dump current UI hierarchy
2. Check for expected elements (e.g., wait for `recycler_view` to appear)
3. Check for error dialogs: `grep "AlertDialog" ui.xml`
4. If state is wrong, abort and report rather than capture a bad screenshot

## Smart Waiting (Decision Point)

Do not use static `sleep` commands. Poll for conditions:

```bash
# Wait for element to appear (max 10 seconds)
for i in {1..20}; do
    adb shell uiautomator dump /sdcard/ui.xml
    adb pull /sdcard/ui.xml /tmp/ui.xml
    if grep -q "target_view_id" /tmp/ui.xml; then
        echo "Element found"
        break
    fi
    sleep 0.5
done
```

## User Simulation Mode

When asked to "simulate a user journey" or "capture Play Store screenshots":

### Persona
You are a QA engineer generating assets for a Play Store listing. Your goal is to showcase the app's best features.

### Standard Journey Protocol

1. **Launch App:**
   ```bash
   adb shell am start -n com.yourpackage/.MainActivity
   ```

2. **Home Screen:**
   - Wait for main UI to load
   - **Decision:** If screen has content, capture as `01_home.png`

3. **Feature Navigation:**
   - Identify primary action buttons via UI dump
   - Tap first major feature
   - Wait for transition
   - Capture as `02_feature_detail.png`

4. **Settings/Profile:**
   - Navigate to Settings (search for gear icon or "Settings" text)
   - Capture as `03_settings.png`

### Output Format

Return a summary:
```
OK Captured: 01_home.png, 02_detail.png, 03_settings.png
FAIL Failed: Search screen (button not found)
```

## Integration with mobilecli

This project has built-in device control capabilities. When available, prefer using mobilecli commands:

| mobilecli Command | ADB Equivalent |
|-------------------|----------------|
| `mobilecli screenshot --device <id>` | `adb shell screencap -p` |
| `mobilecli io tap --device <id> x,y` | `adb shell input tap x y` |
| `mobilecli dump ui --device <id>` | `adb shell uiautomator dump` |

## Guidelines

- Always verify device connection with `adb devices` first
- Always run diagnostics on first use
- If a View ID is not found, the script will list available IDs
- Use semantic navigation instead of hardcoded coordinates
- Validate each screenshot before reporting success
