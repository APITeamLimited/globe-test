docker build -t globe-test .
docker tag globe-test apiteamdevops/globe-test:latest
docker push apiteamdevops/globe-test:latest