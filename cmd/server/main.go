package main

import (
	"fmt"
	v1 "go-person-crud-example/pkg/cmd/server/v1"
	v2 "go-person-crud-example/pkg/cmd/server/v2"
	"golang.org/x/sync/errgroup"
	"os"
)

func main() {
	g := new(errgroup.Group)
	g.Go(func() error {
		if err := v1.RunServer(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		return nil
	})
	g.Go(func() error {
		if err := v2.RunServer(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}