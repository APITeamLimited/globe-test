docker build -t k6 .
docker tag orchestrator registry.gitlab.com/apiteamcloud/orchestrator:latest
docker push registry.gitlab.com/apiteamcloud/orchestrator:latest