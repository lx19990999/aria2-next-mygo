package main

import (
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
)

func (a *app) uiSettingsBody(c *ui.Context, t *ui.Theme) {
	L := a.L()

	ui.Column(c).Grow(1).Children(func() {
		ui.Row(c).Padding(10, 14).Gap(10).AlignItems(ui.Center).Background(t.Surface).Children(func() {
			ui.Text(c, L.UISettings).FontSize(a.fs(15)).FontWeight(600).Grow(1)
			if ui.PrimaryButton(c, L.Save).Padding(5, 12).Clicked() {
				a.commitScaleFromText()
				a.commitFontSizeFromText()
				a.commitAccentHexFromText()
				a.savePrefs()
			}
		})
		ui.Box(c).Height(1).Background(t.Border)

		ui.Scroll(c).Grow(1).Children(func() {
			ui.Column(c).Padding(12, 16).Gap(14).Children(func() {
				ui.Text(c, L.UIHint).FontSize(a.fs(12)).TextColor(t.TextMuted)

				a.uiPrefRow(c, t, L.Language, func() {
					prev := a.lang
					ui.Select(c, &a.lang, []string{"zh", "en"}).Width(180)
					ui.Textf(c, "%s / %s", L.LangZh, L.LangEn).FontSize(a.fs(11)).TextColor(t.TextMuted)
					if a.lang != prev {
						a.savePrefs()
					}
				})

				theme := "light"
				if a.themeDark {
					theme = "dark"
				}
				a.uiPrefRow(c, t, L.Theme, func() {
					ui.Select(c, &theme, []string{"light", "dark"}).Width(180)
					ui.Textf(c, "%s / %s", L.ThemeLight, L.ThemeDark).FontSize(a.fs(11)).TextColor(t.TextMuted)
				})
				if (theme == "dark") != a.themeDark {
					a.themeDark = theme == "dark"
					a.savePrefs()
				}

				a.uiPrefRow(c, t, L.UIScale, func() {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						if ui.Button(c, "−").Padding(4, 12).Clicked() {
							a.commitScaleFromText()
							a.bumpScale(-0.1)
						}
						in := ui.TextInput(c, &a.scaleText).Width(80)
						focused := in.Focused()
						if in.Submitted() || (a.scaleFocused && !focused) {
							if a.commitScaleFromText() {
								a.savePrefs()
							}
						}
						a.scaleFocused = focused
						if ui.Button(c, "+").Padding(4, 12).Clicked() {
							a.commitScaleFromText()
							a.bumpScale(0.1)
						}
						ui.Text(c, L.ScaleStepHint).FontSize(a.fs(11)).TextColor(t.TextMuted)
					})
				})

				a.uiPrefRow(c, t, L.FontFamily, func() {
					prev := a.fontFamily
					ui.Combobox(c, &a.fontFamily, systemFontNames()).Width(320)
					if a.fontFamily != prev && a.fontFamily != "" {
						a.savePrefs()
					}
				})

				a.uiPrefRow(c, t, L.FontSize, func() {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						if ui.Button(c, "−").Padding(4, 12).Clicked() {
							a.commitFontSizeFromText()
							a.bumpFontSize(-1)
						}
						in := ui.TextInput(c, &a.fontSizeText).Width(80)
						focused := in.Focused()
						if in.Submitted() || (a.fontSizeFocused && !focused) {
							if a.commitFontSizeFromText() {
								a.savePrefs()
							}
						}
						a.fontSizeFocused = focused
						if ui.Button(c, "+").Padding(4, 12).Clicked() {
							a.commitFontSizeFromText()
							a.bumpFontSize(1)
						}
						ui.Text(c, "pt").FontSize(a.fs(11)).TextColor(t.TextMuted)
					})
				})

				a.uiPrefRow(c, t, L.AccentColor, func() {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						for _, p := range accentPresets {
							col := ui.Hex(p.Hex)
							selected := strings.EqualFold(colorToHex(a.accentColor), p.Hex)
							b := ui.Button(c, " ").Width(28).Height(28).Padding(0).Radius(4).Background(col)
							if selected {
								b.Border(2, t.Text)
							} else {
								b.Border(1, t.Border)
							}
							if b.Clicked() {
								a.accentColor = col
								a.accentHex = p.Hex
								a.accentHexText = p.Hex
								a.savePrefs()
							}
						}
					})
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						for _, p := range accentPresets {
							ui.Text(c, p.Name).Width(28).FontSize(a.fs(10)).TextColor(t.TextMuted).TextAlign(ui.Center)
						}
					})
					ui.Row(c).Gap(10).AlignItems(ui.Center).Padding(4, 0).Children(func() {
						ui.Text(c, L.CustomColor).FontSize(a.fs(12)).TextColor(t.TextMuted)
						well := ui.ColorWell(c, &a.accentColor).Width(36).Height(28)
						if well.Changed() {
							a.accentHex = colorToHex(a.accentColor)
							a.accentHexText = a.accentHex
							a.savePrefs()
						}
						hex := ui.TextInput(c, &a.accentHexText).Width(100).Placeholder("#0d9488")
						hexFocused := hex.Focused()
						if hex.Submitted() || (a.accentHexFocused && !hexFocused) {
							if a.commitAccentHexFromText() {
								a.savePrefs()
							}
						}
						a.accentHexFocused = hexFocused
						if ui.Button(c, L.ApplyColor).Padding(4, 10).Clicked() {
							if a.commitAccentHexFromText() {
								a.savePrefs()
							} else {
								a.accentHex = colorToHex(a.accentColor)
								a.accentHexText = a.accentHex
								a.savePrefs()
							}
						}
					})
				})

				path, _ := prefsPath()
				ui.Textf(c, "%s  %s", L.PrefsPath, path).FontSize(a.fs(11)).TextColor(t.TextMuted)
			})
		})
	})
}

func (a *app) uiPrefRow(c *ui.Context, t *ui.Theme, label string, body func()) {
	ui.Column(c).Gap(6).Children(func() {
		ui.Text(c, label).FontSize(a.fs(12)).TextColor(t.TextMuted)
		body()
		ui.Box(c).Height(1).Background(t.Border)
	})
}

func (a *app) commitScaleFromText() bool {
	x, err := strconv.ParseFloat(strings.TrimSpace(a.scaleText), 32)
	if err != nil {
		a.scaleText = formatScale(a.uiScale)
		return false
	}
	if x < 0.5 {
		x = 0.5
	}
	if x > 3 {
		x = 3
	}
	a.uiScale = float32(int(x*10+0.5)) / 10
	a.scaleText = formatScale(a.uiScale)
	return true
}

func (a *app) commitFontSizeFromText() bool {
	x, err := strconv.ParseFloat(strings.TrimSpace(a.fontSizeText), 32)
	if err != nil {
		a.fontSizeText = formatFontSize(a.fontSize)
		return false
	}
	if x < 9 {
		x = 9
	}
	if x > 28 {
		x = 28
	}
	a.fontSize = float32(int(x + 0.5))
	a.fontSizeText = formatFontSize(a.fontSize)
	return true
}

func (a *app) commitAccentHexFromText() bool {
	c := parseAccent(a.accentHexText)
	a.accentColor = c
	a.accentHex = colorToHex(c)
	a.accentHexText = a.accentHex
	return true
}

func formatScale(v float32) string {
	return strings.TrimSpace(strings.Replace(strconv.FormatFloat(float64(v), 'f', 1, 32), ",", ".", 1))
}

func formatFontSize(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', 0, 32)
}
