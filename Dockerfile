FROM golang:1.18.3-alpine3.16 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0  \
    GOARCH="amd64" \
    GOOS=linux

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build --ldflags "-extldflags -static" -o main .

FROM alpine:3.16

WORKDIR /www

COPY --from=builder /build/main /www/
COPY --from=builder /build/database/ /www/database/
COPY --from=builder /build/public/ /www/public/
COPY --from=builder /build/storage/ /www/storage/
COPY --from=builder /build/.env /www/.env

ENTRYPOINT ["/www/main"]