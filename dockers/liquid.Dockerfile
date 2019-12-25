FROM quoine/rocksdb:latest AS builder
ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/QuoineFinancial/liquid-chain
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN cd cmd && go build -o /liquid .


# TODO: Make the built image clean by copy binary to scratch
# FROM scratch
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /liquid ./

EXPOSE 26657
ENTRYPOINT ["/liquid"]
