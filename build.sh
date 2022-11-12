docker build -t globe-test .
docker tag globe-test apiteamdevops/globe-test:0.0.3
docker push apiteamdevops/globe-test:0.0.3