FROM golang:1.24.5-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/socialai ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates && addgroup -S socialai && adduser -S -G socialai socialai
WORKDIR /app
COPY --from=build /out/socialai /usr/local/bin/socialai
RUN mkdir -p /data /media && chown -R socialai:socialai /data /media
USER socialai
EXPOSE 8080
ENTRYPOINT ["socialai"]
