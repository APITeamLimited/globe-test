# Remove build directory if it exists
if [ -d build-agent-darwin ]; then
    rm -rf build-agent-darwin
fi

# Create build directory
mkdir build-agent-darwin
cd build-agent-darwin

# Copy resources
mv apiteam-agent build-agent-darwin/apiteam-agent

# Copy files to build directory
cp -r agent/targets/darwin build-agent-darwin
cp apiteam-logo.png build-agent-darwin/apiteam-logo.png

# Clone redis source code into build-agent-darwin directory from github
git clone https://github.com/redis/redis.git redis

# Change directory to redis source code
cd redis

# Build redis
make
cd ..

# Copy redis-server binary to build-agent-darwin directory
cp redis/src/redis-server /build-agent-darwin/APITeam.app/Contents/MacOS/redis-server

# Remove redis source code
rm -rf redis

# Build agent
GOOS=darwin GOARCH=amd64 go build -o build-agent-darwin/APITeam.app/Contents/MacOS/apiteam-agent

# Recursively remove all gitkeep files
find . -name ".gitkeep" -type f -delete

# Build snap
cd build-agent-darwin
snapcraft

# Clean up
rm build-agent-darwin/apiteam-agent
rm build-agent-darwin/redis-server
rm build-agent-darwin/apiteam-agent.desktop
rm build-agent-darwin/snapcraft.yaml
rm build-agent-darwin/run.sh
rm build-agent-darwin/apiteam-logo.png
