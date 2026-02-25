---
name: asset-genie
description: Generates Android app icons, adaptive icon XML, mipmap density variants, notification icons, and splash screen assets using ImageMagick. Triggered by "generate app icon", "create assets", or "prepare icons" requests.
---

# Asset-Genie Skill (Dynamic Icon & Asset Prep)

## When to use
- User asks to "generate app icon", "create icons", or "prepare assets"
- User has a high-res source image and needs all mipmap density variants
- User needs adaptive icon XML generated
- User wants notification icons, splash screens, or launcher icons created
- User wants AI-generated icon prompts based on the app's theme

## Platform constraint
This skill targets **Android primarily**. Asset generation uses **ImageMagick** (must be installed on the system). On Windows, use `magick` command; on Linux/macOS, use `convert` / `magick`.

## Prerequisites
- **ImageMagick** installed and available in PATH
  - Windows: `magick --version`
  - Linux/macOS: `convert --version` or `magick --version`
- A high-resolution source image (512x512 minimum, 1024x1024 recommended) or willingness to use AI generation prompts

## Instructions

### 1. Analyze the app's theme and branding
Scan the codebase for existing branding signals:

- **Colors:** Check `colors.xml`, `Color.kt`, `Theme.kt`, or equivalent for primary/accent colors
- **Existing icons:** Check `res/mipmap-*` and `res/drawable-*` for current app icons
- **App name:** From `AndroidManifest.xml` or `strings.xml`
- **App purpose:** From README, manifest, or codebase analysis

Summarize the brand context:
```markdown
## Brand Context
- **App Name:** [name]
- **Primary Color:** [hex]
- **Accent Color:** [hex]
- **App Category:** [utility/social/game/etc.]
- **Visual Style:** [minimal/bold/playful/professional]
```

### 2. Generate AI icon prompts
If no source icon exists, create `AssetGenie/ICON_PROMPTS.md` with 5 distinct prompt styles for AI image generators (DALL-E, Midjourney, Flux):

```markdown
# App Icon Generation Prompts
Use these with your preferred AI image generator. Request 1024x1024 PNG output.

## Prompt 1: Minimal Flat
> A flat-design app icon for a [app category] app called "[App Name]", [primary color] background, white geometric [symbol] centered, no text, rounded square shape, Google Material Design style, 1024x1024

## Prompt 2: Gradient Modern
> A modern app icon with a smooth gradient from [primary color] to [accent color], featuring a stylized [symbol] in white, subtle shadow, clean edges, app store ready, 1024x1024

## Prompt 3: 3D Rendered
> A 3D-rendered app icon for "[App Name]", glossy [primary color] surface, embossed [symbol], soft studio lighting, subtle reflection, professional mobile app icon, 1024x1024

## Prompt 4: Duotone Bold
> A bold duotone app icon using [primary color] and white, featuring an abstract [symbol], strong contrast, geometric shapes, modern tech aesthetic, 1024x1024

## Prompt 5: Outlined Clean
> A clean outlined app icon on [primary color] background, thin white line-art [symbol] centered, minimalist, elegant, contemporary mobile design, 1024x1024
```

Replace `[symbol]` with a contextually relevant icon concept based on the app's purpose.

### 3. Generate mipmap density variants
Once a source icon is available (user provides or selects from AI output), generate all required Android mipmap densities:

```bash
# Ensure source is 1024x1024 or larger
# Android mipmap density requirements:
#   mdpi:    48x48
#   hdpi:    72x72
#   xhdpi:   96x96
#   xxhdpi:  144x144
#   xxxhdpi: 192x192
#   Play Store: 512x512

# Windows (using magick)
magick <source>.png -resize 48x48 AssetGenie/mipmap-mdpi/ic_launcher.png
magick <source>.png -resize 72x72 AssetGenie/mipmap-hdpi/ic_launcher.png
magick <source>.png -resize 96x96 AssetGenie/mipmap-xhdpi/ic_launcher.png
magick <source>.png -resize 144x144 AssetGenie/mipmap-xxhdpi/ic_launcher.png
magick <source>.png -resize 192x192 AssetGenie/mipmap-xxxhdpi/ic_launcher.png
magick <source>.png -resize 512x512 AssetGenie/playstore-icon.png

# Linux/macOS (using convert or magick)
convert <source>.png -resize 48x48 AssetGenie/mipmap-mdpi/ic_launcher.png
convert <source>.png -resize 72x72 AssetGenie/mipmap-hdpi/ic_launcher.png
convert <source>.png -resize 96x96 AssetGenie/mipmap-xhdpi/ic_launcher.png
convert <source>.png -resize 144x144 AssetGenie/mipmap-xxhdpi/ic_launcher.png
convert <source>.png -resize 192x192 AssetGenie/mipmap-xxxhdpi/ic_launcher.png
convert <source>.png -resize 512x512 AssetGenie/playstore-icon.png
```

Create the output directories first:
```bash
# Windows
mkdir AssetGenie\mipmap-mdpi, AssetGenie\mipmap-hdpi, AssetGenie\mipmap-xhdpi, AssetGenie\mipmap-xxhdpi, AssetGenie\mipmap-xxxhdpi

# Linux/macOS
mkdir -p AssetGenie/mipmap-{mdpi,hdpi,xhdpi,xxhdpi,xxxhdpi}
```

### 4. Generate round icon variants
Android requires both square and round launcher icons:

```bash
# Generate round variants by applying a circular mask
# Windows
magick <source>.png -resize 48x48 ( +clone -threshold 100% -fill white -draw "circle 24,24 24,0" ) -channel-fx "| gray=>alpha" AssetGenie/mipmap-mdpi/ic_launcher_round.png
magick <source>.png -resize 72x72 ( +clone -threshold 100% -fill white -draw "circle 36,36 36,0" ) -channel-fx "| gray=>alpha" AssetGenie/mipmap-hdpi/ic_launcher_round.png
magick <source>.png -resize 96x96 ( +clone -threshold 100% -fill white -draw "circle 48,48 48,0" ) -channel-fx "| gray=>alpha" AssetGenie/mipmap-xhdpi/ic_launcher_round.png
magick <source>.png -resize 144x144 ( +clone -threshold 100% -fill white -draw "circle 72,72 72,0" ) -channel-fx "| gray=>alpha" AssetGenie/mipmap-xxhdpi/ic_launcher_round.png
magick <source>.png -resize 192x192 ( +clone -threshold 100% -fill white -draw "circle 96,96 96,0" ) -channel-fx "| gray=>alpha" AssetGenie/mipmap-xxxhdpi/ic_launcher_round.png
```

### 5. Generate Adaptive Icon XML
Create the Adaptive Icon resource files for Android 8.0+ (API 26+):

**`AssetGenie/adaptive-icon/ic_launcher.xml`:**
```xml
<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
    <background android:drawable="@color/ic_launcher_background"/>
    <foreground android:drawable="@mipmap/ic_launcher_foreground"/>
</adaptive-icon>
```

**`AssetGenie/adaptive-icon/ic_launcher_round.xml`:**
```xml
<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
    <background android:drawable="@color/ic_launcher_background"/>
    <foreground android:drawable="@mipmap/ic_launcher_foreground"/>
</adaptive-icon>
```

**`AssetGenie/adaptive-icon/ic_launcher_background.xml`** (using the app's primary color):
```xml
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <color name="ic_launcher_background">[PRIMARY_COLOR_HEX]</color>
</resources>
```

Also generate foreground images (the icon with safe zone padding — 66% of full size centered):

```bash
# Foreground needs to be 108x108dp (432x432px at xxxhdpi) with the icon in the center 72x72dp area
magick <source>.png -resize 288x288 -gravity center -background none -extent 432x432 AssetGenie/adaptive-icon/ic_launcher_foreground.png
```

### 6. Generate notification icon (optional)
Notification icons must be white-on-transparent, simple silhouette:

```bash
# Convert to white silhouette on transparent background
# Windows
magick <source>.png -resize 24x24 -colorspace gray -threshold 50% -negate -background none -alpha shape AssetGenie/notification/ic_notification_mdpi.png
magick <source>.png -resize 36x36 -colorspace gray -threshold 50% -negate -background none -alpha shape AssetGenie/notification/ic_notification_hdpi.png
magick <source>.png -resize 48x48 -colorspace gray -threshold 50% -negate -background none -alpha shape AssetGenie/notification/ic_notification_xhdpi.png
magick <source>.png -resize 72x72 -colorspace gray -threshold 50% -negate -background none -alpha shape AssetGenie/notification/ic_notification_xxhdpi.png
magick <source>.png -resize 96x96 -colorspace gray -threshold 50% -negate -background none -alpha shape AssetGenie/notification/ic_notification_xxxhdpi.png
```

### 7. Generate installation guide
Create `AssetGenie/INSTALL_GUIDE.md`:

```markdown
# Asset Installation Guide

## Launcher Icons
Copy each density folder's contents to your project:
- `mipmap-mdpi/` → `app/src/main/res/mipmap-mdpi/`
- `mipmap-hdpi/` → `app/src/main/res/mipmap-hdpi/`
- `mipmap-xhdpi/` → `app/src/main/res/mipmap-xhdpi/`
- `mipmap-xxhdpi/` → `app/src/main/res/mipmap-xxhdpi/`
- `mipmap-xxxhdpi/` → `app/src/main/res/mipmap-xxxhdpi/`

## Adaptive Icon (Android 8.0+)
Copy `adaptive-icon/` XML files to `app/src/main/res/mipmap-anydpi-v26/`

## Notification Icon
Copy notification icons to corresponding `app/src/main/res/drawable-*dpi/` directories

## Play Store
Upload `playstore-icon.png` (512x512) as the Hi-res icon in Google Play Console
```

### 8. Final output
All deliverables reside in the `AssetGenie/` directory:
```
AssetGenie/
├── ICON_PROMPTS.md (if no source icon was provided)
├── INSTALL_GUIDE.md
├── playstore-icon.png
├── mipmap-mdpi/
│   ├── ic_launcher.png
│   └── ic_launcher_round.png
├── mipmap-hdpi/
│   └── ...
├── mipmap-xhdpi/
│   └── ...
├── mipmap-xxhdpi/
│   └── ...
├── mipmap-xxxhdpi/
│   └── ...
├── adaptive-icon/
│   ├── ic_launcher.xml
│   ├── ic_launcher_round.xml
│   ├── ic_launcher_background.xml
│   └── ic_launcher_foreground.png
└── notification/
    ├── ic_notification_mdpi.png
    └── ...
```

## Important notes
- **ImageMagick is required.** If it is not installed, inform the user and provide installation instructions (`winget install ImageMagick`, `brew install imagemagick`, or `apt install imagemagick`).
- Source images should be 1024x1024 or larger for best quality. Warn if the source is smaller than 512x512.
- Always use PNG format for icons (lossless, supports transparency).
- Adaptive icon foreground must include safe zone padding — the visible area is only ~66% of the full canvas.
- Notification icons **must** be white silhouette on transparent background per Android guidelines.
- On Windows, use `magick` (not `convert`, which conflicts with a Windows system utility).
- If the project already has icons in `res/mipmap-*`, note that these will be overwritten and suggest backing them up first.
