// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gio.mleku.dev/font"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget"
)

// Theme holds the general theme of an app or window. Different top-level
// windows should have different instances of Theme (with different Shapers;
// see the godoc for [text.Shaper]), though their other fields can be equal.
type Theme struct {
	Shaper *text.Shaper
	*Colors
	TextSize unit.Sp
	Icon     struct {
		CheckBoxChecked   *widget.Icon
		CheckBoxUnchecked *widget.Icon
		RadioChecked      *widget.Icon
		RadioUnchecked    *widget.Icon
	}
	// Face selects the default typeface for text.
	Face font.Typeface

	// FingerSize is the minimum touch target size.
	FingerSize unit.Dp
}

// NewTheme constructs a theme (and underlying text shaper).
func NewTheme() *Theme {
	t := &Theme{Shaper: &text.Shaper{}}
	t.Colors = NewColors()
	t.TextSize = 16

	t.Icon.CheckBoxChecked = mustIcon(widget.NewIcon(icons.ToggleCheckBox))
	t.Icon.CheckBoxUnchecked = mustIcon(widget.NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.Icon.RadioChecked = mustIcon(widget.NewIcon(icons.ToggleRadioButtonChecked))
	t.Icon.RadioUnchecked = mustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked))

	// 38dp is on the lower end of possible finger size.
	t.FingerSize = 38

	return t
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
