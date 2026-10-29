FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H appuser
WORKDIR /app
COPY --from=build /out/api /app/api

USER appuser
ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT [ "/app/api" ]