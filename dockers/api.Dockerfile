FROM quoine/rocksdb:latest AS builder
ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/QuoineFinancial/liquid-chain
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN cd cmd/api && go build -o /api .

# FROM scratch
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /api /

EXPOSE 5555
ENTRYPOINT ["/api"]
