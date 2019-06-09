package main

import (
	"fmt"
	"go-person-crud-example/pkg/cmd/server/v2"
	"os"
)

func main() {
	if err := v1.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}