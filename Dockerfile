FROM golang:latest

RUN go install

WORKDIR /app

COPY . .

RUN go build -o main .

COPY kitabe-dede-qorqud.sql /app/kitabe-dede-qorqud.sql

CMD ["./main"]


