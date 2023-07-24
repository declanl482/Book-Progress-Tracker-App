FROM golang:latest

WORKDIR /app

# Copy the project directory into the container's /app directory.
COPY . .

# Build the Go application inside the /app directory.
RUN go build -o Go-Book-Tracker-App

# Expose the port your Go application is listening on.
EXPOSE 8080

# Define the startup command to run your Go application.
CMD ["./go-book-tracker-app"]
