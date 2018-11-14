# Log [![GoDoc](https://godoc.org/github.com/go-log/log?status.svg)](https://godoc.org/github.com/go-log/log)

Log is a logging interface for Go. That's it. Pass around the interface.

## Rationale

Users want to standardise logging. Sometimes libraries log. We leave the underlying logging implementation to the user 
while allowing libraries to log by simply expecting something that satisfies the Logger interface. This leaves 
the user free to pre-configure structure, output, etc.

## Interface

The interface is minimalistic on purpose

```go
type Logger interface {
    Log(v ...interface{})
    Logf(format string, v ...interface{})
}
```

## Example

Pre-configure a logger using [`WithFields`][logrus.WithFields] and pass it as an option to a library:

```go
import (
	"github.com/go-log/log/print"
	"github.com/lib/foo"
	"github.com/sirupsen/logrus"
)

logger := print.New(logrus.WithFields(logrus.Fields{
	"library": "github.com/lib/foo",
}))

f := foo.New(logger)
```

[logrus.WithFields]: https://godoc.org/github.com/sirupsen/logrus#WithFields
