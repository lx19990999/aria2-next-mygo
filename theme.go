package main

import "github.com/egoist/mygo/ui"

func applyTheme(c *ui.Context, a *app) *ui.Theme {
	a.defaultPrefs()
	t := ui.LightTheme()
	if a.themeDark {
		t = ui.DarkTheme()
		t.Background = ui.Hex("#18181b")
		t.Surface = ui.Hex("#27272a")
		t.SurfaceHover = ui.Hex("#3f3f46")
		t.Border = ui.Hex("#3f3f46")
		t.Text = ui.Hex("#fafafa")
		t.TextMuted = ui.Hex("#a1a1aa")
		t.Success = ui.Hex("#4ade80")
		t.Danger = ui.Hex("#f87171")
	} else {
		t.Background = ui.Hex("#f4f4f5")
		t.Surface = ui.Hex("#ffffff")
		t.SurfaceHover = ui.Hex("#f4f4f5")
		t.Border = ui.Hex("#e4e4e7")
		t.Text = ui.Hex("#18181b")
		t.TextMuted = ui.Hex("#71717a")
		t.Success = ui.Hex("#16a34a")
		t.Danger = ui.Hex("#dc2626")
	}

	accent := a.accentColor
	if accent.A == 0 {
		accent = parseAccent(a.accentHex)
	}
	t.Accent = accent
	t.AccentHover = mixToward(accent, ui.Hex("#000000"), 0.12)
	t.AccentPressed = mixToward(accent, ui.Hex("#000000"), 0.22)
	if a.themeDark {
		t.AccentText = ui.Hex("#042f2e")
		if luminance(accent) < 0.55 {
			t.AccentText = ui.Hex("#ffffff")
		}
	} else {
		t.AccentText = ui.Hex("#ffffff")
		if luminance(accent) > 0.72 {
			t.AccentText = ui.Hex("#18181b")
		}
	}

	t.Font = a.fontFamily
	t.FontSize = a.fontSize
	if t.FontSize <= 0 {
		t.FontSize = 13
	}
	spacing := float32(4) * a.uiScale
	if spacing < 3 {
		spacing = 3
	}
	if spacing > 8 {
		spacing = 8
	}
	t.Spacing = spacing
	t.Radius = 3
	c.SetTheme(t)
	return t
}

func mixToward(c, toward ui.Color, t float32) ui.Color {
	return c.Mix(toward, t)
}

func luminance(c ui.Color) float32 {
	return (0.2126*float32(c.R) + 0.7152*float32(c.G) + 0.0722*float32(c.B)) / 255
}

func (a *app) chromeBG() ui.Color {
	if a.themeDark {
		return ui.Hex("#1f1f23")
	}
	return ui.Hex("#fafafa")
}

func (a *app) selectedBG() ui.Color {
	acc := a.accentColor
	if a.themeDark {
		return acc.Alpha(0.28)
	}
	return acc.Alpha(0.10)
}

func (a *app) selectedMark() ui.Color {
	return a.accentColor
}
