# Vertex

Vertex is a replicated state machine that enables execution and storage of arbitrary functional programs in various languages targeting LLVM IR.
This repo is the official Golang implementation.

### Storage

## Development (macOS)

1. Install [Homebrew](https://brew.sh)
2. Install rocksdb

    ```bash
    brew install rocksdb
    ```

3. Compile and run

    ```bash
    go run main.go
    ```



## Docker

```
docker-compose build
docker-compose run node init
docker-compose run --service-ports node start --api
```

```
docker-compose run node unsafe_reset_all
```

