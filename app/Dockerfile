FROM golang:1.27.0-alpine3.24 AS builder

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/community-app \
    ./cmd/server

FROM scratch

COPY --from=builder /out/community-app /community-app

USER 65532:65532
EXPOSE 8080

ENTRYPOINT ["/community-app"]
