package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	f, err := os.Create(filepath.Join("context", "context_generated.go"))
	check(err)
	defer f.Close()
	_, err = f.Write([]byte(preamble))
	check(err)
	ctxpath := filepath.Join(runtime.GOROOT(), "src", "context")
	// not checking go version in hopes it's reverted lel
	fi, err := os.Stat(ctxpath)
	if err == nil && fi.IsDir() {
		_, err = f.Write([]byte(`"context"`))
	} else {
		_, err = f.Write([]byte(`"golang.org/x/net/context"`))
	}
	check(err)
	_, err = f.Write([]byte(postamble))
	check(err)
	_, err = f.Write([]byte(initfunc))
	check(err)
	cmd := exec.Command("go", "install", ".")
	if b, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(string(b))
	}
}

const preamble = `package context

import ( 
`

const postamble = `
"time"
)
`

const initfunc = `
func init() {
	Canceled = context.Canceled
	DeadlineExceeded = context.DeadlineExceeded
	background = func() Context { return Context(context.Background()) }
	todo = func() Context { return Context(context.TODO()) }
	withcancel = func(c Context) (Context, CancelFunc) {
		ctx, cancel := context.WithCancel(c)
		return Context(ctx), CancelFunc(cancel)
	}
	withdeadline = func(c Context, t time.Time) (Context, CancelFunc) {
		ctx, cancel := context.WithDeadline(c, t)
		return Context(ctx), CancelFunc(cancel)
	}
	withtimeout = func(c Context, t time.Duration) (Context, CancelFunc) {
		ctx, cancel := context.WithTimeout(c, t)
		return Context(ctx), CancelFunc(cancel)
	}
	withvalue = func(c Context, key, value interface{}) Context {
		return Context(context.WithValue(c, key, value))
	}
}
`
