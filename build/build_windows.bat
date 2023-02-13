@ECHO OFF

SET version=1.0.0
SET curpath=%~dp0

:: Install the go-msi package
go get github.com/mat007/go-msi
go install github.com/mat007/go-msi

:: Create service info
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go generate

MKDIR %curpath%tmp\bin\plugins
::XCOPY /s /e /y "plugins" "build\tmp\plugins"

ECHO Path: %curpath%

:: 64bit build process
SET GOARCH=amd64
go build -ldflags "-X github.com/rechecked/rcagent/internal/config.PluginDir=C:\Program Files\rcagent\plugins -X github.com/rechecked/rcagent/internal/config.ConfigDir=C:\Program Files\rcagent\" -o build/bin/rcagent.exe
go-msi make --path build/package/wix.json --msi build/rcagent-install.msi --src build/package/templates --out %curpath%build\tmp --version %version% --arch amd64

:: 32bit build process (if we need later)
::cd ..
::SET GOARCH=386
::go build -o build/tmp/rcagent.exe
::cd build
::go-msi make --path build/package/wix.json -msi build/rcagent-install-32bit.msi --version %version% --arch 386

:: Clean up the go.mod file
go mod tidy

ECHO Build completed.
PAUSE