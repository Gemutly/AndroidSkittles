---
name: design-bridge
description: Maps a reference image or Figma layout to the codebase's existing design system and drafts implementation using project components. Triggered by "design to code", "match this design", or "implement this screen" requests.
---

# Design-to-Code Bridge Skill

## When to use
- User provides a reference image, screenshot, or Figma layout to implement
- User asks to "match this design" or "implement this screen"
- User wants to build a UI screen using existing project components
- User wants to prevent style duplication by reusing the design system

## Platform constraint
This skill targets **Android and cross-platform mobile projects** (Jetpack Compose, React Native/Expo, Flutter). Do not scan for or reference iOS-native UI files (SwiftUI, Storyboards, XIBs).

## Instructions

### 0. Check for existing codebase map (optimization)
If `CodebaseMap/CODEBASE_MAP.md` exists, read the **Feature Index** section to quickly locate:
- UI component files and directories
- Theme/color/typography definitions
- Existing design system structure

This saves time scanning the entire codebase. If the map doesn't exist or is outdated, proceed with manual scanning.

### 1. Catalog the existing design system
Scan the codebase to build a map of all reusable UI primitives:

- **Jetpack Compose:**
  - Theme files (`Theme.kt`, `Color.kt`, `Type.kt`, `Shape.kt`) in any `ui/theme/` directory
  - Reusable `@Composable` functions in `ui/components/`, `components/`, or `common/` directories
  - Material theme overrides (`MaterialTheme`, `Typography`, `ColorScheme`)
- **Android XML:**
  - Style definitions in `res/values/styles.xml`, `res/values/themes.xml`
  - Custom components in `res/layout/` and custom view classes
  - Color and dimension resources (`colors.xml`, `dimens.xml`)
- **React Native / Expo:**
  - Theme/style objects in `theme/`, `styles/`, or `constants/` directories
  - Shared components in `components/` or `ui/` directories
  - Styled-components or StyleSheet definitions
- **Flutter:**
  - `ThemeData` definitions in `lib/theme/` or `lib/config/`
  - Reusable widgets in `lib/widgets/` or `lib/components/`

Create a temporary internal index of:
- Component name, file path, and visual purpose (e.g., "PrimaryButton", "CardContainer", "HeaderText")
- Color tokens and their hex values
- Typography scale (font sizes, weights, families)
- Spacing / padding constants

### 2. Analyze the reference design
Examine the provided reference image or Figma layout and identify:

- **Layout structure:** Column/row hierarchy, nesting, scroll behavior
- **Components:** Buttons, text fields, cards, lists, images, icons, navigation bars
- **Typography:** Heading levels, body text, captions — map to existing type scale
- **Colors:** Background, surface, primary, secondary, accent — map to existing color tokens
- **Spacing:** Margins, padding, gaps — map to existing spacing constants
- **States:** Default, pressed, disabled, error, loading
- **Interactions:** Tap targets, swipe areas, scroll indicators

### 3. Match design elements to existing components
For each visual element in the reference:

1. Search the component index from Step 1 for the closest match
2. If an exact match exists → use it directly
3. If a partial match exists → use it with minor parameter overrides (color, size, text)
4. If no match exists → flag it as a **new component candidate** and note what's missing

Create `DesignBridge/COMPONENT_MAP.md`:

```markdown
# Component Mapping — [Screen Name]
Generated: [date]

## Matched Components
| Design Element | Matched Component | File | Override Needed |
|---|---|---|---|
| Hero button | PrimaryButton | ui/components/Buttons.kt | text="Sign Up" |
| Section header | H2Text | ui/theme/Type.kt | — |

## New Components Needed
| Design Element | Reason | Suggested Implementation |
|---|---|---|
| Gradient card | No existing gradient container | Create GradientCard composable |
```

### 4. Generate implementation code
Draft the screen implementation using **only existing components** wherever possible:

- Import from the project's actual package paths
- Use the project's theme tokens (colors, typography, spacing) — never hardcode values
- Follow the project's existing code patterns (naming conventions, file organization)
- Add `// TODO: New component needed` comments where no existing match was found
- Include state handling and preview annotations where the framework supports it

Save the implementation to `DesignBridge/IMPLEMENTATION_[ScreenName].[ext]` (e.g., `.kt`, `.tsx`, `.dart`).

### 5. Generate a deviation report
Flag any inconsistencies between the reference design and the existing design system:

Create `DesignBridge/DEVIATION_REPORT.md`:

```markdown
# Design System Deviation Report
Generated: [date]

## Color Deviations
- Reference uses #FF6B35 for CTA button, but design system primary is #FF5722
  - **Recommendation:** Use existing primary or propose updating the token

## Typography Deviations
- Reference uses 18sp body text, design system body is 16sp
  - **Recommendation:** Use existing body style for consistency

## New Tokens Proposed
- accent_orange: #FF6B35 (used 3 times in reference)
```

### 6. Final output
All deliverables reside in the `DesignBridge/` directory:
```
DesignBridge/
├── COMPONENT_MAP.md
├── DEVIATION_REPORT.md
└── IMPLEMENTATION_[ScreenName].[ext]
```

## Important notes
- **Never create duplicate styles.** Always prefer reusing existing theme tokens and components over creating new ones.
- If the `component-catch` MCP server is available, use `extract_optimized_design_system` and `find_in_code` to speed up the design system cataloging step.
- If no reference image is provided, ask the user for one before proceeding.
- If the project has no established design system, note this and suggest creating foundational theme files first.
- The deviation report is critical for maintaining design consistency — always generate it.
- Skip any iOS-only files (`.swift`, `.storyboard`, `.xib`).
