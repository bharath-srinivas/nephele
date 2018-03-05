package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
	"github.com/bharath-srinivas/nephele/function"
	"github.com/bharath-srinivas/nephele/internal/spinner"
)

// rds list command.
var startCmd = &cobra.Command{
	Use:     "start [DB instance id]",
	Short:   "Start the specified AWS RDS instance",
	Args:    cobra.ExactArgs(1),
	Example: "  nephele rds start test-db-instance",
	PreRun:  command.PreRun,
	RunE:    startInstance,
}

func init() {
	rdsCmd.AddCommand(startCmd)
}

// run command.
func startInstance(cmd *cobra.Command, args []string) error {
	sp := spinner.Default(spinner.Prefix[2])
	sp.Start()
	sess := rds.New(command.Session)
	instanceId := function.RDS{
		ID: args[0],
	}

	rdsService := &function.RDSService{
		RDS:     instanceId,
		Service: sess,
	}

	resp, err := rdsService.StartInstance()
	if err != nil {
		sp.Stop()
		return err
	}

	sp.Stop()
	fmt.Println("Current State(" + *resp.DBInstance.DBInstanceIdentifier + ") : " +
		*resp.DBInstance.DBInstanceStatus)
	return nil
}
