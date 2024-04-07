FROM golang:latest AS builder

WORKDIR /app

COPY . ./src

RUN cd ./src && \
    GOPROXY="https://goproxy.io" go build -o gateway main.go && \
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