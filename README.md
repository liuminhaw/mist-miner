# mist-miner

Mining cloud resource usage

# Usage

Build main cli

```bash
go build -o mist-miner
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
Plugins are responsible for returning consistent data when there are no changes to the resource. 
The `shared` library offers a `JsonNormalize` helper function that normalizes input JSON strings 
by sorting the keys of each object. Additionally, the `MinerProperty` struct includes a `FormatContentValue` method, 
which formats the data either as a normalized JSON string or as a regular string.

### Example
```go
property := shared.MinerProperty{
	Type: userDetail,
	Label: shared.MinerPropertyLabel{
		Name:   "UserDetail",
		Unique: true,
	},
	Content: shared.MinerPropertyContent{
		Format: shared.FormatJson,
	},
}
if err := property.FormatContentValue(ud.configuration.User); err != nil {
	return properties, fmt.Errorf("generate user detail: %w", err)
}
```

