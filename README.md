# mist-miner

Mining cloud resource usage

# Usage

Build main cli

```bash
go build -o mist-miner
```

Build plugin mm-s3

```bash
go build -o ./plugins/bin/mm-s3 ./plugins/mm-s3
```

Set which plugin binary to use

```bash
export PLUGIN_BINARY="./mm-s3"
```

Execution

```bash
# Fetch cloud resource records
./mist-miner mine

# Show given hash object file content
./mist-miner cat-file <group> <hash>
```

## gRPC build

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
    proto/miner.proto
```

## zlib uncompress

Use CLI to uncompress zlib file for viewing content

```bash
zlib-flate -uncompress < input_file_path
```

# Note

## Plugins
Plugins have the responsibility to return consistent data if there is no change on the resource. 
The shared library provided a `JsonNormalize` helper function to normalize input json string
by sorting the keys of each object. Also `MinerProperty` struct has `FormatContentValue` method
that help to set data format in normalized json string or just normal string.
