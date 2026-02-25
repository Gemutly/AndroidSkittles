---
name: adb-sleuth
description: Automated Android debugger that monitors logcat, diagnoses crashes and errors, cross-references source files, and runs UI sanity checks. Triggered by "debug app", "check logs", "why is my app crashing", or "sanity check" requests.
---

# Adb-Sleuth Skill (Debugger & Tester)

## When to use
- User reports an app crash or runtime error on Android
- User asks to "check logs", "debug", or "why is my app crashing"
- User wants to monitor logcat for a specific package
- User requests a "sanity check" or "smoke test" of all screens
- User needs help interpreting a stack trace or ANR

## Platform constraint
This skill targets **Android only** (physical devices and emulators). Requires `adb` in PATH or a built `mobilecli` binary.

## Instructions

### 0. Check for existing codebase map (optimization)
If `CodebaseMap/CODEBASE_MAP.md` exists, read it to understand:
- Project structure and key directories
- Where to find AndroidManifest.xml or build.gradle for package name
- Location of source files for stack trace cross-referencing

This provides quick context. If the map doesn't exist, proceed with manual discovery.

### 1. Identify the target device and package
Determine the connected device and the app's package name:

```bash
# List connected devices
mobilecli devices

# Get the foreground app's package name (if the app is running)
mobilecli apps foreground --device <device-id>

# Alternatively, check AndroidManifest.xml or build.gradle for the package name
```

If no device is connected, inform the user and stop.

### 2. Pull and filter logcat output
Capture the last 500 lines of logcat filtered to the target package:

```bash
# Using adb (primary method)
adb -s <device-id> logcat -d -t 500 --pid=$(adb -s <device-id> shell pidof -s <package-name>) 2>/dev/null

# Broader filter if PID is unavailable (app has crashed)
adb -s <device-id> logcat -d -t 500 | grep -iE "<package-name>|AndroidRuntime|FATAL|ANR|Exception|Error"
```

If the user reports a crash that just happened, focus on the last 50 lines containing:
- `FATAL EXCEPTION`
- `AndroidRuntime`
- `java.lang.` or `kotlin.` exception classes
- `ANR in`
- `Process: <package-name>`

### 3. Parse and classify the error
Analyze the captured logs and classify the issue:

- **Crash (FATAL EXCEPTION):** Extract the full stack trace, root exception class, and message
- **ANR (Application Not Responding):** Identify the blocked thread and what it was waiting on
- **Non-fatal error:** Extract warning/error log lines and their frequency
- **OOM (OutOfMemoryError):** Note memory allocation that failed and suggest heap analysis
- **SecurityException:** Flag permission issues and the missing permission

Create a structured diagnosis:

```markdown
## Diagnosis
- **Type:** Crash / ANR / Non-fatal / OOM / Security
- **Exception:** [full exception class name]
- **Message:** [exception message]
- **Root Cause Line:** [file:line from stack trace that's in the user's source code]
- **Frequency:** [one-time / recurring]
```

### 4. Cross-reference with source code
From the stack trace, identify lines that belong to the user's source code (not framework/library code):

1. Extract class names and line numbers from the stack trace
2. Search the codebase for matching source files using `grep` or file search
3. Read the relevant source file at the identified line numbers
4. Identify the offending code and its immediate context (10 lines above and below)

Present the findings:

```markdown
## Source Reference
- **File:** `app/src/main/java/com/example/MainActivity.kt`
- **Line:** 42
- **Code Context:**
  [relevant code snippet with the problematic line highlighted]
- **Analysis:** [explanation of why this line caused the error]
```

### 5. Suggest a fix
Based on the diagnosis and source analysis:

1. Explain the root cause in plain language
2. Provide a concrete code fix (diff format preferred)
3. If the fix involves adding a permission, show the manifest change too
4. Note any related issues the same pattern might cause elsewhere in the codebase

Save the full report to `AdbSleuth/DEBUG_REPORT.md`:

```markdown
# Debug Report — [Package Name]
Generated: [date]
Device: [device-id] ([device model])

## Error Summary
[one-line summary]

## Full Stack Trace
[formatted stack trace]

## Diagnosis
[structured diagnosis from Step 3]

## Source Reference
[source cross-reference from Step 4]

## Suggested Fix
[code fix with explanation]

## Additional Recommendations
[any broader concerns or patterns]
```

### 6. Run a UI sanity check (optional, on user request)
When the user asks for a "sanity check" or "smoke test":

1. Get the list of all activities/screens from the manifest or codebase scan
2. For each screen, attempt to navigate to it:

```bash
# Launch a specific activity
adb -s <device-id> shell am start -n <package-name>/.<ActivityName>

# Wait for the screen to render
sleep 2

# Capture a screenshot
mobilecli screenshot --device <device-id> --output AdbSleuth/sanity/<activity_name>.png

# Check logcat for errors during navigation
adb -s <device-id> logcat -d -t 20 | grep -iE "Exception|Error|FATAL|ANR"
```

3. Flag any screens that:
   - Crash on launch (black hole screens)
   - Show error states
   - Produce logcat warnings
   - Fail to render within 5 seconds

Save results to `AdbSleuth/SANITY_REPORT.md`:

```markdown
# UI Sanity Check Report
Generated: [date]

## Screen: [ActivityName]
- **Status:** ✅ OK / ❌ Crash / ⚠️ Warning
- **Screenshot:** ![ActivityName](sanity/activity_name.png)
- **Errors:** [any errors captured]
- **Load Time:** [approximate]
```

### 7. Final output
All deliverables reside in the `AdbSleuth/` directory:
```
AdbSleuth/
├── DEBUG_REPORT.md
├── SANITY_REPORT.md (if sanity check was run)
└── sanity/
    ├── main_activity.png
    ├── settings_activity.png
    └── ...
```

## Important notes
- Always filter logcat by package name or PID to avoid noise from system processes
- If `mobilecli` binary is not built yet, fall back to raw `adb` commands
- On Windows, replace `grep` with `Select-String` or use `findstr` for logcat filtering
- Never expose sensitive data (API keys, tokens, user data) that may appear in logs — redact them in reports
- The sanity check step is optional and should only run when explicitly requested (it can be slow on large apps)
- If the app is not installed or the package name is wrong, inform the user rather than proceeding blindly
