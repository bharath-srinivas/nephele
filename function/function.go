// Package function contains the core functions of aws-go.
package function

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/olekukonko/tablewriter"
	"github.com/bharath-srinivas/aws-go/spinner"
	"github.com/bharath-srinivas/aws-go/store"
	"github.com/bharath-srinivas/aws-go/utils"
)

// list of spinner prefixes.
var spinnerPrefix = []string{
	"",
	"\x1b[36mfetching\x1b[m ",
	"\x1b[36mprocessing\x1b[m ",
}

// initiateSession returns an instance of AWS Session.
func initiateSession() (*session.Session) {
	accessId, secretKey, region := store.GetCredentials()

	creds := credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID: accessId,
		SecretAccessKey: secretKey})

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Credentials: creds, Region: aws.String(region)},
	}))

	return sess
}

// ListInstances renders the list of AWS EC2 instances in an ASCII table on the terminal.
func ListInstances() {
	sp := spinner.Default(spinnerPrefix[1])
	sp.Start()

	svc := ec2.New(initiateSession())

	params := &ec2.DescribeInstancesInput{}

	resp, err := svc.DescribeInstances(params)

	if err != nil {
		sp.Stop()
		fmt.Println(err.Error())
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetRowLine(true)
		table.SetHeader([]string{
			"Instance Name",
			"Instance ID",
			"Instance State",
			"Private IPv4 Address",
			"Public IPv4 Address",
			"Instance Type",
		})

		for _, i := range resp.Reservations {
			for _, t := range i.Instances {
				if *t.State.Name == "terminated" {
					continue
				}

				var instanceName, publicIP string
				if t.Tags != nil {
					for _, tag := range t.Tags {
						if *tag.Key == "Name" {
							instanceName = *tag.Value
						}
					}
				}

				if t.PublicIpAddress != nil {
					publicIP = *t.PublicIpAddress
				}

				tableData := []string{
					instanceName,
					*t.InstanceId,
					*t.State.Name,
					*t.PrivateIpAddress,
					publicIP,
					*t.InstanceType,
				}

				table.Append(tableData)
			}
		}
		sp.Stop()
		table.Render()
	}
}

// StartInstance starts the specified instance and returns the previous and current state of that instance.
func StartInstance(instanceId string, dryRun bool) {
	sp := spinner.Default(spinnerPrefix[2])
	sp.Start()

	svc := ec2.New(initiateSession())

	params := &ec2.StartInstancesInput{
		DryRun: aws.Bool(dryRun),
		InstanceIds: []*string{aws.String(instanceId)},
	}

	resp, err := svc.StartInstances(params)

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

// StopInstance stops the specified instance and returns the previous and current state of that instance.
func StopInstance(instanceId string, dryRun bool) {
	sp := spinner.Default(spinnerPrefix[2])
	sp.Start()

	svc := ec2.New(initiateSession())

	params := &ec2.StopInstancesInput{
		DryRun: aws.Bool(dryRun),
		InstanceIds: []*string{aws.String(instanceId)},
	}

	resp, err := svc.StopInstances(params)

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

// ListLambdaFunctions renders the list of available AWS Lambda functions with their configurations on the terminal.
func ListLambdaFunctions() {
	sp := spinner.Default(spinnerPrefix[1])
	sp.Start()

	svc := lambda.New(initiateSession())

	params := &lambda.ListFunctionsInput{}

	resp, err := svc.ListFunctions(params)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			sp.Stop()
			fmt.Println(aerr.Error())
		} else {
			sp.Stop()
			fmt.Println(err.Error())
		}
	} else {
		sp.Stop()
		for index, function := range resp.Functions {
			var functionDescription string
			if function.Description != nil {
				functionDescription = *function.Description
			}

			fmt.Fprintln(os.Stdout, *function.FunctionName, "\n", " -description:", functionDescription, "\n",
				" -runtime:", *function.Runtime, "\n", " -memory:", *function.MemorySize, "\n",
				" -timeout:", *function.Timeout, "\n", " -handler:", *function.Handler, "\n",
				" -role:", *function.Role, "\n", " -version:", *function.Version)
			if index < len(resp.Functions)-1 {
				fmt.Printf("\n")
			}
		}
	}
}

// InvokeLambdaFunction invokes the given function and returns the status code of the invocation.
func InvokeLambdaFunction(functionName string) {
	sp := spinner.Default(spinnerPrefix[2])
	sp.Start()

	svc := lambda.New(initiateSession())

	params := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
	}

	resp, err := svc.Invoke(params)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			sp.Stop()
			fmt.Println(aerr.Error())
		} else {
			sp.Stop()
			fmt.Println(err.Error())
		}
	} else {
		sp.Stop()
		fmt.Printf("Status Code: %d\n", *resp.StatusCode)
	}
}

// ListRDSInstances renders the list of available AWS RDS instances in an ASCII table on the terminal.
func ListRDSInstances() {
	sp := spinner.Default(spinnerPrefix[1])
	sp.Start()
	svc := rds.New(initiateSession())

	params := &rds.DescribeDBInstancesInput{}

	resp, err := svc.DescribeDBInstances(params)
	if err != nil {
		sp.Stop()
		fmt.Println(err.Error())
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetColWidth(20)
		table.SetRowLine(true)
		table.SetHeader([]string{
			"DB Instance ID",
			"DB Instance Status",
			"Endpoint",
			"DB Instance Class",
			"Engine",
			"Engine Version",
			"Multi-AZ",
		})

		for _, instance := range resp.DBInstances {
			if *instance.DBInstanceStatus == "terminated" {
				continue
			}

			dbInstanceID := utils.WordWrap(*instance.DBInstanceIdentifier, "-", 2)
			endpoint := utils.WordWrap(*instance.Endpoint.Address, ".", 2)

			tableData := []string{
				dbInstanceID,
				*instance.DBInstanceStatus,
				endpoint,
				*instance.DBInstanceClass,
				*instance.Engine,
				*instance.EngineVersion,
				strconv.FormatBool(*instance.MultiAZ),
			}

			table.Append(tableData)
		}
		sp.Stop()
		table.Render()
	}
}