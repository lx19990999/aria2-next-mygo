package main

import "runtime"

// appVersion defaults to "dev". Release and local builds override it with:
//
//	go build -ldflags "-X main.appVersion=$(./scripts/git-version.sh)"
//
// or MYGO_LDFLAGS when using `go tool mygo build`. The version comes from the
// latest git tag (v1.2.3 → 1.2.3), not from editing this file.
var appVersion = "dev"

const (
	aria2NextLatestURL = "https://github.com/AnInsomniacy/aria2-next/releases/latest"
	aria2NextBrewTap   = "aninsomniacy/aria2-next/aria2-next"
	mygoRepoURL        = "https://github.com/egoist/mygo"
	appRepoURL         = "https://github.com/lx19990999/aria2-next-mygo"
)

func hostPlatformLabel() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}
