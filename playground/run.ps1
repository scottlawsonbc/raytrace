# Build and run raytrace playground.
# The reason for the explicit build step is because Windows firewall keeps asking
# for permission to run the executable if it's built on the fly, whereas always
# creating a consistent executable path avoids this issue.
go build -o playground.exe
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed." -ForegroundColor Red;
    exit $LASTEXITCODE;
}
Write-Host "Build succeeded." -ForegroundColor Green;
./playground.exe
