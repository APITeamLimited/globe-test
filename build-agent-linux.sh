# Remove build directory if it exists
#if [ -d build-agent-linux ]; then
#    rm -rf build-agent-linux
#fi
#
## Create build directory
#mkdir build-agent-linux
#cd build-agent-linux
#
## Clone redis source code into build-agent-linux directory from github
#git clone https://github.com/redis/redis.git redis
#
## Change directory to redis source code
#cd redis
#
## Build redis
#make
#cd ..
#
## Copy redis-server binary to build-agent-linux directory
#cp redis/src/redis-server redis-server
#
## Remove redis source code
#rm -rf redis

# Build agent
GOOS=linux GOARCH=amd64 go build -o apiteam-agent

# Copy resources
mv apiteam-agent build-agent-linux/apiteam-agent

# Copy files to build directory
cp agent/targets/linux/apiteam-agent.desktop build-agent-linux/apiteam-agent.desktop
cp agent/targets/linux/snapcraft.yaml build-agent-linux/snapcraft.yaml
cp agent/targets/linux/run.sh build-agent-linux/run.sh
cp apiteam-logo.png build-agent-linux/apiteam-logo.png

# Build snap
cd build-agent-linux
snapcraft

# Clean up
cd ..
#rm build-agent-linux/apiteam-agent
#rm build-agent-linux/redis-server
#rm build-agent-linux/apiteam-agent.desktop
#rm build-agent-linux/snapcraft.yaml
#rm build-agent-linux/run.sh
#rm build-agent-linux/apiteam-logo.png
