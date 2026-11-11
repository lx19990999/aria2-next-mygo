# Detached GUI launch — no console window, does not occupy the terminal.
$ErrorActionPreference = 'Stop'
$here = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $here

& "$here\build.ps1"
Get-Process aria2-next-mygo -ErrorAction SilentlyContinue | Stop-Process -Force
Start-Sleep -Milliseconds 200
Start-Process -FilePath (Join-Path $here 'aria2-next-mygo.exe') -WorkingDirectory $here -WindowStyle Normal
Write-Host 'aria2-next-mygo started (no console)'
