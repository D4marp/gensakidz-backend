# --- Build stage ---
FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# modernc.org/sqlite is pure Go — CGO_ENABLED=0 keeps this a static binary,
# no libc/gcc needed at all in the runtime image.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server .

# --- Runtime stage ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=build /out/server ./server

# Everything that must survive a redeploy (the SQLite file + uploaded
# photos) lives under /app/data, which is meant to be mounted as a volume.
RUN mkdir -p /app/data/uploads && chown -R app:app /app
USER app

ENV DB_PATH=/app/data/gensakidz.db
ENV UPLOADS_DIR=/app/data/uploads
ENV PORT=8080

EXPOSE 8080
VOLUME ["/app/data"]

CMD ["./server"]
