#!/bin/bash
# Validates the screenshot environment
set -e

echo "=== Screenshot Skill Diagnostics ==="
echo ""

PASS=0
FAIL=0

# Test 1: ADB
echo -n "1. ADB installed... "
if command -v adb &>/dev/null; then
    echo "OK"
    ((PASS++))
else
    echo "FAIL (Install Android Platform Tools)"
    ((FAIL++))
fi

# Test 2: Device connected
echo -n "2. Device connected... "
if adb devices 2>/dev/null | grep -q "device$"; then
    echo "OK"
    ((PASS++))
else
    echo "FAIL (Connect device/emulator)"
    ((FAIL++))
fi

# Test 3: ImageMagick
echo -n "3. ImageMagick installed... "
if command -v magick &>/dev/null; then
    echo "OK"
    ((PASS++))
else
    echo "FAIL (Run: brew install imagemagick)"
    ((FAIL++))
fi

# Test 4: UI Automator
echo -n "4. UI Automator working... "
if adb shell uiautomator dump /sdcard/test.xml 2>/dev/null; then
    adb shell rm /sdcard/test.xml 2>/dev/null
    echo "OK"
    ((PASS++))
else
    echo "FAIL"
    ((FAIL++))
fi

# Test 5: Screenshot capability
echo -n "5. Screenshot capability... "
if adb shell screencap -p /sdcard/test_cap.png 2>/dev/null; then
    adb shell rm /sdcard/test_cap.png 2>/dev/null
    echo "OK"
    ((PASS++))
else
    echo "FAIL"
    ((FAIL++))
fi

echo ""
echo "=== Results: $PASS/5 passed ==="

if [ $FAIL -gt 0 ]; then
    echo "Fix the issues above before using this skill."
    exit 1
else
    echo "All systems ready for screenshot capture."
    exit 0
fi
