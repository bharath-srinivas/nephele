package function

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/rds"
	"store"
	"utils"
)

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

func header(service string) string {
	var headers []string

	if service == "ec2" {
		headers = []string{
			"Instance Name",
			"Instance ID",
			"Instance State",
			"Private IPv4 Address",
			"Public IPv4 Address",
			"Instance Type",
		}
	} else if service == "rds" {
		headers = []string{
			"DB Instance ID",
			"DB Instance Status",
			"Endpoint",
			"DB Instance Class",
			"Engine",
			"Engine Version",
			"Multi-AZ",
		}
	}

	var underline []string

	for _, header := range headers {
		underline = append(underline, strings.Repeat("-", len(header)))
	}

	names := strings.Join(headers, "\t")
	uLine := strings.Join(underline, "\t")
	return names + "\n" + uLine + "\t"
}

func ListInstances() {
	spinner := utils.GetSpinner("fetching ")
	spinner.Start()

	svc := ec2.New(initiateSession())

	params := &ec2.DescribeInstancesInput{}

	resp, err := svc.DescribeInstances(params)

	if err != nil {
		spinner.Stop()
		fmt.Println(err.Error())
	} else {
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer writer.Flush()

		spinner.Stop()
		fmt.Fprintln(writer, header("ec2"))
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

				fmt.Fprintln(writer, instanceName + "\t", *t.InstanceId + "\t", *t.State.Name + "\t",
					*t.PrivateIpAddress + "\t", publicIP + "\t", *t.InstanceType + "\t")
			}
		}
	}
}

func StartInstance(instanceId string, dryRun bool) {
	spinner := utils.GetSpinner("processing ")
	spinner.Start()

	svc := ec2.New(initiateSession())

	params := &ec2.StartInstancesInput{
		DryRun: aws.Bool(dryRun),
		InstanceIds: []*string{aws.String(instanceId)},
	}

	resp, err := svc.StartInstances(params)

	if err != nil {
		spinner.Stop()
		fmt.Println(err.Error())
	} else {
		previousState := *resp.StartingInstances[0].PreviousState.Name
		currentState := *resp.StartingInstances[0].CurrentState.Name
		spinner.Stop()
		fmt.Println("Previous State: " + previousState + "\nCurrent State: " + currentState)
	}
}

func StopInstance(instanceId string, dryRun bool) {
	spinner := utils.GetSpinner("processing ")
	spinner.Start()

	svc := ec2.New(initiateSession())

	params := &ec2.StopInstancesInput{
		DryRun: aws.Bool(dryRun),
		InstanceIds: []*string{aws.String(instanceId)},
	}

	resp, err := svc.StopInstances(params)

	if err != nil {
		spinner.Stop()
		fmt.Println(err.Error())
	} else {
		previousState := *resp.StoppingInstances[0].PreviousState.Name
		currentState := *resp.StoppingInstances[0].CurrentState.Name
		spinner.Stop()
		fmt.Println("Previous State: " + previousState + "\nCurrent State: " + currentState)
	}
}

func ListLambdaFunctions() {
	spinner := utils.GetSpinner("fetching ")
	spinner.Start()

	svc := lambda.New(initiateSession())

	params := &lambda.ListFunctionsInput{}

	resp, err := svc.ListFunctions(params)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			spinner.Stop()
			fmt.Println(aerr.Error())
		} else {
			spinner.Stop()
			fmt.Println(err.Error())
		}
	} else {
		spinner.Stop()
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

func InvokeLambdaFunction(functionName string) {
	spinner := utils.GetSpinner("processing ")
	spinner.Start()

	svc := lambda.New(initiateSession())

	params := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
	}

	resp, err := svc.Invoke(params)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			spinner.Stop()
			fmt.Println(aerr.Error())
		} else {
			spinner.Stop()
			fmt.Println(err.Error())
		}
	} else {
		spinner.Stop()
		fmt.Printf("Status Code: %d\n", *resp.StatusCode)
	}
}

func ListRDSInstances() {
	spinner := utils.GetSpinner("fetching ")
	spinner.Start()

	svc := rds.New(initiateSession())

	params := &rds.DescribeDBInstancesInput{}

	resp, err := svc.DescribeDBInstances(params)
	if err != nil {
		spinner.Stop()
		fmt.Println(err.Error())
	} else {
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer writer.Flush()

		spinner.Stop()
		fmt.Fprintln(writer, header("rds"))
		for _, instance := range resp.DBInstances {
			if *instance.DBInstanceStatus == "terminated" {
				continue
			}

			fmt.Fprintln(writer, *instance.DBInstanceIdentifier + "\t", *instance.DBInstanceStatus + "\t",
				*instance.Endpoint.Address + "\t", *instance.DBInstanceClass + "\t", *instance.Engine + "\t",
				*instance.EngineVersion + "\t", strconv.FormatBool(*instance.MultiAZ) + "\t")
		}
	}
}