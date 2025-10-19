# Material Design Color System Showcase

This example demonstrates the complete Material Design 3 color system implemented in Gio, showcasing all color roles and palette colors in both light and dark themes.

## Features

- **Complete Color System**: Displays all 30 Material Design color roles
- **Full Palette**: Shows all 66 palette colors (Primary, Secondary, Tertiary, Error, Neutral, Neutral Variant)
- **Theme Toggle**: Switch between light and dark themes to see how colors adapt
- **Color Information**: Each color shows its name and hex value
- **Responsive Layout**: Organized grid layout that adapts to window size

## Color Roles Displayed

### Primary Colors
- Primary, OnPrimary, PrimaryContainer, OnPrimaryContainer

### Secondary Colors  
- Secondary, OnSecondary, SecondaryContainer, OnSecondaryContainer

### Tertiary Colors
- Tertiary, OnTertiary, TertiaryContainer, OnTertiaryContainer

### Error Colors
- Error, OnError, ErrorContainer, OnErrorContainer

### Background Colors
- Background, OnBackground

### Surface Colors
- Surface, OnSurface, SurfaceVariant, OnSurfaceVariant

### Outline Colors
- Outline, OutlineVariant

### Utility Colors
- Shadow, Scrim, InverseSurface, InverseOnSurface, InversePrimary, SurfaceTint

## Palette Colors

Each color family (Primary, Secondary, Tertiary, Error, Neutral, Neutral Variant) displays all 11 shades from 50 (lightest) to 950 (darkest), plus the base 500 color.

## Running the Example

```bash
cd examples/material-colors
go run .
```

## Usage

- Click the "Switch to Dark/Light" button to toggle between themes
- Observe how color roles automatically adapt to maintain proper contrast
- Use this as a reference for implementing Material Design colors in your Gio applications

## Implementation Notes

This example uses the new Material Design 3 color system from `widget/material/colors.go`, demonstrating proper usage of:
- `ThemeMode` for light/dark theme switching
- Color role getters (`Primary()`, `OnSurface()`, etc.)
- Palette access for custom color variations
- Proper contrast ratios maintained across themes
