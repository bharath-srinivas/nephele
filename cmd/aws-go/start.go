package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/spinner"
)

// dryRun enabled.
var dryRun bool

// start instance command.
var startCmd = &cobra.Command{
	Use:     "start [instance id]",
	Short:   "Start the specified EC2 instance",
	Args:    cobra.ExactArgs(1),
	Example: "aws-go start i-0a12b345c678de",
	Run:     startInstance,
}

func init() {
	Command.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "perform the operation with dry run enabled")
}

// run command.
func startInstance(cmd *cobra.Command, args []string) {
	sp := spinner.Default(spinnerPrefix[2])
	sp.Start()
	sess := ec2.New(Session)

	ec2Service := &function.EC2Service{
		Service: sess,
	}

	resp, err := ec2Service.StartInstance(args[0], dryRun)

	if err != nil {
		sp.Stop()
		fmt.Println(err.Error())
	} else {
		previousState := *resp.StartingInstances[0].PreviousState.Name
		currentState := *resp.StartingInstances[0].CurrentState.Name
		sp.Stop()
		fmt.Println("Previous State: " + previousState + "\nCurrent State: " + currentState)
	}
}
