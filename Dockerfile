FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY cmd/api ./cmd/api
COPY internal ./internal

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /gazebo-backend ./cmd/api

# Expose port 4000
EXPOSE 4000

ENTRYPOINT ["/gazebo-backend"]