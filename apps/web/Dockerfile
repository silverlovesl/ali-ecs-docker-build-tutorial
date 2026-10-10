# =============================================================================
#  Frontend image (build context: repository root)
#
#    docker build -f apps/web/Dockerfile -t ali-ecs-docker-build-tutorial-web:latest .
#
#  Stage 1 builds the static bundle with Bun (workspace-aware: needs the root
#  package.json / bun.lock). Stage 2 serves it with nginx; the host nginx
#  (SSL termination) proxies to this container.
# =============================================================================
FROM oven/bun:1.4-alpine AS build

WORKDIR /app

# Install workspace dependencies first (better layer caching).
# --ignore-scripts skips the root `prepare` -> `lefthook install` script, which
# needs git and makes no sense inside a build image.
COPY package.json bun.lock ./
COPY packages ./packages
COPY apps/web/package.json apps/web/
RUN bun install --frozen-lockfile --ignore-scripts

# Build the static bundle
COPY apps/web ./apps/web
RUN cd apps/web && bun run build

# -----------------------------------------------------------------------------
# Runtime: nginx serving the SPA
# -----------------------------------------------------------------------------
FROM nginx:1.29-alpine

RUN rm /etc/nginx/conf.d/default.conf
COPY apps/web/nginx.conf /etc/nginx/conf.d/app.conf
COPY --from=build /app/apps/web/dist /usr/share/nginx/html

EXPOSE 80
