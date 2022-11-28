@REM Remove build directory if it exists
if [ -d build-agent-windows ]; then
    rm -rf build-agent-windows
fi

@REM Create build directory
mkdir build-agent-windows
cd build-agent-windows


@REM Download windows port of redis

(New-Object System.Net.WebClient).DownloadFile('https://github.com/tporadowski/redis/archive/refs/tags/v5.0.14.1.zip') -OutFile 'redis.zip'

Expand-Archive redis.zip -DestinationPath redis

@REM Copy redis binaries to build directory
cp redis/redis-server.exe .

@REM Remove redis source code
rm -rf redis

# Build agent
GOOS=windows GOARCH=amd64 go build -o apiteam-agent.exe

# Copy resources
mv apiteam-agent.exe build-agent-windows/apiteam-agent.exe

# Copy files to build directory
cp agent/targets/linux/apiteam-agent.desktop build-agent-windows/apiteam-agent.desktop
cp agent/targets/linux/snapcraft.yaml build-agent-windows/snapcraft.yaml
cp agent/targets/linux/run.sh build-agent-windows/run.sh
cp apiteam-logo.png build-agent-windows/apiteam-logo.png

# Build snap
cd build-agent-windows
snapcraft

# Clean up
cd ..
#rm build-agent-windows/apiteam-agent
#rm build-agent-windows/redis-server
#rm build-agent-windows/apiteam-agent.desktop
#rm build-agent-windows/snapcraft.yaml
#rm build-agent-windows/run.sh
#rm build-agent-windows/apiteam-logo.png
