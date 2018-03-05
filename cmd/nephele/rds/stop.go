package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
	"github.com/bharath-srinivas/nephele/function"
	"github.com/bharath-srinivas/nephele/internal/spinner"
)

// RDS snapshot before shutdown.
var snapshot string

// rds list command.
var stopCmd = &cobra.Command{
	Use:     "stop [DB instance id]",
	Short:   "Stop the specified AWS RDS instance",
	Args:    cobra.ExactArgs(1),
	Example: "  nephele rds stop test-db-instance",
	PreRun:  command.PreRun,
	RunE:    stopInstance,
}

func init() {
	rdsCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVarP(&snapshot, "snapshot", "s", "", "take snapshot with specified name before shutting down")
}

// run command.
func stopInstance(cmd *cobra.Command, args []string) error {
	sp := spinner.Default(spinner.Prefix[2])
	sp.Start()
	sess := rds.New(command.Session)
	instanceId := function.RDS{
		ID: args[0],
	}

	if snapshot != "" {
		instanceId.SnapShotID = snapshot
	}

	rdsService := &function.RDSService{
		RDS:     instanceId,
		Service: sess,
	}

	resp, err := rdsService.StopInstance()
	if err != nil {
		sp.Stop()
		return err
	}

	sp.Stop()
	fmt.Println("Current State(" + *resp.DBInstance.DBInstanceIdentifier + ") : " +
		*resp.DBInstance.DBInstanceStatus)
	return nil
}
