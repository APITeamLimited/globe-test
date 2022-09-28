docker build -t orchestrator .
docker tag orchestrator apiteamdevops/orchestrator:latest
docker push apiteamdevops/orchestrator:latest