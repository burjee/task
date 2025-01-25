FROM golang:1.22.6-alpine AS builder

COPY grpc /grpc

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/. .
RUN go build -o bin/task-api ./main.go

CMD ["/app/bin/task-api"]

FROM scratch

COPY --from=builder /app/bin/task-api /app/bin/task-api
COPY api/web /app/bin/web

# ENTRYPOINT
CMD ["/app/bin/task-api"]