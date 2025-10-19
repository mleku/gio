# Colors Example

This example demonstrates the color scheme from fromage integrated into the gio widget package.

## Features Demonstrated

- **Material Design 3 Colors**: Complete color system with light and dark themes
- **Color Roles**: Semantic color roles for consistent theming
- **Theme Switching**: Support for light, dark, and auto theme modes
- **Custom Accents**: Ability to create custom color schemes
- **Surface Tinting**: Subtle surface tinting for enhanced visual hierarchy

## Color System

The color system includes:

### **Color Roles**
- **Primary**: Main brand color
- **Secondary**: Supporting accent color
- **Tertiary**: Additional accent color
- **Error**: Error states and warnings
- **Background**: Main background color
- **Surface**: Card and panel backgrounds
- **Outline**: Borders and dividers

### **Color Palettes**
- **Primary Palette**: 11 shades (50-950) of the primary color
- **Secondary Palette**: 11 shades (50-950) of the secondary color
- **Tertiary Palette**: 11 shades (50-950) of the tertiary color
- **Error Palette**: 11 shades (50-950) of the error color
- **Neutral Palette**: 11 shades (50-950) of neutral grays
- **Neutral Variant Palette**: 11 shades (50-950) of blue-grays

## Code Example

```go
// Create a color scheme
colors := widget.NewColorsWithMode(widget.ThemeModeLight)

// Use semantic color roles
paint.Fill(gtx.Ops, colors.Background())    // Background color
paint.Fill(gtx.Ops, colors.Primary())       // Primary color
paint.Fill(gtx.Ops, colors.Surface())       // Surface color
paint.Fill(gtx.Ops, colors.Outline())       // Outline color

// Switch themes
colors.SetThemeMode(widget.ThemeModeDark)
colors.ToggleTheme()

// Create custom color scheme
customColors := widget.NewColorsWithCustomAccents(0x00BCD4, 0x9C27B0) // Cyan primary, Purple secondary
```

## Key Features

- **Material Design 3**: Follows Google's Material Design 3 color system
- **Semantic Colors**: Use color roles instead of hardcoded colors
- **Theme Support**: Built-in light and dark theme support
- **Custom Accents**: Generate complete palettes from custom colors
- **Surface Tinting**: Subtle color tinting for enhanced hierarchy
- **HSV Support**: Advanced color manipulation with HSV color space
- **Accessibility**: Colors designed for proper contrast ratios

## Running the Example

```bash
go run ./examples/colors/main.go
```

## Visual Result

The example shows:
- Light theme background color filling the entire window
- Demonstrates the Material Design 3 color system
- Shows how to integrate colors with the widget system

## Use Cases

- **Consistent Theming**: Apply consistent colors across your application
- **Theme Switching**: Support light/dark mode switching
- **Brand Colors**: Use your brand colors with proper contrast
- **Accessibility**: Ensure proper color contrast for accessibility
- **Material Design**: Follow Material Design 3 guidelines
- **Custom Themes**: Create custom color schemes for your app
