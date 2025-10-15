# gateway/Dockerfile

# Development stage with Air hot-reload
FROM docker.io/library/golang:1.25-alpine AS dev

WORKDIR /build

# Install Air for hot-reload
RUN go install github.com/air-verse/air@latest

# Copy shared library and service (maintain monorepo structure)
COPY shared/ ./shared/
COPY app/backend/gateway/ ./gateway/

# Overwrite go.mod with container-specific version
COPY app/backend/gateway/go.mod.container ./gateway/go.mod

# Work in service directory
WORKDIR /build/gateway

# Download dependencies (replace directive works because ../shared exists)
RUN go mod download

# Air will watch for changes
# CMD ["air", "-c", ".air.toml"]
# CMD ["/bin/sh", "-c", "sed -i 's|../../../shared/go|../shared/go|g' go.mod && air -c .air.toml"]
CMD ["/bin/sh", "-c", "sed 's|../../../shared/go|../shared/go|g' go.mod > /tmp/go.mod && sed -i 's|go build|go build -modfile /tmp/go.mod|g' .air.toml && air -c .air.toml"]

###############################################
################# PRODUCTION ##################
###############################################
# Production build stage
FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /build

# Copy shared library and service
COPY shared/ ./shared/
COPY app/backend/gateway/ ./gateway/

# Overwrite go.mod with container-specific version
COPY app/backend/gateway/go.mod.container ./gateway/go.mod

WORKDIR /build/gateway

# Download dependencies
RUN go mod download

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Final production image
FROM scratch


COPY --from=builder /build/gateway/main /main

ENTRYPOINT ["/main"]
