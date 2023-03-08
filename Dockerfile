# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY . /go/src/app

# Download dependencies
RUN go mod init https://github.com/RockeyWereWolf/dedeqorqudpostgresearch

# Build the application
RUN go build -o app .

# Expose the default port (8080)
EXPOSE 8080

# Set the environment variables
ENV PGHOST=db
ENV PGPORT=5432
ENV PGDATABASE=mydb
ENV PGUSER=myuser
ENV PGPASSWORD=mypassword

# Start the application
CMD ["/go/src/app/app"]

