---
name: play-store-prep
description: Prepares Android Play Store assets including submission text, screenshots, and AI mockup generation prompts. Triggered by "prepare for Play Store" or "store submission" requests.
---

# Play Store Prep Skill

## When to use
- User asks to "prepare for Play Store" or "store submission"
- User needs to generate store listing assets and screenshots
- User wants mockup prompts for marketing images

## Platform constraint
This skill targets **Android only**. All audits, assets, and prompts are for Google Play Store submission.

## Instructions

### 1. Audit the codebase for Play Store requirements
Check for the following and report any issues:

- **AndroidManifest.xml:** Verify `package`, `versionCode`, `versionName`, required permissions, `<uses-feature>` declarations
- **Signing config:** Check for release keystore configuration in `build.gradle` / `build.gradle.kts`
- **App icons:** Verify adaptive icon resources exist (`mipmap-*` directories) and a 512x512 hi-res icon is available
- **Permissions:** Flag any dangerous or unnecessary permissions that may trigger Play Store review
- **Target SDK:** Ensure `targetSdkVersion` meets current Play Store minimum requirements
- **ProGuard / R8:** Check if code shrinking and obfuscation are configured for release builds

Create `PlayStorePrep/AUDIT_REPORT.md` summarizing findings with pass/fail status for each item.

### 2. Draft store listing text
Create `PlayStorePrep/LISTING.md` with:

```markdown
# Play Store Listing

## App Title (max 30 characters)
[Generated title]

## Short Description (max 80 characters)
[Generated short description]

## Full Description (max 4000 characters)
[Generated full description with feature highlights, formatted for Play Store]

## Category
[Suggested category]

## Tags / Keywords
[Comma-separated list of relevant keywords]

## Content Rating
[Suggested rating based on app content and permissions]

## Privacy Policy URL
[Placeholder or detected URL]

## Contact Email
[Placeholder]
```

Base the listing content on the project README, feature set, and codebase analysis.

### 3. Capture raw Android screenshots
Use `mobilecli` or `adb` to capture screenshots on a connected device/emulator:

```bash
# Preferred: using mobilecli
mobilecli screenshot --device <device-id> --output PlayStorePrep/screenshots/raw/<screen_name>.png

# Fallback: using adb
adb -s <device-id> exec-out screencap -p > PlayStorePrep/screenshots/raw/<screen_name>.png
```

Capture at least the following screens (if applicable):
- Home / main screen
- Key feature screens (2-3)
- Settings or profile screen
- Any onboarding or tutorial screens

Play Store requirements:
- **Phone:** minimum 1080x1920 px (16:9) or 1440x2560
- **Tablet (optional):** 1200x1920 or higher
- Save as PNG format

### 4. Generate mockup prompts
Create `PlayStorePrep/MOCKUP_PROMPTS.md` with 5 distinct style groups for use in AI image generators (Midjourney, DALL-E, Stable Diffusion):

```markdown
# Mockup Generation Prompts
Use these prompts with your preferred AI image generator to create professional Play Store feature graphics and screenshots.

## Style 1: Urban Professional
> Photorealistic mockup of a modern smartphone held in a right hand at a slight angle, displaying [APP_SCREENSHOT], urban coffee shop background with warm ambient lighting, bokeh effect, 4K resolution, professional product photography

## Style 2: Nature Retreat
> Photorealistic mockup of a smartphone resting on a light wooden table, displaying [APP_SCREENSHOT], soft natural forest background with dappled sunlight, warm earth tones, lifestyle product photography, 4K

## Style 3: Futuristic Tech
> Photorealistic mockup of a floating smartphone with subtle shadow, displaying [APP_SCREENSHOT], dark gradient background with neon blue and purple accent lighting, tech-forward aesthetic, clean minimal composition, 4K

## Style 4: Minimal Studio
> Photorealistic mockup of a smartphone centered on a clean white surface, displaying [APP_SCREENSHOT], soft studio lighting with gentle shadows, minimalist background, high-end product catalog style, 4K

## Style 5: On-the-Go Lifestyle
> Photorealistic mockup of a smartphone held in a left hand, slightly rotated, displaying [APP_SCREENSHOT], blurred city street background at golden hour, candid lifestyle photography feel, vibrant warm tones, 4K
```

Replace `[APP_SCREENSHOT]` with the actual screenshot filename when generating.

Include hand position and rotation variations:
- Right hand, slight tilt (Styles 1, 3)
- Flat on surface (Styles 2, 4)
- Left hand, rotated (Style 5)

### 5. Organize final assets
Consolidate all files into the `PlayStorePrep/` directory:

```
PlayStorePrep/
├── AUDIT_REPORT.md
├── LISTING.md
├── MOCKUP_PROMPTS.md
└── screenshots/
    └── raw/
        ├── screen_home.png
        ├── screen_feature_1.png
        ├── screen_feature_2.png
        ├── screen_settings.png
        └── ...
```

## Important notes
- All outputs target Google Play Store specifically (not Apple App Store)
- Screenshots must meet Play Store dimension requirements (minimum 1080x1920 for phone)
- If no Android device/emulator is connected, generate all text assets and note that screenshots need to be captured separately
- The MOCKUP_PROMPTS.md is designed for external AI image generators — the agent does not generate images directly
- Store listing text should be professional, keyword-optimized, and based on actual app features
