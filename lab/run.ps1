Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Write-Host "Starting build process..." -ForegroundColor Cyan

# Clean the dist directory
if (Test-Path ./dist) {
    Remove-Item -Recurse -Force ./dist
    Write-Host "Removed existing ./dist directory." -ForegroundColor Yellow
} else {
    Write-Host "./dist directory does not exist. Creating a new one." -ForegroundColor Yellow
}

New-Item -ItemType Directory -Path ./dist | Out-Null
Write-Host "Created new ./dist directory." -ForegroundColor Green


# Set environment variables and build the client
Write-Host "Building client.wasm..." -ForegroundColor Cyan
$env:GOOS = 'js'
$env:GOARCH = 'wasm'

# Build the client.
go build -o ./dist/client.wasm ./client
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build client.wasm." -ForegroundColor Red
    Exit 1
} else {
    Write-Host "Successfully built client.wasm." -ForegroundColor Green
}
# Build the client worker.
go build -o ./dist/worker.wasm ./worker
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build worker.wasm." -ForegroundColor Red
    Exit 1
} else {
    Write-Host "Successfully built worker.wasm." -ForegroundColor Green
}

$env:GOOS = $null
$env:GOARCH = $null

# Copy static files
Write-Host "Copying static files..." -ForegroundColor Cyan
Copy-Item ./static/favicon.ico ./dist/favicon.ico -Force
Copy-Item ./static/index.html ./dist/index.html -Force
Copy-Item ./static/worker.js ./dist/worker.js -Force
Write-Host "Copied favicon.ico and index.html and worker.js." -ForegroundColor Green

Copy-Item -Recurse ./assets ./dist/assets -Force
Write-Host "Copied assets directory." -ForegroundColor Green

# Copy wasm_exec.js from Go root
$goRoot = go env GOROOT
if (Test-Path "$goRoot/misc/wasm/wasm_exec.js") {
    Copy-Item "$goRoot/misc/wasm/wasm_exec.js" ./dist/wasm_exec.js -Force
    Write-Host "Copied wasm_exec.js from Go root." -ForegroundColor Green
} else {
    Write-Host "wasm_exec.js not found in Go root. Please ensure Go is properly installed." -ForegroundColor Red
    Exit 1
}

# Remove the server binary if it exists.
if (Test-Path ./lab.exe) {
    Remove-Item ./lab.exe -Force
    Write-Host "Removed existing server binary." -ForegroundColor Yellow
}

# Build the server
Write-Host "Building server..." -ForegroundColor Cyan
go build .
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build the server." -ForegroundColor Red
    Exit 1
} else {
    Write-Host "Successfully built the server." -ForegroundColor Green
}

# Run the server
Write-Host "Running server..." -ForegroundColor Cyan
.\lab.exe
