package material

import (
	"fmt"
	"image/color"
	"math"
)

// ThemeMode represents the theme variant
type ThemeMode int

const (
	ThemeModeLight ThemeMode = iota
	ThemeModeDark
	ThemeModeAuto // Follows system preference
)

// Material Design color roles
type ColorRoles struct {
	// Primary colors
	Primary            color.NRGBA
	OnPrimary          color.NRGBA
	PrimaryContainer   color.NRGBA
	OnPrimaryContainer color.NRGBA

	// Secondary colors
	Secondary            color.NRGBA
	OnSecondary          color.NRGBA
	SecondaryContainer   color.NRGBA
	OnSecondaryContainer color.NRGBA

	// Tertiary colors
	Tertiary            color.NRGBA
	OnTertiary          color.NRGBA
	TertiaryContainer   color.NRGBA
	OnTertiaryContainer color.NRGBA

	// Error colors
	Error            color.NRGBA
	OnError          color.NRGBA
	ErrorContainer   color.NRGBA
	OnErrorContainer color.NRGBA

	// Background colors
	Background   color.NRGBA
	OnBackground color.NRGBA

	// Surface colors
	Surface          color.NRGBA
	OnSurface        color.NRGBA
	SurfaceVariant   color.NRGBA
	OnSurfaceVariant color.NRGBA

	// Outline colors
	Outline        color.NRGBA
	OutlineVariant color.NRGBA

	// Shadow and scrim
	Shadow color.NRGBA
	Scrim  color.NRGBA

	// Inverse colors
	InverseSurface   color.NRGBA
	InverseOnSurface color.NRGBA
	InversePrimary   color.NRGBA

	// Surface tint
	SurfaceTint color.NRGBA
}

// Material Design color palette
type ColorPalette struct {
	// Primary palette
	Primary50  color.NRGBA
	Primary100 color.NRGBA
	Primary200 color.NRGBA
	Primary300 color.NRGBA
	Primary400 color.NRGBA
	Primary500 color.NRGBA
	Primary600 color.NRGBA
	Primary700 color.NRGBA
	Primary800 color.NRGBA
	Primary900 color.NRGBA
	Primary950 color.NRGBA

	// Secondary palette
	Secondary50  color.NRGBA
	Secondary100 color.NRGBA
	Secondary200 color.NRGBA
	Secondary300 color.NRGBA
	Secondary400 color.NRGBA
	Secondary500 color.NRGBA
	Secondary600 color.NRGBA
	Secondary700 color.NRGBA
	Secondary800 color.NRGBA
	Secondary900 color.NRGBA
	Secondary950 color.NRGBA

	// Tertiary palette
	Tertiary50  color.NRGBA
	Tertiary100 color.NRGBA
	Tertiary200 color.NRGBA
	Tertiary300 color.NRGBA
	Tertiary400 color.NRGBA
	Tertiary500 color.NRGBA
	Tertiary600 color.NRGBA
	Tertiary700 color.NRGBA
	Tertiary800 color.NRGBA
	Tertiary900 color.NRGBA
	Tertiary950 color.NRGBA

	// Error palette
	Error50  color.NRGBA
	Error100 color.NRGBA
	Error200 color.NRGBA
	Error300 color.NRGBA
	Error400 color.NRGBA
	Error500 color.NRGBA
	Error600 color.NRGBA
	Error700 color.NRGBA
	Error800 color.NRGBA
	Error900 color.NRGBA
	Error950 color.NRGBA

	// Neutral palette
	Neutral50  color.NRGBA
	Neutral100 color.NRGBA
	Neutral200 color.NRGBA
	Neutral300 color.NRGBA
	Neutral400 color.NRGBA
	Neutral500 color.NRGBA
	Neutral600 color.NRGBA
	Neutral700 color.NRGBA
	Neutral800 color.NRGBA
	Neutral900 color.NRGBA
	Neutral950 color.NRGBA

	// Neutral variant palette
	NeutralVariant50  color.NRGBA
	NeutralVariant100 color.NRGBA
	NeutralVariant200 color.NRGBA
	NeutralVariant300 color.NRGBA
	NeutralVariant400 color.NRGBA
	NeutralVariant500 color.NRGBA
	NeutralVariant600 color.NRGBA
	NeutralVariant700 color.NRGBA
	NeutralVariant800 color.NRGBA
	NeutralVariant900 color.NRGBA
	NeutralVariant950 color.NRGBA
}

type Colors struct {
	roles       ColorRoles
	palette     ColorPalette
	themeMode   ThemeMode
	surfaceTint color.NRGBA // Custom surface tint color
}

func NewColors() *Colors {
	return NewColorsWithMode(ThemeModeLight)
}

func NewColorsWithMode(mode ThemeMode) *Colors {
	theme := &Colors{
		themeMode:   mode,
		surfaceTint: hex("#263238"), // Default blue-gray
	}

	// Initialize with Material Design 3 default colors
	theme.initDefaultPalette()
	theme.initRoles()

	return theme
}

func NewColorsWithCustomAccents(primary, secondary uint32) *Colors {
	return NewColorsWithCustomAccentsAndMode(primary, secondary, ThemeModeLight)
}

func NewColorsWithCustomAccentsAndMode(primary, secondary uint32, mode ThemeMode) *Colors {
	theme := &Colors{
		themeMode:   mode,
		surfaceTint: hex("#263238"), // Default blue-gray
	}

	// Initialize with custom colors
	theme.initCustomPalette(primary, secondary)
	theme.initRoles()

	return theme
}

// Initialize default Material Design 3 palette
func (t *Colors) initDefaultPalette() {
	// Primary - Blue
	t.palette.Primary50 = hex("#E0F7FA")
	t.palette.Primary100 = hex("#B2EBF2")
	t.palette.Primary200 = hex("#80DEEA")
	t.palette.Primary300 = hex("#4DD0E1")
	t.palette.Primary400 = hex("#26C6DA")
	t.palette.Primary500 = hex("#00BCD4")
	t.palette.Primary600 = hex("#00ACC1")
	t.palette.Primary700 = hex("#0097A7")
	t.palette.Primary800 = hex("#00838F")
	t.palette.Primary900 = hex("#006064")
	t.palette.Primary950 = hex("#004D40")

	// Secondary - Purple
	t.palette.Secondary50 = hex("#F3E5F5")
	t.palette.Secondary100 = hex("#E1BEE7")
	t.palette.Secondary200 = hex("#CE93D8")
	t.palette.Secondary300 = hex("#BA68C8")
	t.palette.Secondary400 = hex("#AB47BC")
	t.palette.Secondary500 = hex("#9C27B0")
	t.palette.Secondary600 = hex("#8E24AA")
	t.palette.Secondary700 = hex("#7B1FA2")
	t.palette.Secondary800 = hex("#6A1B9A")
	t.palette.Secondary900 = hex("#4A148C")
	t.palette.Secondary950 = hex("#3F0F7A")

	// Tertiary - Teal
	t.palette.Tertiary50 = hex("#E0F2F1")
	t.palette.Tertiary100 = hex("#B2DFDB")
	t.palette.Tertiary200 = hex("#80CBC4")
	t.palette.Tertiary300 = hex("#4DB6AC")
	t.palette.Tertiary400 = hex("#26A69A")
	t.palette.Tertiary500 = hex("#009688")
	t.palette.Tertiary600 = hex("#00897B")
	t.palette.Tertiary700 = hex("#00796B")
	t.palette.Tertiary800 = hex("#00695C")
	t.palette.Tertiary900 = hex("#004D40")
	t.palette.Tertiary950 = hex("#003D33")

	// Error - Red
	t.palette.Error50 = hex("#FFEBEE")
	t.palette.Error100 = hex("#FFCDD2")
	t.palette.Error200 = hex("#EF9A9A")
	t.palette.Error300 = hex("#E57373")
	t.palette.Error400 = hex("#EF5350")
	t.palette.Error500 = hex("#F44336")
	t.palette.Error600 = hex("#E53935")
	t.palette.Error700 = hex("#D32F2F")
	t.palette.Error800 = hex("#C62828")
	t.palette.Error900 = hex("#B71C1C")
	t.palette.Error950 = hex("#A01515")

	// Neutral - Gray
	t.palette.Neutral50 = hex("#FAFAFA")
	t.palette.Neutral100 = hex("#F5F5F5")
	t.palette.Neutral200 = hex("#EEEEEE")
	t.palette.Neutral300 = hex("#E0E0E0")
	t.palette.Neutral400 = hex("#BDBDBD")
	t.palette.Neutral500 = hex("#9E9E9E")
	t.palette.Neutral600 = hex("#757575")
	t.palette.Neutral700 = hex("#616161")
	t.palette.Neutral800 = hex("#424242")
	t.palette.Neutral900 = hex("#212121")
	t.palette.Neutral950 = hex("#0F0F0F")

	// Neutral Variant - Blue-gray
	t.palette.NeutralVariant50 = hex("#F8FAFC")
	t.palette.NeutralVariant100 = hex("#F1F5F9")
	t.palette.NeutralVariant200 = hex("#E2E8F0")
	t.palette.NeutralVariant300 = hex("#CBD5E1")
	t.palette.NeutralVariant400 = hex("#94A3B8")
	t.palette.NeutralVariant500 = hex("#64748B")
	t.palette.NeutralVariant600 = hex("#475569")
	t.palette.NeutralVariant700 = hex("#334155")
	t.palette.NeutralVariant800 = hex("#1E293B")
	t.palette.NeutralVariant900 = hex("#0F172A")
	t.palette.NeutralVariant950 = hex("#020617")
}

// Initialize custom palette with provided colors
func (t *Colors) initCustomPalette(primary, secondary uint32) {
	// Use provided colors as base and generate palette
	t.palette.Primary500 = rgb(primary)
	t.palette.Secondary500 = rgb(secondary)

	// Generate variations (simplified - in practice you'd use proper color theory)
	t.palette.Primary50 = lighten(t.palette.Primary500, 0.9)
	t.palette.Primary100 = lighten(t.palette.Primary500, 0.8)
	t.palette.Primary200 = lighten(t.palette.Primary500, 0.6)
	t.palette.Primary300 = lighten(t.palette.Primary500, 0.4)
	t.palette.Primary400 = lighten(t.palette.Primary500, 0.2)
	t.palette.Primary600 = darken(t.palette.Primary500, 0.1)
	t.palette.Primary700 = darken(t.palette.Primary500, 0.2)
	t.palette.Primary800 = darken(t.palette.Primary500, 0.3)
	t.palette.Primary900 = darken(t.palette.Primary500, 0.4)
	t.palette.Primary950 = darken(t.palette.Primary500, 0.5)

	// Similar for secondary...
	t.palette.Secondary50 = lighten(t.palette.Secondary500, 0.9)
	t.palette.Secondary100 = lighten(t.palette.Secondary500, 0.8)
	t.palette.Secondary200 = lighten(t.palette.Secondary500, 0.6)
	t.palette.Secondary300 = lighten(t.palette.Secondary500, 0.4)
	t.palette.Secondary400 = lighten(t.palette.Secondary500, 0.2)
	t.palette.Secondary600 = darken(t.palette.Secondary500, 0.1)
	t.palette.Secondary700 = darken(t.palette.Secondary500, 0.2)
	t.palette.Secondary800 = darken(t.palette.Secondary500, 0.3)
	t.palette.Secondary900 = darken(t.palette.Secondary500, 0.4)
	t.palette.Secondary950 = darken(t.palette.Secondary500, 0.5)

	// Initialize other palettes with defaults
	t.initDefaultPalette()
}

// Initialize Material Design 3 color roles based on theme mode
func (t *Colors) initRoles() {
	switch t.themeMode {
	case ThemeModeDark:
		t.initDarkRoles()
	case ThemeModeLight:
		fallthrough
	default:
		t.initLightRoles()
	}
}

// Initialize light theme color roles
func (t *Colors) initLightRoles() {
	// Primary roles
	t.roles.Primary = t.palette.Primary500
	t.roles.OnPrimary = t.palette.Neutral50
	t.roles.PrimaryContainer = t.palette.Primary100
	t.roles.OnPrimaryContainer = t.palette.Primary900

	// Secondary roles
	t.roles.Secondary = t.palette.Secondary500
	t.roles.OnSecondary = t.palette.Neutral50
	t.roles.SecondaryContainer = t.palette.Secondary100
	t.roles.OnSecondaryContainer = t.palette.Secondary900

	// Tertiary roles
	t.roles.Tertiary = t.palette.Tertiary500
	t.roles.OnTertiary = t.palette.Neutral50
	t.roles.TertiaryContainer = t.palette.Tertiary100
	t.roles.OnTertiaryContainer = t.palette.Tertiary900

	// Error roles
	t.roles.Error = t.palette.Error500
	t.roles.OnError = t.palette.Neutral50
	t.roles.ErrorContainer = t.palette.Error100
	t.roles.OnErrorContainer = t.palette.Error900

	// Background roles
	t.roles.Background = t.palette.Neutral50
	t.roles.OnBackground = t.palette.Neutral900

	// Surface roles with subtle tint for contrast
	t.roles.Surface = t.applySurfaceTint(t.palette.Neutral50)
	t.roles.OnSurface = t.palette.Neutral900
	t.roles.SurfaceVariant = t.palette.NeutralVariant100
	t.roles.OnSurfaceVariant = t.palette.NeutralVariant800

	// Outline roles
	t.roles.Outline = t.palette.NeutralVariant500
	t.roles.OutlineVariant = t.palette.NeutralVariant300

	// Shadow and scrim
	t.roles.Shadow = hex("#000000")
	t.roles.Scrim = hex("#000000")

	// Inverse roles
	t.roles.InverseSurface = t.palette.Neutral800
	t.roles.InverseOnSurface = t.palette.Neutral100
	t.roles.InversePrimary = t.palette.Primary200

	// Surface tint
	t.roles.SurfaceTint = t.palette.Primary500
}

// Initialize dark theme color roles
func (t *Colors) initDarkRoles() {
	// Primary roles
	t.roles.Primary = t.palette.Primary200
	t.roles.OnPrimary = t.palette.Primary900
	t.roles.PrimaryContainer = t.palette.Primary800
	t.roles.OnPrimaryContainer = t.palette.Primary100

	// Secondary roles
	t.roles.Secondary = t.palette.Secondary200
	t.roles.OnSecondary = t.palette.Secondary900
	t.roles.SecondaryContainer = t.palette.Secondary800
	t.roles.OnSecondaryContainer = t.palette.Secondary100

	// Tertiary roles
	t.roles.Tertiary = t.palette.Tertiary200
	t.roles.OnTertiary = t.palette.Tertiary900
	t.roles.TertiaryContainer = t.palette.Tertiary800
	t.roles.OnTertiaryContainer = t.palette.Tertiary100

	// Error roles
	t.roles.Error = t.palette.Error200
	t.roles.OnError = t.palette.Error900
	t.roles.ErrorContainer = t.palette.Error800
	t.roles.OnErrorContainer = t.palette.Error100

	// Background roles
	t.roles.Background = t.palette.Neutral900
	t.roles.OnBackground = t.palette.Neutral50

	// Surface roles with subtle tint for contrast
	t.roles.Surface = t.applySurfaceTint(t.palette.Neutral900)
	t.roles.OnSurface = t.palette.Neutral50
	t.roles.SurfaceVariant = t.palette.NeutralVariant800
	t.roles.OnSurfaceVariant = t.palette.NeutralVariant200

	// Outline roles
	t.roles.Outline = t.palette.NeutralVariant400
	t.roles.OutlineVariant = t.palette.NeutralVariant700

	// Shadow and scrim
	t.roles.Shadow = hex("#000000")
	t.roles.Scrim = hex("#000000")

	// Inverse roles
	t.roles.InverseSurface = t.palette.Neutral100
	t.roles.InverseOnSurface = t.palette.Neutral800
	t.roles.InversePrimary = t.palette.Primary800

	// Surface tint
	t.roles.SurfaceTint = t.palette.Primary200
}

// Getters for color roles
func (t *Colors) Primary() color.NRGBA            { return t.roles.Primary }
func (t *Colors) OnPrimary() color.NRGBA          { return t.roles.OnPrimary }
func (t *Colors) PrimaryContainer() color.NRGBA   { return t.roles.PrimaryContainer }
func (t *Colors) OnPrimaryContainer() color.NRGBA { return t.roles.OnPrimaryContainer }

func (t *Colors) Secondary() color.NRGBA            { return t.roles.Secondary }
func (t *Colors) OnSecondary() color.NRGBA          { return t.roles.OnSecondary }
func (t *Colors) SecondaryContainer() color.NRGBA   { return t.roles.SecondaryContainer }
func (t *Colors) OnSecondaryContainer() color.NRGBA { return t.roles.OnSecondaryContainer }

func (t *Colors) Tertiary() color.NRGBA            { return t.roles.Tertiary }
func (t *Colors) OnTertiary() color.NRGBA          { return t.roles.OnTertiary }
func (t *Colors) TertiaryContainer() color.NRGBA   { return t.roles.TertiaryContainer }
func (t *Colors) OnTertiaryContainer() color.NRGBA { return t.roles.OnTertiaryContainer }

func (t *Colors) Error() color.NRGBA            { return t.roles.Error }
func (t *Colors) OnError() color.NRGBA          { return t.roles.OnError }
func (t *Colors) ErrorContainer() color.NRGBA   { return t.roles.ErrorContainer }
func (t *Colors) OnErrorContainer() color.NRGBA { return t.roles.OnErrorContainer }

func (t *Colors) Background() color.NRGBA   { return t.roles.Background }
func (t *Colors) OnBackground() color.NRGBA { return t.roles.OnBackground }

func (t *Colors) Surface() color.NRGBA          { return t.roles.Surface }
func (t *Colors) OnSurface() color.NRGBA        { return t.roles.OnSurface }
func (t *Colors) SurfaceVariant() color.NRGBA   { return t.roles.SurfaceVariant }
func (t *Colors) OnSurfaceVariant() color.NRGBA { return t.roles.OnSurfaceVariant }

func (t *Colors) Outline() color.NRGBA        { return t.roles.Outline }
func (t *Colors) OutlineVariant() color.NRGBA { return t.roles.OutlineVariant }

func (t *Colors) Shadow() color.NRGBA { return t.roles.Shadow }
func (t *Colors) Scrim() color.NRGBA  { return t.roles.Scrim }

func (t *Colors) InverseSurface() color.NRGBA   { return t.roles.InverseSurface }
func (t *Colors) InverseOnSurface() color.NRGBA { return t.roles.InverseOnSurface }
func (t *Colors) InversePrimary() color.NRGBA   { return t.roles.InversePrimary }

func (t *Colors) SurfaceTint() color.NRGBA { return t.roles.SurfaceTint }

// Surface tint control methods
func (t *Colors) SetSurfaceTint(tint color.NRGBA) {
	t.surfaceTint = tint
	t.initRoles() // Reinitialize roles to apply new tint
}

func (t *Colors) GetSurfaceTint() color.NRGBA { return t.surfaceTint }

// SetSurfaceTintFromHSV sets the surface tint from HSV values
func (t *Colors) SetSurfaceTintFromHSV(hue, saturation, value float32) *Colors {
	// Convert HSV to RGB
	hueDegrees := hue * 360
	tint := t.hsvToRgb(hueDegrees, saturation, value)
	t.SetSurfaceTint(tint)
	return t
}

// SetSurfaceTintTone sets the tone (brightness/value) of the surface tint
func (t *Colors) SetSurfaceTintTone(tone float32) *Colors {
	h, s, _ := t.rgbToHsv(t.surfaceTint)
	t.SetSurfaceTintFromHSV(h/360, s, tone)
	return t
}

// SetSurfaceTintHue sets the hue of the surface tint
func (t *Colors) SetSurfaceTintHue(hue float32) *Colors {
	_, s, v := t.rgbToHsv(t.surfaceTint)
	t.SetSurfaceTintFromHSV(hue, s, v)
	return t
}

// SetSurfaceTintSaturation sets the saturation of the surface tint
func (t *Colors) SetSurfaceTintSaturation(saturation float32) *Colors {
	h, _, v := t.rgbToHsv(t.surfaceTint)
	t.SetSurfaceTintFromHSV(h/360, saturation, v)
	return t
}

// GetSurfaceTintHSV returns the HSV values of the current surface tint
func (t *Colors) GetSurfaceTintHSV() (hue, saturation, value float32) {
	h, s, v := t.rgbToHsv(t.surfaceTint)
	return h / 360, s, v // Convert hue to 0-1 range
}

// Getters for palette colors
func (t *Colors) Palette() ColorPalette { return t.palette }

// Theme mode methods
func (t *Colors) ThemeMode() ThemeMode { return t.themeMode }

func (t *Colors) SetThemeMode(mode ThemeMode) {
	if t.themeMode != mode {
		t.themeMode = mode
		t.initRoles()
	}
}

func (t *Colors) ToggleTheme() {
	if t.themeMode == ThemeModeLight {
		t.SetThemeMode(ThemeModeDark)
	} else {
		t.SetThemeMode(ThemeModeLight)
	}
}

// Utility functions
func hex(s string) color.NRGBA {
	if len(s) != 7 || s[0] != '#' {
		return color.NRGBA{}
	}

	var r, g, b uint8
	if _, err := fmt.Sscanf(s[1:], "%02x%02x%02x", &r, &g, &b); err != nil {
		return color.NRGBA{}
	}

	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

// ColorToHex converts a color.NRGBA to hex string format
func ColorToHex(c color.NRGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

// Simple color manipulation functions
func lighten(c color.NRGBA, factor float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(float64(c.R) + (255-float64(c.R))*factor),
		G: uint8(float64(c.G) + (255-float64(c.G))*factor),
		B: uint8(float64(c.B) + (255-float64(c.B))*factor),
		A: c.A,
	}
}

func darken(c color.NRGBA, factor float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(float64(c.R) * (1 - factor)),
		G: uint8(float64(c.G) * (1 - factor)),
		B: uint8(float64(c.B) * (1 - factor)),
		A: c.A,
	}
}

// blendColors blends two colors with a given ratio
func (t *Colors) blendColors(base, tint color.NRGBA, ratio float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(float64(base.R)*(1-ratio) + float64(tint.R)*ratio),
		G: uint8(float64(base.G)*(1-ratio) + float64(tint.G)*ratio),
		B: uint8(float64(base.B)*(1-ratio) + float64(tint.B)*ratio),
		A: base.A,
	}
}

// applySurfaceTint applies a custom tint to the surface color for contrast
func (t *Colors) applySurfaceTint(baseColor color.NRGBA) color.NRGBA {
	if t.ThemeMode() == ThemeModeLight {
		// Light mode: use the brightness complement to work with bright white background
		return t.invertValue(t.surfaceTint)
	} else {
		// Dark mode: use the surface tint color as-is
		return t.surfaceTint
	}
}

// invertValue inverts the brightness (value) of a color while preserving hue and saturation
func (t *Colors) invertValue(c color.NRGBA) color.NRGBA {
	// Convert RGB to HSV
	h, s, v := t.rgbToHsv(c)

	// Invert the value (brightness)
	v = 1.0 - v

	// Convert back to RGB
	return t.hsvToRgb(h, s, v)
}

// hsvToRgb converts HSV to RGB
func (t *Colors) hsvToRgb(h, s, v float32) color.NRGBA {
	h = h / 360.0 // Normalize to [0, 1]
	c := v * s
	x := c * (1 - float32(math.Abs(float64(math.Mod(float64(h*6), 2))-1)))
	m := v - c

	var r, g, b float32
	if h < 1.0/6.0 {
		r, g, b = c, x, 0
	} else if h < 2.0/6.0 {
		r, g, b = x, c, 0
	} else if h < 3.0/6.0 {
		r, g, b = 0, c, x
	} else if h < 4.0/6.0 {
		r, g, b = 0, x, c
	} else if h < 5.0/6.0 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}

	return color.NRGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}

// rgbToHsv converts RGB to HSV
func (t *Colors) rgbToHsv(c color.NRGBA) (float32, float32, float32) {
	r := float32(c.R) / 255.0
	g := float32(c.G) / 255.0
	b := float32(c.B) / 255.0

	max := float32(math.Max(float64(r), math.Max(float64(g), float64(b))))
	min := float32(math.Min(float64(r), math.Min(float64(g), float64(b))))
	delta := max - min

	var h float32
	if delta == 0 {
		h = 0
	} else if max == r {
		h = float32(60 * math.Mod(float64((g-b)/delta), 6))
	} else if max == g {
		h = 60 * ((b-r)/delta + 2)
	} else {
		h = 60 * ((r-g)/delta + 4)
	}

	if h < 0 {
		h += 360
	}

	s := float32(0)
	if max != 0 {
		s = delta / max
	}

	v := max

	return h, s, v
}
