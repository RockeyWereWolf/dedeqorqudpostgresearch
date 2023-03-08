FROM golang:alpine

RUN apk add --no-cache git


WORKDIR /app

RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o main .

CMD ["./main"]



