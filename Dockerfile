FROM quoine/rocksdb:latest AS builder
ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/QuoineFinancial/vertex
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /vertex .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /vertex ./
EXPOSE 3000
ENTRYPOINT ["./vertex"]
