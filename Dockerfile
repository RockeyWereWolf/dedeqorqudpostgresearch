FROM golang:latest

RUN go get github.com/lib/pq

WORKDIR /app

COPY . .

RUN go build -o main .

COPY kitabe-dede-qorqud.sql /app/kitabe-dede-qorqud.sql

CMD ["./main"]


