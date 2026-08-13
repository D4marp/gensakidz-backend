# --- Build stage ---
FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# go-sql-driver/mysql is pure Go — CGO_ENABLED=0 keeps this a static binary,
# no libc/gcc needed at all in the runtime image.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server .

# --- Runtime stage ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=build /out/server ./server

# Only uploaded photos need to survive a redeploy here — the database itself
# lives in MySQL (a separate container/volume), not on this filesystem.
RUN mkdir -p /app/data/uploads && chown -R app:app /app
USER app

ENV UPLOADS_DIR=/app/data/uploads
ENV PORT=8080

EXPOSE 8080
VOLUME ["/app/data/uploads"]

CMD ["./server"]
