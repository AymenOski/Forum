FROM golang:1.24.4

WORKDIR /Forum

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o Forum ./cmd/server/main.go

EXPOSE 8080

CMD ["./Forum"]