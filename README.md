# go-uuid

## Features

- Generate RFC 4122 compliant UUIDs.
- Flexible state persistence options.

## Usage

Example for UUID v1:

```go
package main

import (
"fmt"
"log"
"svlada.com/uuid/v1"
)

func main() {
  uuid, err := v1.UUIDv1()
  if err != nil {
    log.Fatal("Error generating UUID:", err)
  }

  fmt.Println("Generated UUID:", uuid)
}
```

## Other Versions

Other versions of the UUID generator library, including v3, v4, v5, v6, v7 and v8 are currently in progress. Stay tuned for updates!