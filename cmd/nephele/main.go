package main

import (
	"fmt"
	"os"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"

	_ "github.com/bharath-srinivas/nephele/cmd/nephele/ec2"
	_ "github.com/bharath-srinivas/nephele/cmd/nephele/env"
	_ "github.com/bharath-srinivas/nephele/cmd/nephele/lambda"
	_ "github.com/bharath-srinivas/nephele/cmd/nephele/rds"
	_ "github.com/bharath-srinivas/nephele/cmd/nephele/s3"
	_ "github.com/bharath-srinivas/nephele/cmd/nephele/upgrade"
	_ "github.com/bharath-srinivas/nephele/cmd/nephele/version"
)

func main() {
	if err := command.Command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
