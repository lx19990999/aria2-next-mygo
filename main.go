package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"aria2-next-mygo/rpc"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

const (
	defaultRPCURL    = "http://127.0.0.1:17800/jsonrpc"
	defaultRPCSecret = "DiRoKMa7b1wT"
	pollInterval     = time.Second
)

type app struct {
	mu sync.Mutex

	win *mygo.Window

	rpcURL    string
	rpcSecret string
	newURI    string

	connected bool
	connErr   string
	version   rpc.Version
	stat      rpc.GlobalStat

	section     string // active | waiting | stopped | manage | settings | ui | about
	active      []rpc.Task
	waiting     []rpc.Task
	stopped     []rpc.Task
	selectedGID string
	statusMsg   string
	statusAt    time.Time

	settingsTab    int // unused; kept for env compat
	settingsAnchor string
	pendingScroll  string
	optFilter      string
	optVals        map[string]string
	optOrig        map[string]string
	optsLoaded     bool
	optsLoading    bool

	lang         string // zh | en
	uiScale      float32
	fontSize     float32
	fontFamily   string
	themeDark    bool
	accentHex         string
	accentHexText     string
	accentColor       ui.Color
	scaleText         string
	fontSizeText      string
	scaleFocused      bool
	fontSizeFocused   bool
	accentHexFocused  bool
	lastTitle         string

	brewCmd   string
	dockerCmd string

	enginePath       string
	engineConf       string
	engineAutoStart  bool
	engineStopOnExit bool
	engineBusy       bool
	engine           engineManager
	rpcPaused        bool // user disconnected; skip auto-reconnect until Reconnect
	refreshMu        sync.Mutex
	rpcFailStreak    int
}

func (a *app) client() *rpc.Client {
	return rpc.New(a.rpcURL, a.rpcSecret)
}

func (a *app) currentTasks() []rpc.Task {
	switch a.section {
	case "waiting":
		return a.waiting
	case "stopped":
		return a.stopped
	default:
		return a.active
	}
}

func (a *app) startPolling() {
	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		a.refresh()
		for range ticker.C {
			a.refresh()
		}
	}()
}

func (a *app) refresh() {
	// Skip overlapping polls (ticker + manual reconnect / options load).
	if !a.refreshMu.TryLock() {
		return
	}
	defer a.refreshMu.Unlock()

	if a.rpcPaused {
		a.apply(func() {
			a.connected = false
			a.connErr = ""
			a.rpcFailStreak = 0
			a.updateWindowTitle()
		})
		return
	}
	c := a.client()
	ver, err := c.GetVersion()
	if err != nil {
		a.noteRPCFailure(err)
		return
	}
	stat, err := c.GetGlobalStat()
	if err != nil {
		a.noteRPCFailure(err)
		return
	}
	active, err1 := c.TellActive()
	waiting, err2 := c.TellWaiting(0, 200)
	stopped, err3 := c.TellStopped(0, 200)
	if err1 != nil || err2 != nil || err3 != nil {
		err = err1
		if err == nil {
			err = err2
		}
		if err == nil {
			err = err3
		}
		a.noteRPCFailure(err)
		return
	}
	a.apply(func() {
		a.connected = true
		a.connErr = ""
		a.rpcFailStreak = 0
		a.version = ver
		a.stat = stat
		a.active = active
		a.waiting = waiting
		a.stopped = stopped
		a.updateWindowTitle()
	})
}

// noteRPCFailure only flips to disconnected after consecutive failures so a
// single transient EOF does not flicker the UI every second.
func (a *app) noteRPCFailure(err error) {
	a.apply(func() {
		a.rpcFailStreak++
		a.connErr = err.Error()
		if a.rpcFailStreak >= 2 || !a.connected {
			a.connected = false
			a.updateWindowTitle()
		}
	})
}

func (a *app) updateWindowTitle() {
	if a.win == nil {
		return
	}
	title := a.windowTitle()
	if title == a.lastTitle {
		return
	}
	a.lastTitle = title
	a.win.SetTitle(title)
}

func (a *app) windowTitle() string {
	L := a.L()
	product := a.version.Product
	if product == "" {
		product = "aria2-next"
	}
	engine := product
	if a.version.Version != "" {
		engine = product + " " + a.version.Version
	}
	status := L.Disconnected
	if a.connected {
		host, _ := rpcHostPort(a.rpcURL)
		if isLoopbackHost(host) {
			status = L.LocalEngine
		} else {
			status = L.RemoteEngine
		}
	}
	return fmt.Sprintf("aria2-next-mygo %s · %s · %s", appVersion, engine, status)
}

func (a *app) apply(fn func()) {
	if a.win == nil {
		a.mu.Lock()
		fn()
		a.mu.Unlock()
		return
	}
	a.win.Update(fn)
}

func (a *app) flash(msg string) {
	a.statusMsg = msg
	a.statusAt = time.Now()
}

func (a *app) runRPC(label string, fn func(*rpc.Client) error) {
	go func() {
		err := fn(a.client())
		a.apply(func() {
			if err != nil {
				a.flash(label + "失败：" + err.Error())
			} else {
				a.flash(label + "成功")
			}
		})
		a.refresh()
	}()
}

func (a *app) view(c *ui.Context) {
	if a.section == "" {
		a.section = "active"
	}
	t := applyTheme(c, a)
	c.Root().Background(t.Background)
	L := a.L()
	a.defaultPrefs()

	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		// Compact sidebar — speeds live only here (not duplicated in footer).
		ui.Column(c).Width(208).Shrink(0).Background(a.chromeBG()).Padding(14, 12).Gap(10).Children(func() {
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Textf(c, "↓ %s", rpc.FormatSpeed(a.stat.DownloadSpeed)).FontSize(a.fs(12)).FontWeight(600).TextColor(t.Accent).SingleLine()
				ui.Textf(c, "↑ %s", rpc.FormatSpeed(a.stat.UploadSpeed)).FontSize(a.fs(12)).FontWeight(500).TextColor(t.TextMuted).SingleLine()
			})

			ui.Box(c).Height(1).Background(t.Border)

			ui.Column(c).Gap(2).Children(func() {
				a.navButton(c, t, "active", L.Downloading, nz(a.stat.NumActive))
				a.navButton(c, t, "waiting", L.Waiting, nz(a.stat.NumWaiting))
				a.navButton(c, t, "stopped", L.Stopped, nz(a.stat.NumStopped))
				ui.Box(c).Height(6)
				a.navButton(c, t, "manage", L.ManageSettings, "")
				a.navButton(c, t, "settings", L.EngineSettings, "")
				a.navButton(c, t, "ui", L.UISettings, "")
				a.navButton(c, t, "about", L.About, "")
			})

			ui.Spacer(c)
			if a.statusMsg != "" && time.Since(a.statusAt) < 5*time.Second {
				ui.Text(c, a.statusMsg).FontSize(a.fs(11)).TextColor(t.TextMuted)
			} else if !a.connected && a.connErr != "" {
				ui.Text(c, a.connErr).FontSize(a.fs(11)).TextColor(t.Danger)
			}
		})

		ui.Box(c).Width(1).Background(t.Border)

		ui.Column(c).Grow(1).MinWidth(0).Background(t.Surface).Children(func() {
			switch a.section {
			case "manage":
				a.manageBody(c, t)
			case "settings":
				a.settingsBody(c, t)
			case "ui":
				a.uiSettingsBody(c, t)
			case "about":
				a.aboutBody(c, t)
			default:
				a.header(c, t)
				ui.Box(c).Height(1).Background(t.Border)
				a.tasksBody(c, t)
			}
		})
	})
}

func (a *app) navButton(c *ui.Context, t *ui.Theme, id, label, count string) {
	selected := a.section == id
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
		ui.Row(c).Grow(1).Padding(7, 10).Gap(8).AlignItems(ui.Center).Children(func() {
			labelColor := t.Text
			weight := 500
			if selected {
				labelColor = t.Accent
				weight = 600
			}
			ui.Text(c, label).FontSize(a.fs(13)).FontWeight(weight).TextColor(labelColor).Grow(1).SingleLine()
			if count != "" {
				ui.Text(c, count).FontSize(a.fs(11)).TextColor(t.TextMuted)
			}
		})
	})
	if row.Clicked() {
		a.section = id
		if id == "settings" && !a.optsLoaded && !a.optsLoading {
			a.loadGlobalOptions()
		}
	}
}

func (a *app) header(c *ui.Context, t *ui.Theme) {
	L := a.L()
	title := L.DownloadTasks
	switch a.section {
	case "waiting":
		title = L.WaitingQueue
	case "stopped":
		title = L.StoppedTitle
	}
	ui.Row(c).Padding(10, 14).Gap(8).AlignItems(ui.Center).Background(t.Surface).Children(func() {
		ui.Text(c, title).FontSize(a.fs(15)).FontWeight(600).Grow(1)
		if ui.Button(c, L.Refresh).Padding(5, 10).Clicked() {
			go a.refresh()
		}
		if a.section == "stopped" {
			if ui.Button(c, L.ClearStopped).Padding(5, 10).Clicked() {
				a.runRPC(L.ClearStopped, func(client *rpc.Client) error {
					return client.PurgeDownloadResult()
				})
			}
		}
		disabled := a.selectedGID == ""
		if ui.Button(c, L.Pause).Disabled(disabled).Padding(5, 10).Clicked() {
			gid := a.selectedGID
			a.runRPC(L.Pause, func(client *rpc.Client) error { return client.Pause(gid) })
		}
		if ui.Button(c, L.Resume).Disabled(disabled).Padding(5, 10).Clicked() {
			gid := a.selectedGID
			a.runRPC(L.Resume, func(client *rpc.Client) error { return client.Unpause(gid) })
		}
		if ui.Button(c, L.Delete).Disabled(disabled).Padding(5, 10).Clicked() {
			gid := a.selectedGID
			a.runRPC(L.Delete, func(client *rpc.Client) error {
				if a.section == "stopped" {
					return client.RemoveDownloadResult(gid)
				}
				return client.Remove(gid)
			})
		}
	})
}

func (a *app) tasksBody(c *ui.Context, t *ui.Theme) {
	L := a.L()
	ui.Column(c).Grow(1).Children(func() {
		ui.Row(c).Padding(10, 14).Gap(8).AlignItems(ui.Center).Background(t.Surface).Children(func() {
			ui.TextInput(c, &a.newURI).Placeholder(L.PasteLink).Grow(1)
			if ui.PrimaryButton(c, L.Add).Padding(6, 14).Clicked() {
				uri := strings.TrimSpace(a.newURI)
				if uri == "" {
					a.flash(L.PasteLink)
				} else {
					a.runRPC(L.Add, func(client *rpc.Client) error {
						_, err := client.AddURI(uri)
						return err
					})
					a.newURI = ""
				}
			}
		})
		ui.Box(c).Height(1).Background(t.Border)

		tasks := a.currentTasks()

		ui.Row(c).Padding(6, 14).Background(t.Background).Children(func() {
			ui.Text(c, L.Name).Grow(1).FontSize(a.fs(11)).TextColor(t.TextMuted)
			ui.Text(c, L.Progress).Width(160).FontSize(a.fs(11)).TextColor(t.TextMuted)
			ui.Text(c, L.Size).Width(130).FontSize(a.fs(11)).TextColor(t.TextMuted)
			ui.Text(c, L.Speed).Width(90).FontSize(a.fs(11)).TextColor(t.TextMuted)
			ui.Text(c, L.Conns).Width(48).FontSize(a.fs(11)).TextColor(t.TextMuted)
			ui.Text(c, L.Status).Width(64).FontSize(a.fs(11)).TextColor(t.TextMuted)
		})
		ui.Box(c).Height(1).Background(t.Border)

		ui.Scroll(c).Grow(1).Children(func() {
			ui.Column(c).Gap(0).Children(func() {
				if len(tasks) == 0 {
					ui.Column(c).Padding(48).Gap(6).Center().Children(func() {
						ui.Text(c, L.NoTasks).FontSize(a.fs(14)).FontWeight(600)
						ui.Text(c, L.NoTasksHint).FontSize(a.fs(12)).TextColor(t.TextMuted)
					})
					return
				}
				for i := range tasks {
					a.taskRow(c, t, tasks[i])
				}
			})
		})

		ui.Box(c).Height(1).Background(t.Border)
		ui.Row(c).Padding(7, 14).Background(t.Surface).Children(func() {
			ui.Textf(c, L.Items, len(tasks)).FontSize(a.fs(11)).TextColor(t.TextMuted).Grow(1)
		})
	})
}

func (a *app) taskRow(c *ui.Context, t *ui.Theme, task rpc.Task) {
	selected := a.selectedGID == task.GID
	progress := task.Progress()
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	bg := t.Surface
	if selected {
		bg = a.selectedBG()
	}

	row := ui.Row(c).Gap(12).AlignItems(ui.Center).Padding(10, 14).Background(bg)
	row.Children(func() {
		ui.Text(c, task.Name()).FontSize(a.fs(13)).FontWeight(500).SingleLine().Grow(1).MinWidth(0)
		ui.Column(c).Width(150).Gap(4).Children(func() {
			ui.Progress(c, progress).FillWidth().Height(5).Radius(2)
			ui.Textf(c, "%.0f%%", progress*100).FontSize(a.fs(11)).TextColor(t.TextMuted)
		})
		ui.Textf(c, "%s / %s",
			rpc.FormatBytes(rpc.ParseInt(task.CompletedLength)),
			rpc.FormatBytes(rpc.ParseInt(task.TotalLength)),
		).Width(128).FontSize(a.fs(12)).TextColor(t.TextMuted).SingleLine()
		ui.Text(c, rpc.FormatSpeed(task.DownloadSpeed)).Width(88).FontSize(a.fs(12)).FontWeight(500).TextColor(t.Text).SingleLine()
		ui.Text(c, nz(task.Connections)).Width(44).FontSize(a.fs(12)).TextColor(t.TextMuted)
		ui.Text(c, statusLabel(task.Status)).Width(60).FontSize(a.fs(12)).FontWeight(500).TextColor(statusColor(t, task.Status))
	})
	if row.Clicked() {
		a.selectedGID = task.GID
	}
	ui.Box(c).Height(1).Background(t.Border)
}

func statusLabel(status string) string {
	switch status {
	case "active":
		return "下载中"
	case "waiting":
		return "等待"
	case "paused":
		return "已暂停"
	case "complete":
		return "已完成"
	case "error":
		return "出错"
	case "removed":
		return "已删除"
	default:
		return status
	}
}

func statusColor(t *ui.Theme, status string) ui.Color {
	switch status {
	case "active":
		return t.Accent
	case "complete":
		return t.Success
	case "error":
		return t.Danger
	default:
		return t.TextMuted
	}
}

func nz(s string) string {
	if s == "" {
		return "0"
	}
	return s
}

func main() {
	section := "active"
	if s := strings.TrimSpace(os.Getenv("ARIA2_MYGO_SECTION")); s != "" {
		section = s
	}
	settingsTab := 0
	if s := strings.TrimSpace(os.Getenv("ARIA2_MYGO_SETTINGS_TAB")); s != "" {
		fmt.Sscanf(s, "%d", &settingsTab)
	}
	anchor := anchorConn
	if s := strings.TrimSpace(os.Getenv("ARIA2_MYGO_SETTINGS_ANCHOR")); s != "" {
		anchor = s
	} else if settingsTab > 0 {
		anchor = anchorConn
	}
	a := &app{
		rpcURL:         defaultRPCURL,
		rpcSecret:      defaultRPCSecret,
		section:        section,
		settingsTab:    settingsTab,
		settingsAnchor: anchor,
		pendingScroll:  anchor,
	}
	a.loadPrefs()
	// Keep GUI RPC endpoint aligned with the managed engine conf (port/secret).
	if strings.TrimSpace(a.engineConf) != "" {
		_ = a.syncRPCFromEngineConf()
	}
	// Persist RPC URL + encrypted secret into prefs (migrate from hardcoded defaults).
	_ = a.writePrefs()
	mygo.App.WhenReady(func() {
		a.win = mygo.NewWindow(mygo.WindowOptions{
			Title:     a.windowTitle(),
			Width:     1180,
			Height:    760,
			MinWidth:  960,
			MinHeight: 620,
			// Avoid restoring a broken previous window geometry.
			StateKey: "",
			Content:  ui.View(a.view),
		})
		a.lastTitle = a.windowTitle()
		a.win.OnClosed(func() {
			a.maybeStopEngineOnExit()
		})
		// Managed auto-start only applies to loopback RPC. A saved remote URL
		// (e.g. LAN host) must not be rewritten or blocked by local engine start.
		host, _ := rpcHostPort(a.rpcURL)
		autoStartLocal := a.engineAutoStart && isLoopbackHost(host)
		if autoStartLocal {
			a.rpcPaused = true
		}
		a.startPolling()
		if autoStartLocal {
			a.maybeAutoStartEngine()
		}
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
