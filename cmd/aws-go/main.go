package main

import (
	"fmt"
	"os"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"

	_ "github.com/bharath-srinivas/aws-go/cmd/aws-go/ec2"
	_ "github.com/bharath-srinivas/aws-go/cmd/aws-go/env"
	_ "github.com/bharath-srinivas/aws-go/cmd/aws-go/lambda"
	_ "github.com/bharath-srinivas/aws-go/cmd/aws-go/rds"
	_ "github.com/bharath-srinivas/aws-go/cmd/aws-go/s3"
	_ "github.com/bharath-srinivas/aws-go/cmd/aws-go/upgrade"
	_ "github.com/bharath-srinivas/aws-go/cmd/aws-go/version"
)

func main() {
	if err := command.Command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
