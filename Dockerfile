FROM alpine:latest

RUN apk update
RUN apk add go

RUN go install github.com/mikerybka/host-router@latest

ENTRYPOINT ["/root/go/bin/host-router"]
