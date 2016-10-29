This is an example of using the go:generate directive with 
go run as an alternive to a Makefile.

The packages "context" in Go 1.7 (and presumably later versions)
is equivalent to "golang.org/x/net/context". There's a certain
difficulty in mixing the two, unless your project has its own
"context" package. The go:generate directive is used to wrap
either the standard context package or a vendored package as
appropriate. 

Inspired by the Go proverb "A little copying is better than a little dependency." A
lot of copying can make dependency management easier.
