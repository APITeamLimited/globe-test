docker build -t worker .
docker tag worker apiteamdevops/worker:latest
docker push apiteamdevops/worker:latest