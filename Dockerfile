FROM golang:buster

RUN apt-get update
RUN apt-get install -y gcc libgtk-3-dev libayatana-appindicator3-dev

COPY . /app
WORKDIR /app
RUN go get
RUN go build -o globe-test

ENTRYPOINT ["./globe-test"]
