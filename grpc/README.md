# gRPC 定義

## 檔案

``` bash
├── task_grpc.pb.go # 產生的 gRPC 檔案
├── task.pb.go      # 產生的 gRPC 檔案
└── task.proto      # gRPC 定義
```

## 安裝 Protocol Buffer Compiler

```bash
sudo apt install -y protobuf-compiler
```

## 安裝 Go plugins for the protocol compiler

```bash
# 安裝 go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# 更新環境變數讓 protoc 編譯器可以找到 go plugins
export PATH="$PATH:$(go env GOPATH)/bin"
```

## 產生 gRPC 程式碼

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative task.proto
```
