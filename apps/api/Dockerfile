# =============================================================================
#  Backend image (build context: apps/api)
#
#    docker build -t ali-ecs-docker-build-tutorial-api:latest apps/api
# =============================================================================
FROM golang:1.25-alpine AS build

WORKDIR /src

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api .

# -----------------------------------------------------------------------------
# Runtime: minimal image, non-root user
# -----------------------------------------------------------------------------
FROM alpine:3.22

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 appuser

COPY --from=build /out/api /usr/local/bin/api

USER appuser
EXPOSE 8080

ENTRYPOINT ["api"]
