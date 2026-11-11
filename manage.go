package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func (a *app) manageBody(c *ui.Context, t *ui.Theme) {
	L := a.L()

	ui.Column(c).Grow(1).Children(func() {
		ui.Row(c).Padding(10, 14).Gap(10).AlignItems(ui.Center).Background(t.Surface).Children(func() {
			ui.Text(c, L.ManageSettings).FontSize(a.fs(15)).FontWeight(600).Grow(1)
			if ui.PrimaryButton(c, L.Save).Padding(5, 12).Clicked() {
				a.savePrefs()
			}
		})
		ui.Box(c).Height(1).Background(t.Border)

		ui.Scroll(c).Grow(1).Children(func() {
			ui.Column(c).Padding(12, 16).Gap(14).Children(func() {
				ui.Text(c, L.ManageHint).FontSize(a.fs(12)).TextColor(t.TextMuted)

				a.uiPrefRow(c, t, L.EngineStatus, func() {
					ui.Text(c, a.engineStatusText()).FontSize(a.fs(13)).FontWeight(600)
					if !a.connected && a.connErr != "" {
						ui.Text(c, a.connErr).FontSize(a.fs(11)).TextColor(t.Danger)
					} else if a.engineAlive() && !a.connected {
						ui.Text(c, "进程已拉起；若标题仍是未连接，点「刷新」重试 RPC。").FontSize(a.fs(11)).TextColor(t.TextMuted)
					}
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						busy := a.engineBusy
						if ui.PrimaryButton(c, L.StartEngine).Disabled(busy || a.connected || a.engineAlive()).Padding(5, 12).Clicked() {
							a.runEngineAction(L.StartEngine, a.startManagedEngine)
						}
						if ui.Button(c, L.StopEngine).Disabled(busy || (!a.engineManaged() && !a.connected && !a.rpcPortOpen())).Padding(5, 12).Clicked() {
							a.runEngineAction(L.StopEngine, a.stopManagedEngine)
						}
						if ui.Button(c, L.Refresh).Padding(5, 12).Clicked() {
							a.rpcPaused = false
							go a.refresh()
						}
					})
				})

				a.uiPrefRow(c, t, L.EngineExecutable, func() {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						ui.TextInput(c, &a.enginePath).Grow(1).MinWidth(0).Placeholder(L.EngineExecutableHint)
						if ui.Button(c, L.Browse).Padding(4, 10).Clicked() {
							a.browseEnginePath()
						}
						if ui.Button(c, L.ClearPath).Disabled(strings.TrimSpace(a.enginePath) == "").Padding(4, 10).Clicked() {
							a.enginePath = ""
							a.savePrefs()
						}
					})
				})

				a.uiPrefRow(c, t, L.EngineConfig, func() {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						ui.TextInput(c, &a.engineConf).Grow(1).MinWidth(0).Placeholder(L.EngineConfigHint)
						if ui.Button(c, L.Browse).Padding(4, 10).Clicked() {
							a.browseEngineConf()
						}
						if ui.Button(c, L.ClearPath).Disabled(strings.TrimSpace(a.engineConf) == "").Padding(4, 10).Clicked() {
							a.engineConf = ""
							a.savePrefs()
						}
					})
				})

				a.uiPrefRow(c, t, L.ManageOptions, func() {
					prevAuto := a.engineAutoStart
					prevStop := a.engineStopOnExit
					ui.Checkbox(c, &a.engineAutoStart, L.AutoStartEngine)
					ui.Checkbox(c, &a.engineStopOnExit, L.StopEngineOnExit)
					if a.engineAutoStart != prevAuto || a.engineStopOnExit != prevStop {
						a.savePrefs()
					}
				})

				a.uiPrefRow(c, t, L.RPCTitle, func() {
					ui.Text(c, L.ManageRPCHint).FontSize(a.fs(12)).TextColor(t.TextMuted)
					ui.Textf(c, "%s  %s", L.Address, a.rpcURL).FontSize(a.fs(12)).TextColor(t.TextMuted)
					ui.Textf(c, "%s  %s", L.Secret, maskSecret(a.rpcSecret)).FontSize(a.fs(12)).TextColor(t.TextMuted)
					if ui.Button(c, L.SyncRPCFromConf).Padding(4, 10).Clicked() {
						if err := a.syncRPCFromEngineConfLocal(); err != nil {
							a.flash(err.Error())
						} else {
							a.savePrefs()
							a.flash(L.SyncRPCFromConfOK)
						}
					}
				})
			})
		})
	})
}

func maskSecret(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "(empty)"
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

func (a *app) runEngineAction(label string, fn func() error) {
	if a.engineBusy {
		return
	}
	a.engineBusy = true
	go func() {
		err := fn()
		a.apply(func() {
			a.engineBusy = false
			if err != nil {
				a.flash(label + "失败：" + err.Error())
			} else {
				a.flash(label + "成功")
			}
		})
		time.Sleep(300 * time.Millisecond)
		a.refresh()
	}()
}

func (a *app) browseEnginePath() {
	go func() {
		filters := []mygo.FileFilter{}
		if runtime.GOOS == "windows" {
			filters = []mygo.FileFilter{{Name: "Executable", Extensions: []string{"exe"}}}
		}
		paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
			Parent:      a.win,
			Title:       a.L().EngineExecutable,
			DefaultPath: a.enginePath,
			Filters:     filters,
		})
		if err != nil || len(paths) == 0 {
			return
		}
		abs, _ := filepath.Abs(paths[0])
		a.apply(func() {
			a.enginePath = abs
			a.savePrefs()
		})
	}()
}

func (a *app) browseEngineConf() {
	go func() {
		paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
			Parent:      a.win,
			Title:       a.L().EngineConfig,
			DefaultPath: a.engineConf,
			Filters: []mygo.FileFilter{
				{Name: "Config", Extensions: []string{"conf", "cfg", "txt"}},
			},
		})
		if err != nil || len(paths) == 0 {
			return
		}
		abs, _ := filepath.Abs(paths[0])
		a.apply(func() {
			a.engineConf = abs
			_ = a.syncRPCFromEngineConfLocal()
			a.savePrefs()
		})
	}()
}

func (a *app) maybeAutoStartEngine() {
	if !a.engineAutoStart {
		return
	}
	go func() {
		time.Sleep(300 * time.Millisecond)
		finish := func(err error) {
			// Clear pause before refresh so the poll does not race with async UI apply.
			a.rpcPaused = false
			a.apply(func() {
				if err != nil {
					a.flash(a.L().StartEngine + "失败：" + err.Error())
				} else {
					a.flash(a.L().StartEngine + "成功")
				}
			})
			a.refresh()
		}

		if strings.TrimSpace(a.enginePath) == "" || strings.TrimSpace(a.engineConf) == "" {
			finish(fmt.Errorf("未配置可执行文件或配置文件，已跳过自动启动"))
			return
		}
		if err := a.syncRPCFromEngineConfLocal(); err != nil {
			finish(fmt.Errorf("读取引擎配置失败：%w", err))
			return
		}
		_ = a.writePrefs()

		if a.engineAlive() {
			finish(nil)
			return
		}

		// If the configured conf port is already serving, just connect to it.
		if a.rpcPortOpen() {
			if _, err := a.client().GetVersion(); err == nil {
				finish(nil)
				return
			}
			// Port busy but not our RPC — free it, then launch configured binary.
			_ = a.client().Shutdown()
			_ = a.client().ForceShutdown()
			deadline := time.Now().Add(6 * time.Second)
			for time.Now().Before(deadline) {
				if !a.rpcPortOpen() {
					break
				}
				time.Sleep(150 * time.Millisecond)
			}
			if a.rpcPortOpen() {
				finish(fmt.Errorf("RPC 端口仍被占用，无法启动托管引擎；请先停止旧实例"))
				return
			}
		}

		finish(a.startManagedEngine())
	}()
}

func (a *app) maybeStopEngineOnExit() {
	if !a.engineStopOnExit || !a.engineManaged() {
		return
	}
	_ = a.stopManagedEngine()
}
