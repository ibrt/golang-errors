# golang-errors
[![Go Reference](https://pkg.go.dev/badge/github.com/ibrt/golang-errors.svg)](https://pkg.go.dev/github.com/ibrt/golang-errors)
![CI](https://github.com/ibrt/golang-errors/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/ibrt/golang-errors/branch/main/graph/badge.svg?token=BQVP881F9Z)](https://codecov.io/gh/ibrt/golang-errors)

Attach metadata and stack traces to Go errors.

### Basic Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"github.com/ibrt/golang-errors/errorz"
)

func main() {
    buf, err := json.MarshalIndent(errorz.ToSummary(Parent()), "", "  ")
    errorz.MaybeMustWrap(err)
    fmt.Println(string(buf))
}

func Parent() error {
    if err := Child(); err != nil {
        return errorz.Wrap(err,
            errorz.Prefix("parent"),
            errorz.Status(http.StatusInternalServerError),
            errorz.M("parent-key", "parent-value"))
    }

    return nil
}

func Child() error {
    return errorz.Errorf("child: %v",
        errorz.A("something went wrong"),
        errorz.ID("child-error"),
        errorz.M("child-key", "child-value"))
}

// Outputs:

{
    "id": "child-error",
    "status": 500,
    "metadata": {
        "child-key": "child-value",
        "parent-key": "parent-value"
    },
    "message": "parent: child: something went wrong",
    "stackTrace": [
        "package.Child (.../my/package/file.go:30)",
        ...
    ]
}

```

### Developers

Contributions are welcome, please check in on proposed implementation before sending a PR. You can validate your changes using the `./test.sh` script.
