FROM golang:1.22.6-alpine AS builder

COPY grpc /grpc

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

COPY db/go.mod db/go.sum .
RUN go mod download

COPY db/. .
RUN go build -o bin/task-db ./main.go

CMD ["/app/bin/task-db"]

FROM scratch

COPY --from=builder /app/bin/task-db /app/bin/task-db

# ENTRYPOINT
CMD ["/app/bin/task-db"]