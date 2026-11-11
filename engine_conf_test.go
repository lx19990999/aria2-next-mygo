package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncRPCFromEngineConf(t *testing.T) {
	dir := t.TempDir()
	conf := filepath.Join(dir, "aria2next.conf")
	body := "enable-rpc=true\nrpc-listen-port=6800\nrpc-secret=DiRoKMa7b1wT\n"
	if err := os.WriteFile(conf, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	a := &app{
		rpcURL:     "http://127.0.0.1:17800/jsonrpc",
		rpcSecret:  "old",
		engineConf: conf,
	}
	if err := a.syncRPCFromEngineConf(); err != nil {
		t.Fatal(err)
	}
	if a.rpcURL != "http://127.0.0.1:6800/jsonrpc" {
		t.Fatalf("url %q", a.rpcURL)
	}
	if a.rpcSecret != "DiRoKMa7b1wT" {
		t.Fatalf("secret %q", a.rpcSecret)
	}
}

func TestSyncRPCFromEngineConfKeepsRemote(t *testing.T) {
	dir := t.TempDir()
	conf := filepath.Join(dir, "aria2next.conf")
	body := "enable-rpc=true\nrpc-listen-port=6800\nrpc-secret=local-secret\n"
	if err := os.WriteFile(conf, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	a := &app{
		rpcURL:     "http://192.168.10.109:16800/jsonrpc",
		rpcSecret:  "DiRoKMa7b1wT",
		engineConf: conf,
	}
	if err := a.syncRPCFromEngineConf(); err != nil {
		t.Fatal(err)
	}
	if a.rpcURL != "http://192.168.10.109:16800/jsonrpc" {
		t.Fatalf("remote url rewritten to %q", a.rpcURL)
	}
	if a.rpcSecret != "DiRoKMa7b1wT" {
		t.Fatalf("remote secret rewritten to %q", a.rpcSecret)
	}
}

func TestSyncRPCFromEngineConfLocalForcesLoopback(t *testing.T) {
	dir := t.TempDir()
	conf := filepath.Join(dir, "aria2next.conf")
	body := "enable-rpc=true\nrpc-listen-port=6800\nrpc-secret=local-secret\n"
	if err := os.WriteFile(conf, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	a := &app{
		rpcURL:     "http://192.168.10.109:16800/jsonrpc",
		rpcSecret:  "DiRoKMa7b1wT",
		engineConf: conf,
	}
	if err := a.syncRPCFromEngineConfLocal(); err != nil {
		t.Fatal(err)
	}
	if a.rpcURL != "http://127.0.0.1:6800/jsonrpc" {
		t.Fatalf("url %q", a.rpcURL)
	}
	if a.rpcSecret != "local-secret" {
		t.Fatalf("secret %q", a.rpcSecret)
	}
}
