
ARG BUILD_DATE
ARG GIT_SHA

FROM golang:1.25-alpine3.22 AS builder

WORKDIR /usr/src/app

COPY go.mod .
RUN go mod tidy && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X main.commitSHA=${GIT_SHA} -X main.buildDate=${BUILD_DATE}"  \
    -o /usr/bin/whoami .

FROM scratch AS final

COPY --from=builder /usr/bin/whoami .

ENTRYPOINT [ "./whoami" ]
CMD [ ]
