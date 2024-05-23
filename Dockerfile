ARG GO_VERSION=1

FROM oven/bun:alpine as builder_frontend

WORKDIR /app
COPY package.json bun.lockb ./
RUN bun install
COPY styles.css ./
RUN bun run css:build

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
RUN go run github.com/a-h/templ/cmd/templ@latest generate && go build -o weatherotg .

FROM alpine:latest

COPY --from=builder /app/weatherotg /app/
COPY ./static /app/
COPY --from=builder_frontend /app/static/styles.css /app/static/
CMD ["/app/weatherotg"]
