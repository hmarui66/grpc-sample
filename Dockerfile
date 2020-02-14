FROM golang:1.13

WORKDIR /go/src/github.com/hmarui66/grpc-sample
COPY .. .

ENTRYPOINT ["./server/start.sh"]
