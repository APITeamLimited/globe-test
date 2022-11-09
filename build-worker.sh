docker build -t worker-base:latest . -f worker/DockerfileBase
docker build -t worker . -f worker/Dockerfile
docker tag worker apiteamdevops/worker:latest
docker push apiteamdevops/worker:latest