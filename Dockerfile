FROM golang:1.22.2-bookworm

WORKDIR /app

# Install Air
RUN go install github.com/cosmtrek/air@latest

# install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy only the necessary Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source files
COPY . .

# Build the application
# RUN GO_ENABLED=1 GOOS=linux go build -o /ssh-client cmd/api/main.go

# Expose port 7070
EXPOSE 7070

RUN chmod +x /app/startup.sh

# Command to run the application
# CMD ["/ssh-client"]
#CMD ["air", "-c", ".air.toml"]
CMD ["./startup.sh"]