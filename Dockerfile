# syntax=docker/dockerfile:1.2
FROM cgr.dev/chainguard/go as build

WORKDIR /work

# Use build args for cache keys
ARG CACHEBUST=1

# Copy only go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy the rest of the application code
COPY . ./
RUN go build -o /usr/local/bin/mini-redis ./cli/main.go


# Final image for CLI
FROM cgr.dev/chainguard/go 
COPY --from=build /usr/local/bin/mini-redis /usr/local/bin/mini-redis
ENTRYPOINT ["mini-redis"]

