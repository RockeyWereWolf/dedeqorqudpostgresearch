# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY . .

# Verify that the files were copied correctly
RUN ls -al /go/src/app

# Download dependencies
RUN go mod init github.com/RockeyWereWolf/dedeqorqudpostgresearch
RUN go mod download
RUN go get github.com/lib/pq
RUN go get github.com/sirupsen/logrus
# Build the application
RUN go build -o app .

# Expose the default port (8080)
EXPOSE 8080
EXPOSE 5432

# Set the environment variables
ENV PGHOST=localhost
ENV PGPORT=5432
ENV PGDATABASE=mydb
ENV PGUSER=myuser
ENV PGPASSWORD=mypassword

# Start the application
CMD ["/go/src/app/app"]

