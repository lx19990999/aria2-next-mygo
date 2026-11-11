package main

import (
	"strings"
	"testing"

	"aria2-next-mygo/rpc"

	"github.com/egoist/mygo/ui"
)

func TestCommitScaleAndFontFromText(t *testing.T) {
	a := &app{uiScale: 1, fontSize: 13, scaleText: "1.0", fontSizeText: "13"}
	a.scaleText = "1.7"
	if !a.commitScaleFromText() || a.uiScale != 1.7 {
		t.Fatalf("scale commit got %v text=%q", a.uiScale, a.scaleText)
	}
	a.fontSizeText = "18"
	if !a.commitFontSizeFromText() || a.fontSize != 18 {
		t.Fatalf("font commit got %v text=%q", a.fontSize, a.fontSizeText)
	}
	a.accentHexText = "#2563eb"
	a.commitAccentHexFromText()
	if a.accentHex != "#2563eb" {
		t.Fatalf("accent hex %q", a.accentHex)
	}
}

func TestWindowTitleIncludesVersions(t *testing.T) {
	a := &app{
		rpcURL:    "http://127.0.0.1:6800/jsonrpc",
		connected: true,
		version:   rpc.Version{Version: "3.1.0", Product: "aria2-next"},
	}
	got := a.windowTitle()
	want := "aria2-next-mygo " + appVersion + " · aria2-next 3.1.0 · 本机已连接"
	if got != want {
		t.Fatalf("title %q want %q", got, want)
	}
	a.rpcURL = "http://192.168.10.109:16800/jsonrpc"
	got = a.windowTitle()
	want = "aria2-next-mygo " + appVersion + " · aria2-next 3.1.0 · 远程已连接"
	if got != want {
		t.Fatalf("remote title %q want %q", got, want)
	}
	a.connected = false
	got = a.windowTitle()
	want = "aria2-next-mygo " + appVersion + " · aria2-next 3.1.0 · 未连接"
	if got != want {
		t.Fatalf("disconnected title %q want %q", got, want)
	}
}

func TestViewRendersTasksAndSettings(t *testing.T) {
	a := &app{
		rpcURL:    defaultRPCURL,
		rpcSecret: defaultRPCSecret,
		section:   "active",
		connected: true,
		version:   rpc.Version{Version: "3.1.0", Product: "aria2-next"},
		stat: rpc.GlobalStat{
			NumActive: "1", NumWaiting: "0", NumStopped: "2", DownloadSpeed: "1363148",
		},
		active: []rpc.Task{{
			GID: "abc", Status: "active", TotalLength: "1000", CompletedLength: "400",
			DownloadSpeed: "1363148", Connections: "8",
			Files: []rpc.File{{Path: `C:\Users\Administrator\Downloads\demo.iso`}},
		}},
	}
	tt := ui.NewTester(a.view, 1180, 760)
	for _, want := range []string{"下载任务", "demo.iso", "添加", "下载中"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q texts=%q", want, tt.Texts())
		}
	}
	if tt.HasText("aria2-next 3.1.0") || tt.HasText("本机已连接") {
		t.Fatalf("sidebar should not show version/status texts=%q", tt.Texts())
	}

	a.section = "settings"
	a.connected = true
	a.optsLoaded = true
	a.lang = "zh"
	a.uiScale = 1
	a.fontSize = 13
	a.fontFamily = "system-ui"
	a.accentHex = "#0d9488"
	a.accentColor = ui.Hex("#0d9488")
	a.optVals = map[string]string{"dir": `C:\Downloads`, "split": "16"}
	a.optOrig = map[string]string{"dir": `C:\Downloads`, "split": "16"}
	tt = ui.NewTester(a.view, 1180, 760)
	for _, want := range []string{"托管设置", "引擎设置", "界面设置", "重新连接", "断开连接", "RPC 连接", "基本", "下载路径"} {
		if !tt.HasText(want) {
			t.Fatalf("settings missing %q texts=%q", want, tt.Texts())
		}
	}

	a.section = "manage"
	a.enginePath = `C:\source\20261110\aria2-next-3.1.0-windows-x86_64.exe`
	a.engineConf = `C:\source\20261110\aria2next.conf`
	tt = ui.NewTester(a.view, 1180, 760)
	for _, want := range []string{"托管设置", "启动引擎", "停止引擎", "可执行文件", "配置文件", "自动启动", "清空"} {
		if !tt.HasText(want) {
			t.Fatalf("manage missing %q texts=%q", want, tt.Texts())
		}
	}

	a.section = "ui"
	tt = ui.NewTester(a.view, 1180, 760)
	for _, want := range []string{"界面设置", "界面缩放", "界面字体", "强调色", "关于"} {
		if !tt.HasText(want) {
			t.Fatalf("ui settings missing %q texts=%q", want, tt.Texts())
		}
	}

	a.section = "about"
	tt = ui.NewTester(a.view, 1180, 760)
	for _, want := range []string{"关于", "本应用", "aria2-next", "releases/latest", "复制", "Homebrew", "Docker", "github.com/egoist/mygo"} {
		if !tt.HasText(want) {
			t.Fatalf("about missing %q texts=%q", want, tt.Texts())
		}
	}
	for _, bad := range []string{"打开仓库", "打开最新 Release", "立即下载", "GitHub Latest"} {
		if tt.HasText(bad) {
			t.Fatalf("about still shows %q texts=%q", bad, tt.Texts())
		}
	}
	a.ensureAboutCommands()
	if !strings.Contains(a.brewCmd, "brew install") || a.dockerCmd == "" {
		t.Fatalf("copyable commands not initialized brew=%q docker=%q", a.brewCmd, a.dockerCmd)
	}
	if tt.HasText("从源码构建") {
		t.Fatalf("about should not mention building from source texts=%q", tt.Texts())
	}
}
