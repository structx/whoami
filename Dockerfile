
FROM golang:1.25-alpine3.22 AS builder

WORKDIR /usr/src/app

COPY go.mod .
RUN go mod tidy && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /usr/bin/whoami .

FROM scratch AS final

COPY --from=builder /usr/bin/whoami .

USER 1000

ENTRYPOINT [ "./whoami" ]
CMD [ ]
