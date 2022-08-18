FROM k6-base:latest

COPY . /app
WORKDIR /app
RUN go build -o main
ENTRYPOINT ["./main"]
