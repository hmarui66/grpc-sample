# gRPC sample

gRPC の動作確認用リポジトリ。
 
以下のソースをベースに動作確認用の処理を追加。

https://github.com/grpc/grpc-go/tree/master/examples/route_guide
 
## Start server

```
go run server/*.go
```

## Run client

`go run client/*.go [command]`

### commands

- `unary`: Unary RPC を wait 3 sec をはさみつつ直列で実行
- `unary-non-reuse-cli`: 上記を client を毎回生成しつつ実行
- `unary-non-reuse-conn`: 上記を connection を毎回生成しつつ実行
- `stream`: Stream RPC を wait 3 sec をはさみつつ直列で実行
- `stream-non-reuse-cli`: 上記を client を毎回生成しつつ実行
- `stream-non-reuse-conn`: 上記を connection を毎回生成しつつ実行
- `keep`: Unary RPC 実行後 1 min 経過後に再度 Unary RPC を実行
- `keep-without-first-call`: 1 min 経過後に Unary PRC を実行
- `many-conn`: 同時に多数のコネクションから 1 ~ 5 sec に一度 Unary RPC を実行
- `many-conn-stream`: 同時に多数のコネクションから 1 ~ 5 sec に一度 Stream RPC を実行
- `stream-image`: 7 MB の画像をアップロードする Stream API を実行
- `many-conn-stream-image`: 同時に多数のコネクションをから 1 ~ 5 sec に一度 画像アップロード Stream API を実行

## Refresh proto file

proto ファイルを修正した場合は以下を実行

`./protoc.sh`

## Run with Envoy

### Start server

```
docker build -t grpc-sample-server .
docker run --rm -it --name grpc-sample-server1 grpc-sample-server
docker run --rm -it --name grpc-sample-server2 grpc-sample-server
```

### Start Envoy

```
docker run \
    --name envoy --rm --publish 8080:80 --publish 8081:8081 \
    --link grpc-sample-server1 \
    --link grpc-sample-server2 \
    -v $PWD/envoy:/etc/envoy \
    envoyproxy/envoy:v1.9.0
``` 