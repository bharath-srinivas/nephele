package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/spinner"
)

// stop instance command.
var stopCmd = &cobra.Command{
	Use:     "stop [instance id]",
	Short:   "Stop the specified EC2 instance",
	Args:    cobra.ExactArgs(1),
	Example: "aws-go stop i-0a12b345c678de",
	Run:     stopInstance,
}

func init() {
	Command.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "perform the operation with dry run enabled")
}

// run command.
func stopInstance(cmd *cobra.Command, args []string) {
	sp := spinner.Default(spinnerPrefix[2])
	sp.Start()
	sess := ec2.New(Session)

	ec2Service := &function.EC2Service{
		Service: sess,
	}

	resp, err := ec2Service.StopInstance(args[0], dryRun)

	if err != nil {
		sp.Stop()
		fmt.Println(err.Error())
	} else {
		previousState := *resp.StoppingInstances[0].PreviousState.Name
		currentState := *resp.StoppingInstances[0].CurrentState.Name
		sp.Stop()
		fmt.Println("Previous State: " + previousState + "\nCurrent State: " + currentState)
	}
}
