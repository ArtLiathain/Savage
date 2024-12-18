# syntax=docker/dockerfile:1

################################################################################
# Create a stage for building the application.
ARG GO_VERSION=1.23.2
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Copy go.mod and go.sum into the container
COPY go.mod go.sum ./

# Download dependencies as a separate step to take advantage of Docker's caching.
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application with GOARCH and ensure CGO is enabled for sqlite3
ARG TARGETARCH
RUN GOARCH=$TARGETARCH CGO_ENABLED=1 go build -o /bin/server ./cmd/api/main.go

################################################################################
# Create a new stage for running the application with minimal runtime dependencies.
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy the binary
COPY --from=build /bin/server /bin/

# Copy other files
COPY cmd/api/config.json /cmd/api/config.json
COPY SqlFiles /SqlFiles/

EXPOSE 8080
ENTRYPOINT ["/bin/server"]


