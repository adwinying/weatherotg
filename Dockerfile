ARG GO_VERSION=1

FROM oven/bun:alpine as builder_frontend

WORKDIR /app
COPY package.json bun.lockb ./
RUN bun install
COPY . ./
RUN bun run css:build

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
COPY --from=builder_frontend /app/static/styles.css /app/static/styles.css
RUN go run github.com/a-h/templ/cmd/templ@latest generate && go build -o weatherotg .

FROM alpine:latest

COPY --from=builder /app/weatherotg /app/
CMD ["/app/weatherotg"]
