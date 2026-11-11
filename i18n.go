package main

type langPack struct {
	AppName        string
	Downloading    string
	Waiting        string
	Stopped        string
	EngineSettings string
	ManageSettings string
	DownloadTasks  string
	WaitingQueue   string
	StoppedTitle   string
	Refresh        string
	Pause          string
	Resume         string
	Delete         string
	ClearStopped   string
	Add            string
	PasteLink      string
	NoTasks        string
	NoTasksHint    string
	Items          string
	LocalEngine    string
	RemoteEngine   string
	Disconnected   string
	UI             string
	UISettings     string
	About          string
	AboutApp       string
	AboutAppDesc   string
	AboutAppFramework string
	AppVersion     string
	HostPlatform   string
	AboutEngine       string
	AboutEngineDesc   string
	InstallEngine     string
	InstallHint       string
	InstallBrew       string
	InstallBrewNote   string
	InstallDocker     string
	InstallDockerCmd  string
	Copy              string
	Copied            string
	AfterInstall      string
	AfterInstallSteps string
	Connection        string
	Reload         string
	Save           string
	SaveN          string
	FilterParams   string
	Param          string
	Description    string
	Value          string
	RPCTitle       string
	RPCHint        string
	Address        string
	Secret         string
	Reconnect      string
	Disconnect     string
	Connected      string
	NotConnected   string
	ClearPath      string
	CatalogHint    string
	NoMatch        string
	FontSize       string
	FontFamily     string
	UIScale        string
	Language       string
	Theme          string
	ThemeLight     string
	ThemeDark      string
	ScaleStepHint  string
	LangZh         string
	LangEn         string
	UIHint         string
	AccentColor    string
	CustomColor    string
	ApplyColor     string
	PrefsPath      string
	Default        string
	Name           string
	Progress       string
	Size           string
	Speed          string
	Conns          string
	Status         string
	ManageHint     string
	EngineStatus   string
	EngineStatusStopped string
	EngineStatusStarting string
	EngineStatusRPCFailed string
	EngineStatusPortBusy string
	EngineStatusManagedConnected string
	EngineStatusExternalConnected string
	StartEngine    string
	StopEngine     string
	EngineExecutable string
	EngineExecutableHint string
	EngineConfig   string
	EngineConfigHint string
	Browse         string
	ManageOptions  string
	AutoStartEngine string
	StopEngineOnExit string
	ManageRPCHint  string
	SyncRPCFromConf string
	SyncRPCFromConfOK string
}

var i18nZH = langPack{
	AppName: "Aria2 Next", Downloading: "下载中", Waiting: "等待中", Stopped: "已停止",
	EngineSettings: "引擎设置", ManageSettings: "托管设置", DownloadTasks: "下载任务", WaitingQueue: "等待队列", StoppedTitle: "已停止",
	Refresh: "刷新", Pause: "暂停", Resume: "继续", Delete: "删除", ClearStopped: "清空", Add: "添加",
	PasteLink: "粘贴链接：HTTP / magnet / ed2k / m3u8",
	NoTasks: "暂无任务", NoTasksHint: "在上方粘贴链接后添加", Items: "%d 项",
	LocalEngine: "本机已连接", RemoteEngine: "远程已连接", Disconnected: "未连接",
	UI: "界面", UISettings: "界面设置", Connection: "连接", Reload: "重新加载", Save: "保存", SaveN: "保存 · %d",
	About: "关于", AboutApp: "本应用", AboutAppDesc: "跨平台原生 GUI，通过 JSON-RPC 连接 aria2-next。可在「托管设置」中代管引擎进程，也可只连接已运行的实例。",
	AboutAppFramework: "界面基于 mygo 原生桌面框架构建：",
	AppVersion: "版本", HostPlatform: "当前平台",
	AboutEngine: "aria2-next 引擎", AboutEngineDesc: "新一代 aria2 兼容下载引擎：保留命令行、配置、会话与 JSON-RPC，底层使用 libcurl / libtorrent / GPAC 等成熟库，支持 HTTP(S)、BitTorrent、ED2K、HLS/DASH 等协议。",
	InstallEngine: "安装 aria2-next", InstallHint: "请先从上方最新 Release 链接下载 aria2-next，再在「托管设置」中指定路径并启动，或自行运行后连接 RPC。",
	InstallBrew: "Homebrew（macOS / Linux）", InstallBrewNote: "安装后即可在终端使用 aria2-next 命令。",
	InstallDocker: "Docker", InstallDockerCmd: "docker run -d --name aria2-next -p 6800:6800 -v \"$PWD/downloads:/downloads\" -v \"$PWD/config:/config\" ghcr.io/aninsomniacy/aria2-next:latest",
	Copy: "复制", Copied: "已复制到剪贴板",
	AfterInstall: "安装后", AfterInstallSteps: "1. 在「托管设置」指定可执行文件与配置并启动，或手动启动 aria2-next\n2. 在「引擎设置」中确认地址与密钥\n3. 回到下载列表即可管理任务",
	FilterParams: "筛选全部参数…", Param: "参数", Description: "说明", Value: "值",
	RPCTitle: "RPC 连接", RPCHint: "连接本机 aria2-next；托管启停请到「托管设置」。",
	Address: "地址", Secret: "密钥", Reconnect: "重新连接", Disconnect: "断开连接", Connected: "已连接", NotConnected: "未连接",
	ClearPath: "清空",
	CatalogHint: "左侧为锚点：点击跳到对应段落。筛选作用于全部参数。",
	NoMatch: "没有匹配的参数",
	FontSize: "字体大小", FontFamily: "界面字体", UIScale: "界面缩放", Language: "界面语言", Theme: "主题",
	ThemeLight: "浅色", ThemeDark: "深色", ScaleStepHint: "步进 0.1",
	LangZh: "简体中文", LangEn: "English",
	UIHint: "仅影响本 GUI，不会写入 aria2-next。配置保存在用户目录 .config/aria2-next-mygo.conf。",
	AccentColor: "强调色", CustomColor: "自定义", ApplyColor: "应用颜色", PrefsPath: "配置文件",
	Default: "默认", Name: "名称", Progress: "进度", Size: "大小", Speed: "速度", Conns: "连接", Status: "状态",
	ManageHint: "指定本机 aria2-next 可执行文件与配置文件后，可由本 GUI 启动/停止引擎，无需手动开终端。路径默认为空，可随时清空。",
	EngineStatus: "引擎状态",
	EngineStatusStopped: "未运行",
	EngineStatusStarting: "进程已启动，正在等待 RPC 就绪…",
	EngineStatusRPCFailed: "进程在跑，但 RPC 连不上",
	EngineStatusPortBusy: "端口已有服务（可能是外部实例）",
	EngineStatusManagedConnected: "已连接（本应用托管）",
	EngineStatusExternalConnected: "已连接（外部实例）",
	StartEngine: "启动引擎", StopEngine: "停止引擎",
	EngineExecutable: "可执行文件", EngineExecutableHint: "aria2-next 或 aria2-next-*.exe 路径",
	EngineConfig: "配置文件", EngineConfigHint: "aria2next.conf 路径",
	Browse: "浏览",
	ManageOptions: "托管选项",
	AutoStartEngine: "打开本应用时自动启动引擎",
	StopEngineOnExit: "退出本应用时停止由本应用启动的引擎",
	ManageRPCHint: "启动/自动启动时会从引擎配置文件同步 rpc-listen-port 与 rpc-secret，避免连错端口。",
	SyncRPCFromConf: "从配置同步 RPC",
	SyncRPCFromConfOK: "已从引擎配置同步 RPC 地址与密钥",
}

var i18nEN = langPack{
	AppName: "Aria2 Next", Downloading: "Downloading", Waiting: "Waiting", Stopped: "Stopped",
	EngineSettings: "Engine", ManageSettings: "Managed", DownloadTasks: "Downloads", WaitingQueue: "Queue", StoppedTitle: "Stopped",
	Refresh: "Refresh", Pause: "Pause", Resume: "Resume", Delete: "Remove", ClearStopped: "Clear", Add: "Add",
	PasteLink: "Paste URI: HTTP / magnet / ed2k / m3u8",
	NoTasks: "No tasks", NoTasksHint: "Paste a link above to add one", Items: "%d items",
	LocalEngine: "Local connected", RemoteEngine: "Remote connected", Disconnected: "Disconnected",
	UI: "Interface", UISettings: "UI Settings", Connection: "Connection", Reload: "Reload", Save: "Save", SaveN: "Save · %d",
	About: "About", AboutApp: "This app", AboutAppDesc: "A cross-platform native GUI for aria2-next over JSON-RPC. Use Managed settings to start/stop the engine, or connect to an already running instance.",
	AboutAppFramework: "Built with the mygo native desktop framework:",
	AppVersion: "Version", HostPlatform: "Platform",
	AboutEngine: "aria2-next engine", AboutEngineDesc: "A next-generation aria2-compatible download engine. It keeps the CLI, config, session and JSON-RPC surface while rebuilding the core on libcurl, libtorrent, GPAC and more. Protocols include HTTP(S), BitTorrent, ED2K and HLS/DASH.",
	InstallEngine: "Install aria2-next", InstallHint: "Download aria2-next from the latest release link above, then point Managed settings at the binary or run it yourself.",
	InstallBrew: "Homebrew (macOS / Linux)", InstallBrewNote: "After install, the aria2-next command is available in your shell.",
	InstallDocker: "Docker", InstallDockerCmd: "docker run -d --name aria2-next -p 6800:6800 -v \"$PWD/downloads:/downloads\" -v \"$PWD/config:/config\" ghcr.io/aninsomniacy/aria2-next:latest",
	Copy: "Copy", Copied: "Copied to clipboard",
	AfterInstall: "After install", AfterInstallSteps: "1. Set the binary and config under Managed and start it, or launch aria2-next yourself\n2. Confirm URL and secret under Engine settings\n3. Manage downloads from the task lists",
	FilterParams: "Filter all options…", Param: "Option", Description: "Description", Value: "Value",
	RPCTitle: "RPC connection", RPCHint: "Connects to local aria2-next; use Managed settings to start/stop the process.",
	Address: "URL", Secret: "Secret", Reconnect: "Reconnect", Disconnect: "Disconnect", Connected: "Connected", NotConnected: "Disconnected",
	ClearPath: "Clear",
	CatalogHint: "Left rail jumps to sections. Filter searches every option.",
	NoMatch: "No matching options",
	FontSize: "Font size", FontFamily: "Font family", UIScale: "UI scale", Language: "Language", Theme: "Theme",
	ThemeLight: "Light", ThemeDark: "Dark", ScaleStepHint: "step 0.1",
	LangZh: "简体中文", LangEn: "English",
	UIHint: "Affects this GUI only. Saved to ~/.config/aria2-next-mygo.conf.",
	AccentColor: "Accent", CustomColor: "Custom", ApplyColor: "Apply color", PrefsPath: "Config file",
	Default: "default", Name: "Name", Progress: "Progress", Size: "Size", Speed: "Speed", Conns: "Conn", Status: "Status",
	ManageHint: "Point this GUI at a local aria2-next binary and config to start/stop the engine without a terminal. Paths default to empty and can be cleared.",
	EngineStatus: "Engine status",
	EngineStatusStopped: "Stopped",
	EngineStatusStarting: "Process started, waiting for RPC…",
	EngineStatusRPCFailed: "Process running, but RPC unreachable",
	EngineStatusPortBusy: "Port already in use (likely an external instance)",
	EngineStatusManagedConnected: "Connected (managed by this app)",
	EngineStatusExternalConnected: "Connected (external instance)",
	StartEngine: "Start engine", StopEngine: "Stop engine",
	EngineExecutable: "Executable", EngineExecutableHint: "Path to aria2-next binary",
	EngineConfig: "Config file", EngineConfigHint: "Path to aria2next.conf",
	Browse: "Browse",
	ManageOptions: "Options",
	AutoStartEngine: "Start engine automatically when this app opens",
	StopEngineOnExit: "Stop the engine started by this app when quitting",
	ManageRPCHint: "Start/auto-start syncs rpc-listen-port and rpc-secret from the engine config so the GUI does not hit the wrong port.",
	SyncRPCFromConf: "Sync RPC from config",
	SyncRPCFromConfOK: "RPC URL and secret synced from engine config",
}

func (a *app) L() langPack {
	if a.lang == "en" {
		return i18nEN
	}
	return i18nZH
}

func (a *app) fs(size float32) float32 {
	s := a.uiScale
	if s < 0.5 {
		s = 1
	}
	return size * s
}
