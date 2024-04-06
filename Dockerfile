FROM golang:latest

WORKDIR /app

COPY . ./src

RUN cd ./src && \
    GOPROXY="https://goproxy.io" go build -o gateway main.go && \
    cp gateway ../ && \
    cp *yaml ../ && \
    cd ../ && \
    rm -rf ./src

EXPOSE 8080

CMD ["./gateway"]