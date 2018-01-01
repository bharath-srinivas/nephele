package main

import (
	"fmt"
	"os"

	"cmd/aws-go"
)

func main () {
	if err := cmd.Command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}