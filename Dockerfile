FROM golang:1.12 as builder
ENV GO111MODULE=on
WORKDIR /go/src/uac
ADD cmd ./cmd
ADD pkg ./pkg
ADD go.mod ./
ADD go.sum ./
RUN cd cmd/uac && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/uac

FROM alpine:latest
WORKDIR /uac
RUN apk --no-cache add ca-certificates
COPY --from=builder /tmp/uac ./
RUN chmod +x ./uac
USER 65535:65535
CMD  ["/uac/uac", "server"]