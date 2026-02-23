# Hanzo S3 Admin Go SDK

The Hanzo S3 Admin Go SDK provides APIs to manage [Hanzo S3](https://github.com/hanzoai/s3) services.

This document assumes that you have a working [Go setup](https://golang.org/doc/install).

## Initialize Hanzo S3 Admin Client object.

```go
package main

import (
    "fmt"

    // Note: The Go module path remains github.com/minio/madmin-go for upstream compatibility.
    "github.com/minio/madmin-go/v4"
)

func main() {
    // Use a secure connection.
    ssl := true

    // Initialize Hanzo S3 admin client object.
    mdmClnt, err := madmin.New("your-s3.example.com:9000", "YOUR-ACCESSKEYID", "YOUR-SECRETKEY", ssl)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Fetch service status.
    info, err := mdmClnt.ClusterInfo(context.Background())
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("%#v\n", info)
}
```

## Documentation

All documentation is available [here](https://pkg.go.dev/github.com/minio/madmin-go/v4).

> Note: The Go module path remains `github.com/minio/madmin-go` for upstream compatibility.

## License

This SDK is licensed under [GNU AGPLv3](https://github.com/hanzos3/admin-go/blob/main/LICENSE).
