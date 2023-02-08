FROM golang:1.19-alpine3.16 as builder

WORKDIR build

COPY ./go.mod ./
COPY ./go.sum ./

RUN mkdir -p /etc/ssl/certs/ && update-ca-certificates
RUN apk update && apk add --no-cache ca-certificates && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o analytics-service ./cmd/...

FROM scratch

COPY --from=builder /go/build/analytics-service /analytics-service
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/analytics-service"]