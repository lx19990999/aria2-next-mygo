# aria2-next-mygo

跨平台原生 GUI（Windows / Linux / macOS），用于管理 [aria2-next](https://github.com/AnInsomniacy/aria2-next)。

通过 JSON-RPC 连接本机或局域网中的 aria2-next，也可在「托管设置」中由本应用启动 / 停止引擎进程。

界面基于 [mygo](https://github.com/egoist/mygo) 原生桌面框架。

## 功能

- 下载任务列表：下载中 / 等待中 / 已停止
- 添加 HTTP、magnet、ed2k、m3u8 等链接
- 引擎全局选项浏览与修改（对齐 aria2-next 参数目录）
- **托管设置**：指定可执行文件与配置文件，自动启动、退出时停止
- **RPC 连接**：本机或远程地址；密钥加密保存在用户配置中
- 中英文界面与基础主题 / 缩放设置

## 下载

从 [Releases](https://github.com/lx19990999/aria2-next-mygo/releases) 获取构建产物：

| 平台 | 产物 |
| --- | --- |
| Windows | `aria2-next-mygo-*.exe`（可直接运行） |
| Linux | `.deb`、`.rpm` |
| macOS | `.dmg` |

推送以 `v` 开头的 tag（如 `v0.1.0`）会触发 GitHub Actions 自动打包并发布。

## 使用

1. 安装并准备好 [aria2-next](https://github.com/AnInsomniacy/aria2-next/releases/latest)（或仅连接已有实例）。
2. 打开本应用：
   - **托管本机引擎**：在「托管设置」填写可执行文件与 `aria2next.conf`，可勾选打开时自动启动。
   - **只连接 RPC**：在「引擎设置 → 连接」填写地址与密钥，例如：
     - 本机：`http://127.0.0.1:6800/jsonrpc`
     - 远程：`http://192.168.x.x:16800/jsonrpc`
3. 密钥需与引擎配置中的 `rpc-secret` 一致。

用户偏好保存在：

```text
~/.config/aria2-next-mygo.conf
```

（Windows 一般为 `%USERPROFILE%\.config\aria2-next-mygo.conf`）

## 从源码构建

### 依赖

- Go（版本见 `go.mod`）
- 相邻目录中的 [mygo](https://github.com/egoist/mygo) 源码（`go.mod` 的 `replace` 指向 `../mygo`）

```text
parent/
├── mygo/
└── aria2-next-mygo/   ← 本仓库
```

### 本地编译可执行文件

```bash
# Linux / macOS / Git Bash
./build.sh

# Windows PowerShell
./build.ps1
```

Windows 默认使用 `-ldflags "-H windowsgui"`（无控制台窗口）。

交叉编译示例：

```bash
GOOS=linux   GOARCH=amd64 ./build.sh
GOOS=darwin  GOARCH=arm64 OUT=aria2-next-mygo-darwin-arm64 LDFLAGS= ./build.sh
GOOS=windows GOARCH=amd64 OUT=aria2-next-mygo.exe ./build.sh
```

### 打包（deb / rpm / dmg / exe）

使用 mygo 打包（需本机或 CI 上具备对应平台条件；**dmg 需在 macOS 上生成**）：

```bash
go tool mygo build -platform windows/amd64
go tool mygo build -platform linux/amd64
go tool mygo build -platform darwin/arm64,darwin/amd64
```

产物在 `build/` 下。应用元数据见 `mygo.json`。

### 运行

```bash
# Windows
./start-gui.ps1

# Linux / macOS
./start-gui.sh
```

## 版本号

版本只来自 **git tag**（去掉前缀 `v`），不必手改 `version.go` / `mygo.json`：

| 来源 | 说明 |
| --- | --- |
| `VERSION=1.2.3` 环境变量 | 最高优先 |
| `GITHUB_REF_NAME`（如 `v1.2.3`） | CI 打 tag 发布时 |
| `git describe --tags` | 本地构建：当前精确 tag，否则最近 tag |
| 回退 | `dev` |

`./build.sh` / `./build.ps1` 会把版本注入 `-X main.appVersion=…`。  
打包前可运行 `./scripts/sync-version.sh`，把同一版本写入 `mygo.json`（仓库里的 `0.0.0` 只是占位）。

## 发布流程

直接打 tag 并推送即可：

```bash
git tag v0.2.0
git push origin v0.2.0
```

GitHub Actions（`.github/workflows/release.yml`）会从该 tag 同步版本，构建 Windows exe、Linux deb/rpm、macOS dmg，并发布到 Release。

日常 PR / push 的编译检查见 `.github/workflows/build.yml`。

## 相关项目

- [aria2-next](https://github.com/AnInsomniacy/aria2-next) — 下载引擎
- [mygo](https://github.com/egoist/mygo) — 原生 GUI 框架

## 许可证

[GPL-3.0](LICENSE)
