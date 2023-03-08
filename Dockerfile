FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

COPY kitabe-dede-qorqud.sql /app/kitabe-dede-qorqud.sql

CMD ["./main"]
