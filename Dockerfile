FROM golang:1.10-alpine AS builder
RUN apk add --no-cache git make

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go get github.com/dpetzold/kube-top/cmd/kube-top

FROM scratch
COPY --from=builder /tmp /tmp
COPY --from=builder /go/bin/kube-top /

ENTRYPOINT ["/kube-top"]
