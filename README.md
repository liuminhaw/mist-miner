# mist-miner
Mining cloud resource usage

# Usage
Build main cli
```bash
go build -o mist-miner
```
Build plugin mm-s3
```bash
go build -o mm-s3 ./plugins/mm-s3
```
Set which plugin binary to use
```bash
export PLUGIN_BINARY="./mm-s3"
```
Execution
```bash
./mine
```

## gRPC build
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
    proto/miner.proto
```
