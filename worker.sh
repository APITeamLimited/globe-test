docker build -t worker .
docker tag worker registry.gitlab.com/apiteamcloud/worker:worker-latest
docker push registry.gitlab.com/apiteamcloud/worker:worker-latest