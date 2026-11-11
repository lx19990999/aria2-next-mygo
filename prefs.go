package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/egoist/mygo/ui"
)

const prefsFileName = "aria2-next-mygo.conf"

var (
	fontListOnce sync.Once
	fontList     []string
)

var accentPresets = []struct {
	Name string
	Hex  string
}{
	{"青绿", "#0d9488"},
	{"蓝色", "#2563eb"},
	{"靛蓝", "#4f46e5"},
	{"玫红", "#e11d48"},
	{"琥珀", "#d97706"},
	{"翠绿", "#059669"},
	{"石板", "#475569"},
}

func prefsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config")
	return filepath.Join(dir, prefsFileName), nil
}

func systemFontNames() []string {
	fontListOnce.Do(func() {
		fams := ui.FontFamilies()
		fontList = make([]string, 0, len(fams)+1)
		fontList = append(fontList, "system-ui")
		seen := map[string]bool{"system-ui": true}
		for _, f := range fams {
			if f.Name == "" || seen[f.Name] {
				continue
			}
			seen[f.Name] = true
			fontList = append(fontList, f.Name)
		}
	})
	return fontList
}

func colorToHex(c ui.Color) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func parseAccent(s string) ui.Color {
	fallback := ui.Hex("#0d9488")
	s = strings.TrimSpace(s)
	if s == "" {
		return fallback
	}
	if !strings.HasPrefix(s, "#") {
		s = "#" + s
	}
	var out ui.Color
	ok := true
	func() {
		defer func() {
			if recover() != nil {
				ok = false
			}
		}()
		out = ui.Hex(s)
	}()
	if !ok {
		return fallback
	}
	return out
}

func (a *app) defaultPrefs() {
	if a.lang == "" {
		a.lang = "zh"
	}
	if a.uiScale == 0 {
		a.uiScale = 1
	}
	if a.fontSize == 0 {
		a.fontSize = 13
	}
	if a.fontFamily == "" {
		a.fontFamily = "system-ui"
	}
	if a.accentHex == "" {
		a.accentHex = "#0d9488"
	}
	// Never rewrite editor text here — view() calls this every frame.
	if a.accentColor == (ui.Color{}) {
		a.accentColor = parseAccent(a.accentHex)
		a.accentHex = colorToHex(a.accentColor)
	}
	if a.scaleText == "" {
		a.scaleText = formatScale(a.uiScale)
	}
	if a.fontSizeText == "" {
		a.fontSizeText = formatFontSize(a.fontSize)
	}
	if a.accentHexText == "" {
		a.accentHexText = a.accentHex
	}
}

func (a *app) loadPrefs() {
	a.defaultPrefs()
	path, err := prefsPath()
	if err != nil {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		switch k {
		case "lang":
			if v == "zh" || v == "en" {
				a.lang = v
			}
		case "theme":
			a.themeDark = v == "dark"
		case "ui-scale":
			if x, err := strconv.ParseFloat(v, 32); err == nil && x >= 0.5 && x <= 3 {
				a.uiScale = float32(x)
			}
		case "font-family":
			if v != "" {
				a.fontFamily = v
			}
		case "font-size":
			if x, err := strconv.ParseFloat(v, 32); err == nil && x >= 9 && x <= 28 {
				a.fontSize = float32(x)
			}
		case "accent":
			a.accentHex = v
			a.accentColor = parseAccent(v)
		case "engine-path":
			a.enginePath = v
		case "engine-conf":
			a.engineConf = v
		case "engine-auto-start":
			a.engineAutoStart = v == "true" || v == "1" || v == "yes"
		case "engine-stop-on-exit":
			a.engineStopOnExit = v == "true" || v == "1" || v == "yes"
		case "rpc-url":
			if v != "" {
				a.rpcURL = v
			}
		case "rpc-secret-enc":
			if plain, err := decryptSecret(v); err == nil {
				a.rpcSecret = plain
			}
		case "rpc-secret":
			// Legacy plaintext; keep for migration and re-encrypt on save.
			if v != "" && a.rpcSecret == "" {
				a.rpcSecret = v
			}
		}
	}
	a.scaleText = formatScale(a.uiScale)
	a.fontSizeText = formatFontSize(a.fontSize)
	a.accentHexText = a.accentHex
}

func (a *app) savePrefs() {
	if err := a.writePrefs(); err != nil {
		a.flash("保存配置失败：" + err.Error())
		return
	}
	a.flash("配置已保存")
}

func (a *app) writePrefs() error {
	a.accentHex = colorToHex(a.accentColor)
	a.accentHexText = a.accentHex
	path, err := prefsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	theme := "light"
	if a.themeDark {
		theme = "dark"
	}
	autoStart := "false"
	if a.engineAutoStart {
		autoStart = "true"
	}
	stopOnExit := "false"
	if a.engineStopOnExit {
		stopOnExit = "true"
	}
	secretEnc, err := encryptSecret(a.rpcSecret)
	if err != nil {
		return fmt.Errorf("加密 RPC 密钥失败：%w", err)
	}
	if a.rpcURL == "" {
		a.rpcURL = defaultRPCURL
	}
	body := strings.Join([]string{
		"# aria2-next-mygo UI preferences",
		"lang=" + a.lang,
		"theme=" + theme,
		fmt.Sprintf("ui-scale=%.1f", a.uiScale),
		"font-family=" + a.fontFamily,
		fmt.Sprintf("font-size=%.0f", a.fontSize),
		"accent=" + a.accentHex,
		"rpc-url=" + a.rpcURL,
		"rpc-secret-enc=" + secretEnc,
		"engine-path=" + a.enginePath,
		"engine-conf=" + a.engineConf,
		"engine-auto-start=" + autoStart,
		"engine-stop-on-exit=" + stopOnExit,
		"",
	}, "\n")
	return os.WriteFile(path, []byte(body), 0o644)
}

func (a *app) bumpScale(delta float32) {
	a.uiScale += delta
	if a.uiScale < 0.5 {
		a.uiScale = 0.5
	}
	if a.uiScale > 3 {
		a.uiScale = 3
	}
	a.uiScale = float32(int(a.uiScale*10+0.5)) / 10
	a.scaleText = formatScale(a.uiScale)
	a.savePrefs()
}

func (a *app) bumpFontSize(delta float32) {
	a.fontSize += delta
	if a.fontSize < 9 {
		a.fontSize = 9
	}
	if a.fontSize > 28 {
		a.fontSize = 28
	}
	a.fontSize = float32(int(a.fontSize + 0.5))
	a.fontSizeText = formatFontSize(a.fontSize)
	a.savePrefs()
}
