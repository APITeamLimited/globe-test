docker build -t globe-test .
docker tag globe-test apiteamdevops/globe-test:globe-test:0.1.0
docker push apiteamdevops/globe-test:0.1.0