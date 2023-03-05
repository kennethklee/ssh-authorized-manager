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
FROM golang:1.18.3-alpine AS gobuilder
ARG VERSION=dev

ENV PATH="/app:${PATH}"
WORKDIR /app

# Install debug dependencies
RUN apk add --no-cache sqlite \
    && go install github.com/cosmtrek/air@v1.40.4

COPY app/ .

RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=${VERSION}" -o ssham

# development mode
EXPOSE 8090
HEALTHCHECK --start-period=5s --retries=2 CMD wget --no-verbose --tries=1 --spider 0:8090/api/me
CMD ["/go/bin/air", "-c", "/app/air.toml"]



# Run App
# =======
FROM scratch

ENV PATH="/app:${PATH}"
ENV APP_ENV="production"
WORKDIR /app

COPY --from=gobuilder /app/ssham /app/ssham
COPY --from=nodebuilder /web/dist/ /app/static

EXPOSE 8090
VOLUME /tmp
VOLUME /app/pb_data
CMD ["ssham", "serve", "--http", ":8090"]
