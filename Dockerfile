# =============================================================================
#  Single deployment image: the Go API also serves the built frontend.
#
#    docker build -t ali-ecs-docker-build-tutorial:latest .
#
#  Build context must be the repository root (Bun needs the workspace lockfile,
#  Go needs apps/api). One image, one process, one port — the host nginx
#  terminates TLS and proxies everything to 8080.
# =============================================================================

# ---- 1. frontend: Bun builds the SPA ---------------------------------------
FROM oven/bun:1.4-alpine AS web

# Public mount point of the SPA. The ECS nginx serves it under /tutorial/ (and
# strips the prefix before proxying to the container), so the bundle must
# reference /tutorial/assets/... instead of /assets/... Override with
# `--build-arg` to build an image meant to be served from the domain root.
ARG VITE_BASE_PATH=/tutorial/
ENV VITE_BASE_PATH=$VITE_BASE_PATH

WORKDIR /app

# Install workspace dependencies first (better layer caching).
# --ignore-scripts skips the root `prepare` -> `lefthook install` script, which
# needs git and makes no sense inside a build image.
COPY package.json bun.lock ./
COPY packages ./packages
COPY apps/web/package.json apps/web/
RUN bun install --frozen-lockfile --ignore-scripts

COPY apps/web ./apps/web
RUN cd apps/web && bun run build

# ---- 2. backend: compile a static Go binary --------------------------------
FROM golang:1.25-alpine AS api

WORKDIR /src

COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download

COPY apps/api ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api .

# ---- 3. runtime: API + SPA, no nginx, no supervisor ------------------------
FROM alpine:3.22

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 appuser

COPY --from=api /out/api /usr/local/bin/api
COPY --from=web /app/apps/web/dist /app/public

ENV APP_ENV=production \
    PORT=8080 \
    WEB_DIR=/app/public \
    WEB_BASE_PATH=/tutorial

USER appuser
EXPOSE 8080

ENTRYPOINT ["api"]
