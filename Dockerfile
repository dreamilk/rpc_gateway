FROM golang:1.21 AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY ./ ./src

RUN cd ./src && \
    go build -o gateway main.go 

RUN cd ./src && \
    cp gateway ../ && \
    cp *yaml ../ && \
    cd ../ && \
    rm -rf ./src && \
    ls -al

FROM busybox
WORKDIR /app

COPY --from=builder /app ./

EXPOSE 8080

CMD ["./gateway"]