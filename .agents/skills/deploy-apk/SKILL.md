---
name: deploy-apk
description: Builds an Android project with Gradle, locates the output APK, and installs it on a connected device via ADB. Optionally launches the app after install. Triggered by "deploy app", "build and install", "run on device", "push APK", "install APK", or "package and deploy" requests.
---

# Deploy-APK Skill (Build → Push → Install)

## Platform constraint
Android only. Requires `adb` in PATH. Requires a Gradle-based Android project (presence of `gradlew` or `gradlew.bat`).

## Instructions

### 1. Locate the Android project root
Find the Gradle wrapper and project files:

```bash
# Look for gradlew / gradlew.bat in the current directory or parent directories
# Key markers: gradlew, gradlew.bat, build.gradle, build.gradle.kts, settings.gradle, settings.gradle.kts
```

If the user provides an explicit project path, use that. Otherwise search the current working directory and its parents.
If no Gradle project is found, inform the user and stop.

### 2. Detect build variant
Check `app/build.gradle` or `app/build.gradle.kts` for available build variants:
- Default to `assembleDebug` unless the user requests a release or specific flavor.
- If the user says "release", use `assembleRelease` (warn that signing config must be set up).
- For flavor builds, use `assemble<Flavor><BuildType>` (e.g., `assembleFreeDebug`).

### 3. Build the APK

```powershell
# Windows
& "<project-root>\gradlew.bat" -p "<project-root>" assembleDebug
```

```bash
# Linux/macOS
chmod +x <project-root>/gradlew
<project-root>/gradlew -p "<project-root>" assembleDebug
```

If the build fails:
- Show the error output to the user.
- Check for common issues: missing SDK, wrong Java version, dependency resolution failures.
- Do NOT retry automatically — let the user fix the issue first.

### 4. Locate the output APK
Search for the built APK in the standard Gradle output directories:

```bash
# Typical locations (debug variant):
#   app/build/outputs/apk/debug/app-debug.apk
#   app/build/outputs/apk/<flavor>/debug/app-<flavor>-debug.apk
#
# For release:
#   app/build/outputs/apk/release/app-release.apk
```

Use file search to find `.apk` files under `app/build/outputs/apk/`. Pick the most recently modified APK matching the requested build variant.
If multiple APKs are found, prefer the one matching the build command (e.g., `*-debug.apk` for `assembleDebug`).

### 5. Detect connected device

```bash
adb devices
```

- If exactly one device is connected, auto-select it.
- If multiple devices are connected and the user did not specify one, list them and ask which to target.
- If no devices are connected, inform the user and stop. Suggest checking USB debugging or running an emulator.

### 6. Install the APK

```bash
# Fresh install or update (replace existing)
adb -s <device-id> install -r "<path-to-apk>"
```

If `mobilecli` / `androidskittles` binary is built and available, prefer:
```bash
androidskittles apps install "<path-to-apk>" --device <device-id>
```

Handle common install errors:
- **INSTALL_FAILED_UPDATE_INCOMPATIBLE**: Signatures differ. Ask user if they want to uninstall first (`adb -s <device-id> uninstall <package-name>`), then retry install.
- **INSTALL_FAILED_INSUFFICIENT_STORAGE**: Inform user device storage is full.
- **INSTALL_FAILED_OLDER_SDK**: APK minSdk is higher than device API level.

### 7. Launch the app (optional, default: yes)
After successful install, launch the app unless the user said "install only" or "don't launch":

```bash
# Extract package name and launcher activity from the APK
adb -s <device-id> shell cmd package resolve-activity --brief -c android.intent.category.LAUNCHER <package-name> | tail -n1

# Launch it
adb -s <device-id> shell am start -n <package-name>/<launcher-activity>
```

Alternatively, if the package name is known from `build.gradle`:
```bash
adb -s <device-id> shell monkey -p <package-name> -c android.intent.category.LAUNCHER 1
```

### 8. Report result
Provide a concise summary:

```
✅ Built: app-debug.apk (2.4 MB)
✅ Installed on: Pixel 6 (emulator-5554)
✅ Launched: com.example.myapp
```

If any step failed, report which step and the error.

## Important notes
- On Windows, use `gradlew.bat`; on Linux/macOS use `./gradlew` (ensure it is executable).
- Always use `install -r` to allow reinstall over existing app without uninstalling first.
- Never expose signing keys, keystore passwords, or other secrets from `build.gradle` or `local.properties`.
- If the project uses multiple app modules (not just `app/`), check `settings.gradle` for module names and ask the user which to build.
- For large projects, `assembleDebug` can take several minutes — inform the user that the build is running.
- If `JAVA_HOME` is not set or the wrong Java version is detected, suggest setting it before building.
