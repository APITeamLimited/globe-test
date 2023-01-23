#docker build -t globe-test-base:latest -f Dockerfile.base .
docker build -t globe-test:latest .
docker tag globe-test:latest apiteamdevops/globe-test:0.1.0
docker push apiteamdevops/globe-test:0.1.0