# ronn2docopt
Start with ronn man page format, and parse out synopsis for docopt usage

## Quick Start Example

1. `go run ronn2docopt.go ./examples/basic/docs/thingy.1.ronn`

This prints to stdout the [docopt usage string](http://docopt.org/)

From here, you can embed that usage string as an example by using `go-bindata`:

```
go run ronn2docopt.go ./examples/basic/docs/thingy.1.ronn > ./examples/basic/docs/docopt.txt
go-bindata --prefix "$(pwd)/examples/basic" -o ./examples/basic/lib/bindata.go ./examples/basic/docs/docopt.txt
go run ./examples/basic/thingy.go --help
```
