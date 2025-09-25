
<#
.SYNOPSIS
Generate a prompt by concatenating specific Go source files.

.DESCRIPTION
This script concatenates the content of predefined `.go` files, prepends each with a comment indicating the file path,
and copies the result directly to the clipboard. It also generates an ASCII histogram showing the estimated token
contribution from each file.

.NOTES
Requires PowerShell 5.1+ on Windows.
#>

# Function to calculate estimated token count
function Estimate-Tokens {
    param (
        [string]$Text
    )
    # Approximation: Tokens are roughly 4 characters each (adjust for better accuracy if needed).
    return ([math]::Ceiling($Text.Length / 4))
}


# List of specific files
$files = @(
    "./client/app.go",
    "./client/event.go",
    "./client/main.go",
    "./client/texture.go",
    "./worker/worker.go",
    "./worker/httpfs.go",
    "./static/worker.js",
    "./static/index.html",
    "../r2/vec.go",
    "../r2/point.go",
    "../r3/vec.go",
    "../r3/point.go",
    "./event/key/key.go",
    "./event/mouse/mouse.go",
    "./event/wheel/wheel.go",
    # "./event/touch/touch.go",
    "./event/size/size.go",
    "./event/paint/paint.go",
    "./event/lifecycle/lifecycle.go",
    "../phys/scene.go",
    "./main.go",
    "./run.ps1"
)

# Initialize an empty string for the prompt
$promptText = ""
$tokenStats = @{}

foreach ($file in $files) {
    if (Test-Path $file) {
        # Get file content
        $fileContent = Get-Content -Path $file -Raw
        # Add the file path as a comment
        $promptText += "`n// $file`n"
        # Append the file content
        $promptText += $fileContent + "`n"
        # Estimate tokens for this file
        $tokenStats[$file] = Estimate-Tokens -Text $fileContent
    } else {
        Write-Host "File not found: $file" -ForegroundColor Yellow
    }
}

# Copy to clipboard
$promptText | Set-Clipboard

# Generate the ASCII histogram
$maxTokens = ($tokenStats.Values | Measure-Object -Maximum).Maximum
Write-Host "Token Contribution Histogram:" -ForegroundColor Cyan
foreach ($file in $tokenStats.Keys) {
    $tokens = $tokenStats[$file]
    $barLength = [math]::Ceiling(($tokens / $maxTokens) * 50) # Scale to a maximum width of 50 characters
    $bar = "*" * $barLength
    Write-Host ("{0,-40} [{1,5} tokens] {2}" -f $file, $tokens, $bar)
}

# Notify user
Write-Host "Prompt generated and copied to clipboard." -ForegroundColor Green
