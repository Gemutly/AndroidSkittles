#!/bin/bash
# Auto-crop screenshot by Android View ID
# Usage: ./capture_view.sh <output_file> [view_id]
# Example: ./capture_view.sh home.png id/skittle_container

set -euo pipefail

OUTPUT_FILE="${1:-screenshot_$(date +%Y%m%d_%H%M%S).png}"
VIEW_ID="$2"
TEMP_XML="/sdcard/ui_dump.xml"
TEMP_SCREEN="/sdcard/screen.png"

# Cleanup trap
cleanup() {
    adb shell rm -f "$TEMP_XML" "$TEMP_SCREEN" 2>/dev/null || true
    rm -f /tmp/ui_dump.xml /tmp/full_screen.png
}
trap cleanup EXIT

# Prerequisite checks
echo "Checking prerequisites..."
command -v magick >/dev/null 2>&1 || { echo "ERROR: ImageMagick not installed. Run: brew install imagemagick"; exit 1; }
command -v adb >/dev/null 2>&1 || { echo "ERROR: ADB not in PATH"; exit 1; }

if ! adb devices | grep -q "device$"; then
    echo "ERROR: No device connected. Check 'adb devices'"
    exit 1
fi

# If no view_id provided, capture full screen
if [ -z "$VIEW_ID" ]; then
    echo "No View ID provided. Capturing full screen..."
    adb shell screencap -p "$TEMP_SCREEN"
    adb pull "$TEMP_SCREEN" "$OUTPUT_FILE"
    echo "Full screenshot saved to: $OUTPUT_FILE"
    exit 0
fi

# Dump UI hierarchy
echo "Dumping UI hierarchy..."
adb shell uiautomator dump "$TEMP_XML"
adb pull "$TEMP_XML" /tmp/ui_dump.xml

# Parse bounds for the view_id
BOUNDS=$(grep -o "resource-id=\"$VIEW_ID\"[^>]*bounds=\"\[[0-9,]*\]\[[0-9,]*\]" /tmp/ui_dump.xml 2>/dev/null | \
         grep -o "\[[0-9,]*\]\[[0-9,]*\]" | \
         sed 's/\]\[/,/;s/\[//g;s/\]//g' | head -1)

if [ -z "$BOUNDS" ]; then
    echo "ERROR: View ID '$VIEW_ID' not found in UI hierarchy"
    echo "Available resource-ids:"
    grep -o 'resource-id="[^"]*"' /tmp/ui_dump.xml | sort -u | head -20
    exit 1
fi

# Parse coordinates: left,top,right,bottom
LEFT=$(echo "$BOUNDS" | cut -d',' -f1)
TOP=$(echo "$BOUNDS" | cut -d',' -f2)
RIGHT=$(echo "$BOUNDS" | cut -d',' -f3)
BOTTOM=$(echo "$BOUNDS" | cut -d',' -f4)

WIDTH=$((RIGHT - LEFT))
HEIGHT=$((BOTTOM - TOP))

echo "Found view '$VIEW_ID' at bounds: ${WIDTH}x${HEIGHT}+${LEFT}+${TOP}"

# Capture and crop
echo "Capturing screenshot..."
adb shell screencap -p "$TEMP_SCREEN"
adb pull "$TEMP_SCREEN" /tmp/full_screen.png

magick /tmp/full_screen.png -crop "${WIDTH}x${HEIGHT}+${LEFT}+${TOP}" "$OUTPUT_FILE"

# Validate output
if ! magick identify "$OUTPUT_FILE" &>/dev/null; then
    echo "ERROR: Output file is corrupt"
    exit 1
fi

echo "Cropped screenshot saved to: $OUTPUT_FILE"
