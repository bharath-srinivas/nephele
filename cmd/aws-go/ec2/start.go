package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/internal/spinner"
)

// dryRun enabled.
var dryRun bool

// start instance command.
var startCmd = &cobra.Command{
	Use:     "start [instance id]",
	Short:   "Start the specified EC2 instance",
	Args:    cobra.ExactArgs(1),
	Example: "aws-go start i-0a12b345c678de",
	PreRun:  command.PreRun,
	Run:     startInstance,
}

func init() {
	command.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "perform the operation with dry run enabled")
}

// run command.
func startInstance(cmd *cobra.Command, args []string) {
	sp := spinner.Default(spinner.Prefix[2])
	sp.Start()
	sess := ec2.New(command.Session)

	instanceId := function.EC2{
		ID: args[0],
	}

	ec2Service := &function.EC2Service{
		EC2:     instanceId,
		Service: sess,
	}

	resp, err := ec2Service.StartInstance(dryRun)

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
