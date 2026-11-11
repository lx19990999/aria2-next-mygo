# Cross-platform oriented Windows build helper for aria2-next-mygo.
$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

function Get-AppVersion {
	if ($env:VERSION) {
		return $env:VERSION.Trim().TrimStart('v')
	}
	if ($env:GITHUB_REF_NAME -and $env:GITHUB_REF_NAME.StartsWith('v')) {
		return $env:GITHUB_REF_NAME.TrimStart('v')
	}
	$tag = ''
	try { $tag = (git describe --tags --exact-match HEAD 2>$null) } catch { }
	if (-not $tag) {
		try { $tag = (git describe --tags --abbrev=0 2>$null) } catch { }
	}
	if ($tag) {
		return ([string]$tag).Trim().TrimStart('v')
	}
	return 'dev'
}

$out = if ($env:OUT) { $env:OUT } else { 'aria2-next-mygo.exe' }
$ver = Get-AppVersion
$ldflags = if ($env:LDFLAGS) { $env:LDFLAGS } else { '-H windowsgui' }
$ldflags = "$ldflags -X main.appVersion=$ver"

Write-Host "building $out (version=$ver windows ldflags=$ldflags)"
go build -ldflags $ldflags -o $out .
if ($LASTEXITCODE -ne 0) { throw 'go build failed' }
Write-Host "ok: $out"
