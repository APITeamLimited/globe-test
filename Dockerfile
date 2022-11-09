FROM golang:alpine

COPY . /app
WORKDIR /app
RUN go get
RUN go build -o main

ENTRYPOINT ["./main"]
