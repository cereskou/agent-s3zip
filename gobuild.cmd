@echo off
set GOOS=linux
::set GOARCH=amd64

echo clean ...
go clean

if not exist go.mod (
    echo golang mod init...
    go mod init
)

echo build...
go build -ldflags="-d -s -w" -o main
if %errorlevel% equ 0 (
    echo Zip...
    build-lambda-zip.exe --output agent-s3zip.zip main
)