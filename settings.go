package main

import (
	"fmt"
	"strings"

	"aria2-next-mygo/options"

	"github.com/egoist/mygo/ui"
)

const (
	anchorConn = "conn"
)

func (a *app) ensureOptMaps() {
	if a.optVals == nil {
		a.optVals = map[string]string{}
	}
	if a.optOrig == nil {
		a.optOrig = map[string]string{}
	}
}

func (a *app) loadGlobalOptions() {
	a.ensureOptMaps()
	a.optsLoading = true
	a.flash("正在加载全局选项…")
	go func() {
		opts, err := a.client().GetGlobalOption()
		a.apply(func() {
			a.optsLoading = false
			a.ensureOptMaps()
			if err != nil {
				a.optsLoaded = false
				a.flash("加载选项失败：" + err.Error())
				return
			}
			next := map[string]string{}
			for _, sec := range options.Catalog {
				for _, item := range sec.Items {
					if v, ok := opts[item.Key]; ok {
						next[item.Key] = v
					} else if item.Default != "" {
						next[item.Key] = item.Default
					} else {
						next[item.Key] = ""
					}
				}
			}
			for k, v := range opts {
				if _, ok := next[k]; !ok {
					next[k] = v
				}
			}
			a.optVals = next
			a.optOrig = cloneMap(next)
			a.optsLoaded = true
			a.flash(fmt.Sprintf("已加载 %d 项全局选项", len(opts)))
		})
	}()
}

func (a *app) dirtyOptions() map[string]string {
	a.ensureOptMaps()
	out := map[string]string{}
	for k, v := range a.optVals {
		if a.optOrig[k] != v {
			out[k] = v
		}
	}
	return out
}

func (a *app) saveGlobalOptions() {
	dirty := a.dirtyOptions()
	if len(dirty) == 0 {
		a.flash("没有需要保存的修改")
		return
	}
	a.flash(fmt.Sprintf("正在保存 %d 项…", len(dirty)))
	go func() {
		err := a.client().ChangeGlobalOption(dirty)
		a.apply(func() {
			if err != nil {
				a.flash("保存失败：" + err.Error())
				return
			}
			for k, v := range dirty {
				a.optOrig[k] = v
			}
			a.flash(fmt.Sprintf("已保存 %d 项", len(dirty)))
		})
		a.refresh()
	}()
}

func (a *app) settingsBody(c *ui.Context, t *ui.Theme) {
	L := a.L()
	a.ensureOptMaps()
	if a.section == "settings" && !a.optsLoaded && !a.optsLoading && a.connected {
		a.loadGlobalOptions()
	}
	if a.settingsAnchor == "" {
		a.settingsAnchor = anchorConn
	}
	dirtyN := len(a.dirtyOptions())
	filter := strings.ToLower(strings.TrimSpace(a.optFilter))

	ui.Column(c).Grow(1).Children(func() {
		ui.Row(c).Padding(10, 14).Gap(10).AlignItems(ui.Center).Background(t.Surface).Children(func() {
			ui.Text(c, L.EngineSettings).FontSize(a.fs(15)).FontWeight(600)
			ui.TextInput(c, &a.optFilter).Placeholder(L.FilterParams).Grow(1)
			if ui.Button(c, L.Reload).Padding(5, 10).Clicked() {
				a.optsLoaded = false
				a.loadGlobalOptions()
			}
			saveLabel := L.Save
			if dirtyN > 0 {
				saveLabel = fmt.Sprintf(L.SaveN, dirtyN)
			}
			if ui.PrimaryButton(c, saveLabel).Padding(5, 12).Clicked() {
				a.saveGlobalOptions()
			}
		})
		ui.Box(c).Height(1).Background(t.Border)

		ui.Row(c).Grow(1).AlignItems(ui.Stretch).Children(func() {
			ui.Column(c).Width(140).Shrink(0).Background(a.chromeBG()).Padding(8, 6).Gap(2).Children(func() {
				a.settingsRailItem(c, t, anchorConn, L.Connection)
				ui.Box(c).Height(8)
				for _, sec := range options.Catalog {
					a.settingsRailItem(c, t, sec.ID, secTitle(sec, a.lang))
				}
			})
			ui.Box(c).Width(1).Background(t.Border)

			ui.Scroll(c).Grow(1).Children(func() {
				ui.Column(c).Padding(12, 16).Gap(18).Children(func() {
					a.connectionSettings(c, t)

					shownAny := false
					for _, sec := range options.Catalog {
						n := a.optionSection(c, t, sec, filter)
						if n > 0 {
							shownAny = true
						}
					}
					if filter != "" && !shownAny {
						ui.Text(c, L.NoMatch).FontSize(a.fs(12)).TextColor(t.TextMuted)
					}
				})
			})
		})
	})
}

func secTitle(sec options.Section, lang string) string {
	if lang == "en" {
		return sec.ID
	}
	return sec.Title
}

func (a *app) settingsRailItem(c *ui.Context, t *ui.Theme, id, label string) {
	selected := a.settingsAnchor == id
	row := ui.Row(c).FillWidth().Gap(0).AlignItems(ui.Stretch).Radius(3)
	if selected {
		row.Background(a.selectedBG())
	}
	row.Children(func() {
		mark := t.Background.Alpha(0)
		if selected {
			mark = a.selectedMark()
		}
		ui.Box(c).Width(3).Background(mark)
		color := t.Text
		weight := 500
		if selected {
			color = t.Accent
			weight = 600
		}
		ui.Text(c, label).FontSize(a.fs(12)).FontWeight(weight).TextColor(color).
			Padding(7, 10).SingleLine().Grow(1)
	})
	if row.Clicked() {
		a.settingsAnchor = id
		a.pendingScroll = id
	}
}

func (a *app) sectionHeader(c *ui.Context, t *ui.Theme, id, title, hint string) {
	L := a.L()
	h := ui.Column(c.Key("sec-"+id)).Gap(4).Padding(4, 0, 8, 0)
	h.Children(func() {
		ui.Text(c, title).FontSize(a.fs(14)).FontWeight(600)
		if hint != "" {
			ui.Text(c, hint).FontSize(a.fs(11)).TextColor(t.TextMuted)
		}
		_ = L
	})
	if a.pendingScroll == id {
		h.ScrollIntoView(ui.ScrollIntoViewOptions{Block: ui.ScrollStart})
		a.pendingScroll = ""
	}
}

func (a *app) connectionSettings(c *ui.Context, t *ui.Theme) {
	L := a.L()
	a.sectionHeader(c, t, anchorConn, L.RPCTitle, L.RPCHint)

	ui.Column(c).Gap(8).Children(func() {
		a.denseField(c, t, L.Address, func() {
			ui.TextInput(c, &a.rpcURL).FillWidth()
		})
		a.denseField(c, t, L.Secret, func() {
			ui.TextInput(c, &a.rpcSecret).FillWidth()
		})
		ui.Row(c).Gap(10).AlignItems(ui.Center).Padding(4, 0).Children(func() {
			if ui.PrimaryButton(c, L.Reconnect).Padding(6, 14).Clicked() {
				a.rpcPaused = false
				a.savePrefs()
				a.flash("正在连接…")
				a.optsLoaded = false
				go a.refresh()
				a.loadGlobalOptions()
			}
			if ui.Button(c, L.Disconnect).Disabled(a.rpcPaused).Padding(6, 14).Clicked() {
				a.rpcPaused = true
				a.connected = false
				a.connErr = ""
				a.optsLoaded = false
				a.optVals = nil
				a.optOrig = nil
				a.updateWindowTitle()
				a.flash(L.Disconnected)
			}
			if a.connected {
				ui.Text(c, L.Connected).FontSize(a.fs(12)).TextColor(t.Success)
				ui.Textf(c, "%s %s", a.version.Product, a.version.Version).FontSize(a.fs(12)).TextColor(t.TextMuted)
			} else if a.rpcPaused {
				ui.Text(c, L.Disconnected).FontSize(a.fs(12)).TextColor(t.TextMuted)
			} else {
				ui.Text(c, L.NotConnected).FontSize(a.fs(12)).TextColor(t.Danger)
			}
		})
		if a.connErr != "" {
			ui.Text(c, a.connErr).FontSize(a.fs(12)).TextColor(t.Danger)
		}
		ui.Text(c, L.CatalogHint).FontSize(a.fs(11)).TextColor(t.TextMuted)
	})
}

func (a *app) denseField(c *ui.Context, t *ui.Theme, label string, control func()) {
	ui.Row(c).Gap(12).AlignItems(ui.Center).Padding(4, 0).Children(func() {
		ui.Text(c, label).Width(56).FontSize(a.fs(12)).TextColor(t.TextMuted)
		ui.Column(c).Grow(1).Children(control)
	})
}

func (a *app) optionSection(c *ui.Context, t *ui.Theme, sec options.Section, filter string) int {
	L := a.L()
	type row struct {
		item options.Option
		name string
		desc string
	}
	rows := make([]row, 0, len(sec.Items))
	for _, item := range sec.Items {
		name, desc := options.Display(item.Key, item.Desc, a.lang)
		if filter != "" {
			blob := strings.ToLower(item.Key + " " + name + " " + desc + " " + item.Default)
			if !strings.Contains(blob, filter) {
				continue
			}
		}
		rows = append(rows, row{item: item, name: name, desc: desc})
	}
	if len(rows) == 0 {
		return 0
	}

	title := secTitle(sec, a.lang)
	a.sectionHeader(c, t, sec.ID, fmt.Sprintf("%s · %d", title, len(rows)), "")

	ui.Row(c).Padding(6, 8).Background(t.Background).Children(func() {
		ui.Text(c, L.Param).Width(200).FontSize(a.fs(11)).TextColor(t.TextMuted)
		ui.Text(c, L.Description).Grow(1).FontSize(a.fs(11)).TextColor(t.TextMuted)
		ui.Text(c, L.Value).Width(250).FontSize(a.fs(11)).TextColor(t.TextMuted)
	})
	ui.Box(c).Height(1).Background(t.Border)

	for _, r := range rows {
		a.optionRow(c, t, r.item, r.name, r.desc)
	}
	return len(rows)
}

func (a *app) optionRow(c *ui.Context, t *ui.Theme, item options.Option, name, desc string) {
	L := a.L()
	a.ensureOptMaps()
	if _, ok := a.optVals[item.Key]; !ok {
		a.optVals[item.Key] = item.Default
	}
	dirty := a.optOrig[item.Key] != a.optVals[item.Key]
	bg := t.Surface
	if dirty {
		bg = a.selectedBG()
	}

	ui.Row(c).Gap(10).AlignItems(ui.Center).Padding(7, 8).Background(bg).Children(func() {
		ui.Column(c).Width(200).Gap(1).Children(func() {
			ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
				if dirty {
					ui.Box(c).Width(6).Height(6).Radius(99).Background(t.Accent)
				} else {
					ui.Box(c).Width(6).Height(6)
				}
				ui.Text(c, name).FontSize(a.fs(12)).FontWeight(600).SingleLine().Grow(1)
			})
			ui.Text(c, item.Key).FontSize(a.fs(10)).TextColor(t.TextMuted).SingleLine()
		})

		d := desc
		if item.Default != "" {
			d = d + "  ·  " + L.Default + " " + item.Default
		}
		ui.Text(c, d).FontSize(a.fs(11)).TextColor(t.TextMuted).SingleLine().Grow(1)

		ui.Row(c).Width(250).AlignItems(ui.Center).Children(func() {
			switch item.Kind {
			case options.KindFlag:
				on := strings.EqualFold(a.optVals[item.Key], "true")
				before := on
				ui.Switch(c, &on)
				if on != before {
					a.optVals[item.Key] = fmt.Sprintf("%t", on)
				}
			case options.KindEnum:
				val := a.optVals[item.Key]
				ui.Select(c, &val, item.Choices).FillWidth()
				if val != a.optVals[item.Key] {
					a.optVals[item.Key] = val
				}
			default:
				val := a.optVals[item.Key]
				ui.TextInput(c, &val).FillWidth()
				if val != a.optVals[item.Key] {
					a.optVals[item.Key] = val
				}
			}
		})
	})
	ui.Box(c).Height(1).Background(t.Border)
}

func cloneMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
