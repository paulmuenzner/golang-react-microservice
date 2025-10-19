
ARG GO_VERSION

###############################################
################# DEVELOPMENT #################
###############################################

# Development stage with Air hot-reload
FROM docker.io/library/golang:${GO_VERSION}-alpine AS dev

# Install dependencies für Health Checks
RUN apk add --no-cache curl ca-certificates

WORKDIR /build

# Install Air for hot-reload
RUN go install github.com/air-verse/air@latest

# Copy shared library and service (maintain monorepo structure)
COPY shared/ ./shared/
COPY app/backend/gateway/ ./app/backend/gateway/

# Work in service directory
WORKDIR /build/app/backend/gateway

# Download dependencies (replace directive works because ../shared exists)
RUN go mod download

# Start Air 
CMD ["air", "-c", ".air.toml"]


###############################################
################# PRODUCTION ##################
###############################################
# Production build stage
FROM docker.io/library/golang:${GO_VERSION}-alpine AS builder

WORKDIR /build

# Copy shared library and service
COPY shared/ ./shared/
COPY app/backend/gateway/ ./app/backend/gateway/

WORKDIR /build/app/backend/gateway

# Download dependencies
RUN go mod download

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Final production image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/gateway/main /main

ENTRYPOINT ["/main"]
