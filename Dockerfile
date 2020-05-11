FROM golang:1.14-alpine AS builder

ENV PROJECT_DIR /go/src/github.com/chez-shanpu/slack-stamp
COPY ./ ${PROJECT_DIR}/
WORKDIR ${PROJECT_DIR}/
RUN go build -i -o ./.build/slamp


FROM alpine:3.11 AS prod

COPY --from=builder /go/src/github.com/chez-shanpu/slack-stamp/.build/slamp /
ENTRYPOINT ["/slamp"]