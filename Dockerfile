FROM golang:1.16-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/kms-plugin/app/
COPY . .
RUN go get -d -v ./...
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/kms-plugin ./main.go

FROM alpine:3.12
COPY --from=builder /go/bin/kms-plugin /bin/kms-plugin
WORKDIR /bin
CMD ["/bin/kms-plugin"]