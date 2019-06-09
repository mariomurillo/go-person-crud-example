package main

import (
	"fmt"
	"grpc-persona-crud-example/pkg/cmd/server"
	"os"
)

func main() {
	if err := server.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}