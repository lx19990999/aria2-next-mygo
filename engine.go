package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type engineManager struct {
	mu sync.Mutex

	cmd     *exec.Cmd
	started bool // true if this GUI launched the current process
}

func (a *app) engineAlive() bool {
	a.engine.mu.Lock()
	defer a.engine.mu.Unlock()
	return a.engine.cmd != nil && a.engine.cmd.Process != nil && a.engine.cmd.ProcessState == nil
}

func (a *app) engineManaged() bool {
	a.engine.mu.Lock()
	defer a.engine.mu.Unlock()
	return a.engine.started && a.engine.cmd != nil && a.engine.cmd.Process != nil && a.engine.cmd.ProcessState == nil
}

func (a *app) rpcPortOpen() bool {
	host, port := rpcHostPort(a.rpcURL)
	if host == "" || port == "" {
		return false
	}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 400*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func rpcHostPort(raw string) (host, port string) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" {
		return "127.0.0.1", "17800"
	}
	h := u.Hostname()
	p := u.Port()
	if h == "" {
		h = "127.0.0.1"
	}
	if p == "" {
		p = "17800"
	}
	return h, p
}

func (a *app) validateEnginePaths() error {
	exe := strings.TrimSpace(a.enginePath)
	conf := strings.TrimSpace(a.engineConf)
	if exe == "" {
		return fmt.Errorf("未设置 aria2-next 可执行文件路径")
	}
	if conf == "" {
		return fmt.Errorf("未设置配置文件路径")
	}
	if st, err := os.Stat(exe); err != nil || st.IsDir() {
		return fmt.Errorf("可执行文件不存在：%s", exe)
	}
	if st, err := os.Stat(conf); err != nil || st.IsDir() {
		return fmt.Errorf("配置文件不存在：%s", conf)
	}
	return nil
}

func isLoopbackHost(host string) bool {
	h := strings.ToLower(strings.TrimSpace(host))
	return h == "127.0.0.1" || h == "localhost" || h == "::1" || h == "[::1]"
}

// syncRPCFromEngineConf reads rpc-listen-port / rpc-secret from the managed
// aria2 conf so the GUI does not keep probing a hardcoded old port.
// forceLocal rewrites the URL to 127.0.0.1 (managed engine). Otherwise the
// port/secret are only applied when the current RPC host is already loopback,
// so a remote endpoint like 192.168.x.x:16800 is never overwritten.
func (a *app) syncRPCFromEngineConf() error {
	return a.syncRPCFromEngineConfMode(false)
}

func (a *app) syncRPCFromEngineConfLocal() error {
	return a.syncRPCFromEngineConfMode(true)
}

func (a *app) syncRPCFromEngineConfMode(forceLocal bool) error {
	conf := strings.TrimSpace(a.engineConf)
	if conf == "" {
		return nil
	}
	f, err := os.Open(conf)
	if err != nil {
		return err
	}
	defer f.Close()

	port := ""
	secret := ""
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
		case "rpc-listen-port":
			port = v
		case "rpc-secret":
			secret = v
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}

	host, _ := rpcHostPort(a.rpcURL)
	if forceLocal {
		host = "127.0.0.1"
	} else if !isLoopbackHost(host) {
		return nil
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port != "" {
		a.rpcURL = fmt.Sprintf("http://%s:%s/jsonrpc", host, port)
	}
	if secret != "" {
		a.rpcSecret = secret
	}
	return nil
}

func (a *app) startManagedEngine() error {
	if err := a.validateEnginePaths(); err != nil {
		return err
	}
	if err := a.syncRPCFromEngineConfLocal(); err != nil {
		return fmt.Errorf("读取引擎配置失败：%w", err)
	}
	_ = a.writePrefs()

	if a.engineAlive() {
		return fmt.Errorf("本应用已在托管一个引擎进程")
	}
	if a.connected || a.rpcPortOpen() {
		return fmt.Errorf("检测到 RPC 端口已有服务在监听，请先连接现有实例或更换端口")
	}

	exe := strings.TrimSpace(a.enginePath)
	conf := strings.TrimSpace(a.engineConf)
	cmd := exec.Command(exe, "--conf-path="+conf)
	cmd.Dir = filepath.Dir(conf)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	configureEngineCmd(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	a.engine.mu.Lock()
	a.engine.cmd = cmd
	a.engine.started = true
	a.engine.mu.Unlock()

	go func() {
		_ = cmd.Wait()
		a.engine.mu.Lock()
		if a.engine.cmd == cmd {
			a.engine.cmd = nil
			a.engine.started = false
		}
		a.engine.mu.Unlock()
	}()

	// Give the engine a moment to bind RPC.
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		if a.rpcPortOpen() {
			return nil
		}
		if !a.engineAlive() {
			return fmt.Errorf("引擎进程已退出，请检查配置与日志")
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("引擎已启动，但 RPC 端口尚未就绪")
}

func (a *app) stopManagedEngine() error {
	owned := a.engineManaged()
	if a.connected || a.rpcPortOpen() {
		if err := a.client().Shutdown(); err != nil {
			_ = a.client().ForceShutdown()
		}
		// Wait briefly for port to close.
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if !a.rpcPortOpen() {
				break
			}
			time.Sleep(150 * time.Millisecond)
		}
	}
	if owned {
		a.engine.mu.Lock()
		cmd := a.engine.cmd
		a.engine.mu.Unlock()
		if cmd != nil && cmd.Process != nil && cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_, _ = cmd.Process.Wait()
		}
		a.engine.mu.Lock()
		a.engine.cmd = nil
		a.engine.started = false
		a.engine.mu.Unlock()
	} else if !a.connected && !a.rpcPortOpen() && !a.engineAlive() {
		return fmt.Errorf("没有可停止的托管进程")
	}
	return nil
}

func (a *app) engineStatusText() string {
	L := a.L()
	switch {
	case a.connected:
		if a.engineManaged() {
			return L.EngineStatusManagedConnected
		}
		return L.EngineStatusExternalConnected
	case a.engineAlive():
		// Port already accepting TCP but JSON-RPC still failing.
		if a.rpcPortOpen() && a.connErr != "" {
			return L.EngineStatusRPCFailed
		}
		return L.EngineStatusStarting
	case a.rpcPortOpen():
		return L.EngineStatusPortBusy
	default:
		return L.EngineStatusStopped
	}
}

