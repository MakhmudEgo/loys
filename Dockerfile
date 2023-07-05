FROM golang:1.20.3-alpine3.17 as builder

ENV JWT_USER_SECRET_LOYS hehflkdjflkjlgajdlskfjij39u0f

WORKDIR /usr/src/app
RUN apk update && \
        apk add --no-cache git ca-certificates openssh-client &&\
        update-ca-certificates


COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./bin/loys ./cmd/loys/main.go


FROM alpine:3.17.1

COPY --from=builder /usr/src/app/bin/loys /app/loys
COPY --from=builder /usr/src/app/configs/config.yml /app/configs/config.yml
WORKDIR /app


CMD ["./loys"]
