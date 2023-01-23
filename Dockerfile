FROM globe-test-base:latest

WORKDIR /

RUN rm -rf /app

COPY . /app

WORKDIR /app

RUN go build -o globe-test

ENTRYPOINT ["./globe-test"]
