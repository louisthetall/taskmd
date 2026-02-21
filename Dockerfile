# Stage 1: Build web frontend
FROM node:20-alpine AS web-builder

RUN corepack enable && corepack prepare pnpm@9 --activate

WORKDIR /build

# Copy workspace config for pnpm
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml ./
COPY apps/web/package.json apps/web/

# Install dependencies
RUN pnpm install --frozen-lockfile --filter @taskmd/web

# Copy web source and build
COPY apps/web/ apps/web/
RUN pnpm --filter @taskmd/web build

# Stage 2: Build Go binary
FROM golang:1.24-alpine AS go-builder

WORKDIR /build/apps/cli

# Download Go dependencies first (cache layer)
COPY apps/cli/go.mod apps/cli/go.sum ./
RUN go mod download

# Copy web dist from stage 1 into the embed location
COPY --from=web-builder /build/apps/web/dist /build/apps/cli/internal/web/static/dist

# Copy Go source
COPY apps/cli/ ./

# Build args for version info
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=0 go build \
    -tags embed_web \
    -ldflags="-s -w \
      -X 'github.com/driangle/taskmd/apps/cli/internal/cli.Version=${VERSION}' \
      -X 'github.com/driangle/taskmd/apps/cli/internal/cli.GitCommit=${GIT_COMMIT}' \
      -X 'github.com/driangle/taskmd/apps/cli/internal/cli.BuildDate=${BUILD_DATE}'" \
    -o /taskmd ./cmd/taskmd

# Stage 3: Minimal runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tini

# Run as non-root user
RUN adduser -D -h /home/taskmd taskmd
USER taskmd

COPY --from=go-builder /taskmd /usr/local/bin/taskmd

WORKDIR /tasks

EXPOSE 8080

ENTRYPOINT ["tini", "--"]
CMD ["taskmd", "web", "start"]
