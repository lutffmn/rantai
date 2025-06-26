rantai is an http middleware chainer that provides easy way to chains middlewares.
It provides some features to manage the chainer such as Extend and Exclude.

## Installation

Before installing, [download and install Go](https://go.dev/doc/install). Go 1.24.3 or higher required.
Installation is done using `go get` command.

```bash
go get -u github.com/lutffmn/rantai
```

Or using the alternative way.
Import the module directly in source code.

```go
import "github.com/lutffmn/rantai"
```

Then install it using `go mod` tool.

```bash
go mod tidy
```

## Features

- Chainer
- Extend Middleware
- Exclude Middleware

## Examples

```go

package main

import (
  "fmt"

  "github.com/lutffmn/rantai"
)

func main(){
  // the rest of your code

  rt := rantai.New(Logger, CORS) // construct new instance of Rantai

  mux.Handle("GET /user/{id}", rt.ChainF(getUser))
  mux.Handle("POST /user", rt.Extend(Auth).ChainF(createUser))
  mux.Handle("GET /", rt.Exclude(CORS).ChainF(user))

  // the rest of your code
}
```

## Issues

If you discover an issue, please inform me by submit an issue.
