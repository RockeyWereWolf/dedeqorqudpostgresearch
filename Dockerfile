FROM golang:alpine

RUN apk add --no-cache git
RUN go mod tidy

WORKDIR /app

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]



