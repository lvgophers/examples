package main

import (
	"log"
	"time"

	"github.com/lvgophers/examples/codegeneration/pseudomake/context"
)

//go:generate go run Makefile/Makefile.go

func main() {
	if context.Canceled == nil {
		log.Fatal("Please go generate github.com/j7b/codegeneration/smartexample")
	}
	ctx, cc := context.WithTimeout(context.Background(), time.Second)
	defer cc()
	select {
	case <-ctx.Done():
		log.Println("Done")
	case <-time.After(time.Minute):
		log.Println("Timeout")
	}
}
