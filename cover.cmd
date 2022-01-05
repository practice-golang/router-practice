@echo off

go test ./... -coverprofile=coverage.out

curl.exe --progress-bar -Lo codecov.exe https://uploader.codecov.io/latest/windows/codecov.exe 

@REM  codecov.exe -t ${CODECOV_TOKEN}
codecov.exe -t %1

del codecov.exe /Q /S
