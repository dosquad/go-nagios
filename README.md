# Nagios helper library for applications

![CI](https://github.com/dosquad/go-nagios/workflows/CI/badge.svg)
[![GoDoc](https://godoc.org/github.com/dosquad/go-nagios?status.svg)](https://godoc.org/github.com/dosquad/go-nagios)

Go package for a writing nagios plugins

## Installation

```shell
go get -u github.com/dosquad/go-nagios
```

## Example

```golang
package main

import (
    "fmt"
    "log"

    "github.com/dosquad/go-nagios"
)

func main() {
    // Check something
    err := errors.New("Error returned by check")

    res := &nagios.Result{}

    if err != nil {
        res.Error(err)
    }

    if exitCode := nagios.Print(res); exitCode != 0 {
        os.Exit(exitCode)
    }
}
```
