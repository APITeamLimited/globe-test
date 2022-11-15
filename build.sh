docker build -t globe-test .
docker tag globe-test apiteamdevops/globe-test:0.0.4
docker push apiteamdevops/globe-test:0.0.4