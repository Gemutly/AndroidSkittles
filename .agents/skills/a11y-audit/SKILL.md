---
name: a11y-audit
description: Audits Android UI for WCAG accessibility compliance — checks for missing contentDescription, analyzes color contrast ratios, validates touch target sizes, and generates A11Y_FIXES.md with code snippets. Triggered by "accessibility audit", "a11y check", or "WCAG compliance" requests.
---

# Accessibility Advocate Skill (a11y-audit)

## When to use
- User asks for an "accessibility audit", "a11y check", or "WCAG compliance" review
- User wants to ensure the app meets accessibility standards before release
- User needs to find missing `contentDescription` attributes
- User wants color contrast analysis of their UI
- A `UXReview/UX_INVENTORY.md` exists and the user wants it audited for accessibility

## Platform constraint
This skill targets **Android and cross-platform mobile projects** (Jetpack Compose, XML layouts, React Native/Expo, Flutter). Do not scan for or reference iOS-native accessibility APIs (UIAccessibility).

## Instructions

### 0. Check for existing codebase map (optimization)
If `CodebaseMap/CODEBASE_MAP.md` exists, read the **Feature Index** to locate:
- UI component directories (XML layouts, Compose files, React Native screens)
- Design system color definitions (for contrast analysis)
- Project structure and technology stack

This accelerates the audit. If the map doesn't exist, proceed with full codebase scanning.

### 1. Scan for missing content descriptions
Search all UI files for interactive and informational elements that lack accessibility labels:

**Android XML layouts:**
```bash
# Find ImageView, ImageButton, and custom views missing contentDescription
grep -rnl "ImageView\|ImageButton\|FloatingActionButton" --include="*.xml" app/src/main/res/layout/ | xargs grep -L "contentDescription"
```

Check for:
- `ImageView` without `android:contentDescription`
- `ImageButton` without `android:contentDescription`
- `FloatingActionButton` without `android:contentDescription`
- `CheckBox`, `Switch`, `ToggleButton` without meaningful labels
- Decorative images that should have `android:importantForAccessibility="no"`

**Jetpack Compose:**
```bash
# Find Image/Icon composables missing contentDescription or semantics
grep -rn "Image(\|Icon(" --include="*.kt" app/src/main/
```

Check for:
- `Image()` without `contentDescription` parameter (or set to `null` without `semantics { }` block)
- `Icon()` without `contentDescription`
- `IconButton()` without `contentDescription` on the inner `Icon`
- Clickable elements without `Modifier.semantics { }` labels
- `Modifier.clickable` without `onClickLabel`

**React Native / Expo:**
```bash
# Find Image/TouchableOpacity missing accessibility props
grep -rn "Image\|TouchableOpacity\|Pressable" --include="*.tsx" --include="*.jsx" src/
```

Check for:
- `Image` without `accessibilityLabel`
- `TouchableOpacity` / `Pressable` without `accessibilityLabel` and `accessibilityRole`
- Missing `accessibilityHint` on complex interactions

**Flutter:**
```bash
# Find Image/Icon widgets missing semantics
grep -rn "Image\.\|Icon(" --include="*.dart" lib/
```

Check for:
- `Image` without `Semantics` wrapper or `semanticLabel`
- `Icon` without `semanticLabel`
- `GestureDetector` without `Semantics` wrapper

### 2. Analyze color contrast ratios
Check the design system and UI files for WCAG 2.1 contrast compliance:

**Extract color pairs from the codebase:**
- Text color vs. background color combinations
- Button text vs. button background
- Placeholder/hint text vs. input background
- Icon color vs. surrounding background

**WCAG 2.1 minimum contrast ratios:**
- **AA Normal text (< 18sp):** 4.5:1
- **AA Large text (≥ 18sp or ≥ 14sp bold):** 3:1
- **AAA Normal text:** 7:1
- **AAA Large text:** 4.5:1
- **Non-text UI components & graphics:** 3:1

Calculate contrast ratio using the relative luminance formula:
```
Contrast Ratio = (L1 + 0.05) / (L2 + 0.05)
where L1 is the lighter color's relative luminance and L2 is the darker
```

For each identified color pair, report:
- The two colors (hex values)
- The calculated contrast ratio
- Pass/fail for AA and AAA at the relevant text size
- Where in the codebase this combination is used

If screenshots exist in `UXReview/screenshots/`, visually analyze them for obvious low-contrast areas.

### 3. Validate touch target sizes
WCAG 2.5.8 and Android guidelines require minimum touch target sizes:

- **Minimum:** 48x48dp (Android) / 44x44pt (general WCAG)
- Check for buttons, icons, and interactive elements smaller than this

**Android XML:**
```bash
# Find elements with explicit small dimensions
grep -rn "layout_width\|layout_height\|minWidth\|minHeight" --include="*.xml" app/src/main/res/layout/ | grep -E "[0-9]+(dp|dip)" 
```

**Jetpack Compose:**
```bash
# Find explicit small size modifiers
grep -rn "\.size(\|\.width(\|\.height(" --include="*.kt" app/src/main/ | grep -E "[0-9]+"
```

Flag any interactive element with dimensions below 48dp.

### 4. Check navigation and focus order
Verify that screen readers can navigate the app logically:

- Check for `android:focusable`, `android:nextFocusDown`, `android:nextFocusForward` attributes
- In Compose, check for `Modifier.focusOrder` or `Modifier.focusProperties`
- Verify that `TabRow` / `BottomNavigation` items are properly labeled
- Check that dialog/modal content traps focus correctly
- Ensure headings use proper heading semantics (`accessibilityHeading = true` or `heading()`)

### 5. Audit for additional a11y concerns
Check for:

- **Text scaling:** Ensure text uses `sp` units (not `dp` or `px`) so it respects system font size
- **Color-only information:** Flag UI where color is the only indicator of state (e.g., red for error without an icon or text)
- **Animation:** Check for `prefers-reduced-motion` / `AccessibilityManager.isEnabled` handling
- **Live regions:** Verify dynamic content updates use `android:accessibilityLiveRegion` or `LiveRegionMode` in Compose
- **Custom views:** Ensure custom-drawn views implement `AccessibilityDelegate` or `semantics { }` blocks

### 6. Cross-reference with UX_INVENTORY.md (if available)
If `UXReview/UX_INVENTORY.md` exists:

1. Read each screen entry
2. For each screen, run the above checks against its source file
3. If screenshots exist, note any visible contrast or sizing concerns
4. Add accessibility findings as annotations to the inventory

### 7. Generate A11Y_FIXES.md
Create `A11yAudit/A11Y_FIXES.md` with categorized, actionable findings:

```markdown
# Accessibility Audit Report
Generated: [date]
Standard: WCAG 2.1 Level AA

## Summary
- **Critical issues:** [count] (must fix)
- **Major issues:** [count] (should fix)
- **Minor issues:** [count] (nice to fix)
- **Pass rate:** [percentage of elements passing]

## Critical: Missing Content Descriptions
### [File Path]:[Line Number]
- **Element:** `ImageButton` / `Icon` / etc.
- **Issue:** No contentDescription — screen readers will skip or announce "unlabeled"
- **Fix:**
  ```kotlin
  // Before
  Icon(Icons.Default.Search, contentDescription = null)
  // After
  Icon(Icons.Default.Search, contentDescription = "Search")
  ```

## Critical: Insufficient Color Contrast
### [File Path] — [Element Description]
- **Foreground:** #AAAAAA
- **Background:** #FFFFFF
- **Ratio:** 2.3:1 (requires 4.5:1 for AA)
- **Fix:** Change foreground to #767676 or darker (4.5:1 minimum)

## Major: Touch Targets Below 48dp
### [File Path]:[Line Number]
- **Element:** Icon button
- **Current size:** 32x32dp
- **Fix:**
  ```kotlin
  // Add minimum touch target
  Modifier.sizeIn(minWidth = 48.dp, minHeight = 48.dp)
  ```

## Major: Text Using dp Instead of sp
### [File Path]:[Line Number]
- **Issue:** `fontSize = 14.dp` — will not scale with system font size
- **Fix:** `fontSize = 14.sp`

## Minor: Missing Live Regions
### [File Path]:[Line Number]
- **Element:** Dynamic status text
- **Fix:** Add `Modifier.semantics { liveRegion = LiveRegionMode.Polite }`

## Color Contrast Matrix
[List all color pairs with their ratios and pass/fail status]
```

### 8. Generate fix summary with priority
Create `A11yAudit/FIX_PRIORITY.md`:

```markdown
# Accessibility Fix Priority Guide

## P0 — Ship Blockers (fix before release)
- [ ] Add contentDescription to all interactive ImageViews/Icons
- [ ] Fix color contrast below 3:1 ratio

## P1 — High Priority (fix in next sprint)
- [ ] Increase touch targets below 48dp
- [ ] Convert dp text sizes to sp
- [ ] Add accessibility labels to unlabeled buttons

## P2 — Medium Priority (scheduled improvement)
- [ ] Add live regions to dynamic content
- [ ] Improve focus order in complex layouts
- [ ] Add accessibilityHint to non-obvious interactions

## P3 — Nice to Have
- [ ] Add reduced-motion support
- [ ] Add AAA contrast compliance (7:1)
```

### 9. Final output
All deliverables reside in the `A11yAudit/` directory:
```
A11yAudit/
├── A11Y_FIXES.md
└── FIX_PRIORITY.md
```

## Important notes
- **WCAG 2.1 Level AA** is the default target standard. If the user requests AAA, adjust thresholds accordingly.
- This audit is static analysis only — it cannot test with a real screen reader. Recommend the user also test with TalkBack (Android) for full coverage.
- Color contrast calculations require hex/RGB values. If colors are defined as theme references, resolve them to actual values first.
- Skip iOS-only files (`.swift`, `.storyboard`, `.xib`).
- If neither XML layouts nor Compose files are found, check for React Native/Expo or Flutter equivalents before reporting "no UI files found."
- Decorative images (purely visual, no informational content) should be marked as `importantForAccessibility="no"` rather than given a contentDescription.
- This skill complements the `ux-review` skill — run `ux-review` first to generate the `UX_INVENTORY.md` that this skill can audit.
