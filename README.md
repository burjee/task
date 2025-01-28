# task

專案用來練習 gRPC，api 伺服器接收請求，請求再使用 gRPC 傳送給 db 伺服器對 mongodb 進行操作

## 資料夾

``` bash
├── api  # api 伺服器，接收使用者請求傳送 gRPC 請求
├── db   # db 伺服器，接收 gRPC 請求來操作 mongodb
├── grpc # gRPC 定義
└── web  # 前端頁面
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
cd grpc
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative task.proto
```

## 執行

在本地連接埠 8000 開啟網頁

```bash
docker compose up -d --build
```

## 刪除

```bash
docker compose down
```

## 測試

測試 api

```bash
cd api
go test -v task/api/routes
```

測試 db

```bash
cd db
go test -v task/db/server
```