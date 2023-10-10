# Build Frontend
# ==============
FROM node:18.4.0-alpine3.16 AS nodebuilder
ARG VERSION=dev

WORKDIR /web
COPY web/package*.json ./
RUN npm install

COPY web/ .
RUN APP_VERSION=$VERSION npm run build

# development mode
EXPOSE 3000
# CMD npm start
CMD npm ci && npm start



# Build Backend
# =============
FROM golang:1.21.2-alpine AS gobuilder
ARG VERSION=dev

ENV PATH="/app:${PATH}"
ENV GOCACHE="/app/.cache"
WORKDIR /app

# Install debug dependencies
RUN apk add --no-cache sqlite \
  && go install github.com/cosmtrek/air@latest

COPY app/ .

# Slow build
RUN CGO_ENABLED=0 go build -v -ldflags "-s -w -X main.Version=$VERSION" -tags timetzdata -o ssham

# development mode
EXPOSE 8090
HEALTHCHECK --start-period=5s --retries=2 --interval=30s CMD ssham healthcheck
CMD ["/go/bin/air", "-c", "/app/air.toml"]



# Run App
# =======
FROM scratch

ENV PATH="/app:${PATH}"
ENV APP_ENV="production"
ENV TMPDIR="/app/pb_data/.tmp"
WORKDIR /app

COPY --from=gobuilder /app/ssham /app/ssham
COPY --from=nodebuilder /web/dist/ /app/static

EXPOSE 8090
VOLUME /app/pb_data

HEALTHCHECK --start-period=5s --retries=2 --interval=30s CMD ssham healthcheck
CMD ["ssham", "serve", "--http", "0.0.0.0:8090"]
