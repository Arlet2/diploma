FROM golang:1.24.2-alpine AS build

WORKDIR /src

COPY . .

RUN go build -o ./push ./cmd/main.go

FROM alpine

COPY --from=build --chown=nobody:nogroup /src/push ./app/push

USER nobody:nogroup
WORKDIR /app

ENTRYPOINT [ "./push" ]