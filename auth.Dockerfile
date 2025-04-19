FROM golang:1.24.2-alpine AS build

WORKDIR /src

COPY . .

RUN go build -o ./auth ./cmd/auth/auth.go

FROM alpine

COPY --from=build --chown=nobody:nogroup /src/auth ./app/auth

USER nobody:nogroup
WORKDIR /app

ENTRYPOINT [ "./auth", "auth-server" ]
