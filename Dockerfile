FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o gateway main.go

EXPOSE 8080

CMD ["./gateway"]