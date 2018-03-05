package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
	"github.com/bharath-srinivas/nephele/function"
	"github.com/bharath-srinivas/nephele/internal/spinner"
)

// dryRun enabled.
var dryRun bool

// start instance command.
var startCmd = &cobra.Command{
	Use:     "start [instance id]",
	Short:   "Start the specified EC2 instance",
	Args:    cobra.MinimumNArgs(1),
	Example: " nephele start i-0a12b345c678de",
	PreRun:  command.PreRun,
	RunE:    startInstance,
}

func init() {
	ec2Cmd.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "perform the operation with dry run enabled")
}

// run command.
func startInstance(cmd *cobra.Command, args []string) error {
	sp := spinner.Default(spinner.Prefix[2])
	sp.Start()
	sess := ec2.New(command.Session)

	instanceId := function.EC2{
		IDs: args,
	}

	ec2Service := &function.EC2Service{
		EC2:     instanceId,
		Service: sess,
	}

	resp, err := ec2Service.StartInstances(dryRun)
	if err != nil {
		sp.Stop()
		return err
	}

	sp.Stop()
	for _, data := range resp.StartingInstances {
		fmt.Println("Previous State(" + *data.InstanceId + ") : " + *data.PreviousState.Name)
		fmt.Println("Current State(" + *data.InstanceId + ")  : " + *data.CurrentState.Name)
		fmt.Printf("\n")
	}

	return nil
}
