// Package function contains the core functions of aws-go.
package function

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/bharath-srinivas/aws-go/store"
	"github.com/bharath-srinivas/aws-go/utils"
)

// Invocation type for invoking the Lambda function.
const InvocationType = "RequestResponse"

// EC2 represents the AWS EC2 instance fields.
type EC2 struct {
	Name         string // Name of the EC2 instance
	ID           string // EC2 instance ID
	State        string // Current State of the EC2 instance
	PrivateIP    string // Private IP address of the EC2 instance
	PublicIP     string // Public/Elastic IP address of the EC2 instance
	InstanceType string // EC2 instance type
}

// EC2Service represents the EC2 interface.
type EC2Service struct {
	EC2
	Service ec2iface.EC2API
}

// Function represents the Lambda function fields.
type Function struct {
	Name        string // Name of the Lambda function
	Description string // Description provided for the Lambda function, if any
	Runtime     string // Runtime of the Lambda function
	Memory      int64  // Memory allocated for the Lambda function
	Timeout     int64  // Timeout set for the Lambda function
	Handler     string // Lambda function handler
	Role        string // IAM role assigned for the Lambda function
	Version     string // Version of the Lambda function
}

// LambdaService represents the Lambda interface.
type LambdaService struct {
	Function
	Service lambdaiface.LambdaAPI
}

// RDS represents the RDS instance fields.
type RDS struct {
	InstanceID    string // RDS instance ID
	Status        string // Current status of the RDS instance
	Endpoint      string // Endpoint of the RDS instance
	InstanceClass string // RDS instance class
	Engine        string // RDS engine
	EngineVersion string // RDS engine version
	MultiAZ       bool   // Multi-AZ availability of the RDS instance
}

// RDSService represents the RDS interface.
type RDSService struct {
	RDS
	Service rdsiface.RDSAPI
}

// NewSession returns an instance of AWS Session.
func NewSession() *session.Session {
	userCreds := store.GetCredentials()

	creds := credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID:     userCreds.AccessId,
		SecretAccessKey: userCreds.SecretKey})

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Credentials: creds, Region: aws.String(userCreds.Region)},
	}))

	return sess
}

// GetInstances returns the list of specific fields of AWS EC2 instances as multidimensional slice suitable
// for rendering on a terminal ASCII table.
func (e *EC2Service) GetInstances() ([][]string, error) {
	params := &ec2.DescribeInstancesInput{}
	resp, err := e.Service.DescribeInstances(params)

	if err != nil {
		return nil, err
	}

	var result [][]string
	for _, i := range resp.Reservations {
		var ec2List []string
		for _, t := range i.Instances {
			if *t.State.Name == "terminated" {
				continue
			}

			if t.Tags != nil {
				ec2List = append(ec2List, getInstanceName(t.Tags))
			}

			var publicIP string
			if t.PublicIpAddress != nil {
				publicIP = *t.PublicIpAddress
			}

			ec2List = append(ec2List, *t.InstanceId, *t.State.Name, *t.PrivateIpAddress, publicIP, *t.InstanceType)
		}
		result = append(result, ec2List)
	}
	return result, nil
}

// StartInstances starts the specified instance and returns the state change information of that instance.
func (e *EC2Service) StartInstance(dryRun bool) (*ec2.StartInstancesOutput, error) {
	params := &ec2.StartInstancesInput{
		DryRun:      aws.Bool(dryRun),
		InstanceIds: []*string{aws.String(e.ID)},
	}
	return e.Service.StartInstances(params)
}

// StopInstances stops the specified instance and returns the state change information of that instance.
func (e *EC2Service) StopInstance(dryRun bool) (*ec2.StopInstancesOutput, error) {
	params := &ec2.StopInstancesInput{
		DryRun:      aws.Bool(dryRun),
		InstanceIds: []*string{aws.String(e.ID)},
	}
	return e.Service.StopInstances(params)
}

// GetFunctions returns the list of all the Lambda functions with their configurations.
func (l *LambdaService) GetFunctions() (*lambda.ListFunctionsOutput, error) {
	params := &lambda.ListFunctionsInput{}
	return l.Service.ListFunctions(params)
}

// InvokeFunction invokes the specified function in RequestResponse invocation type and returns the status code.
func (l *LambdaService) InvokeFunction() (*lambda.InvokeOutput, error) {
	params := &lambda.InvokeInput{
		FunctionName:   aws.String(l.Name),
		InvocationType: aws.String(InvocationType),
	}

	return l.Service.Invoke(params)
}

// GetRDSInstances returns the list of specific fields of AWS RDS instances as multidimensional slice suitable
// for rendering on a terminal ASCII table.
func (r *RDSService) GetRDSInstances() ([][]string, error) {
	params := &rds.DescribeDBInstancesInput{}
	resp, err := r.Service.DescribeDBInstances(params)

	if err != nil {
		return nil, err
	}

	var result [][]string
	for _, instance := range resp.DBInstances {
		if *instance.DBInstanceStatus == "terminated" {
			continue
		}

		var rdsList []string
		dbInstanceID := utils.WordWrap(*instance.DBInstanceIdentifier, "-", 2)
		endpoint := utils.WordWrap(*instance.Endpoint.Address, ".", 2)

		rdsList = append(rdsList, dbInstanceID, *instance.DBInstanceStatus, endpoint, *instance.DBInstanceClass,
			*instance.Engine, *instance.EngineVersion, strconv.FormatBool(*instance.MultiAZ))

		result = append(result, rdsList)
	}
	return result, nil
}

// getInstanceName is a helper function which will return the instance name from the given tag list.
func getInstanceName(tag []*ec2.Tag) string {
	var instanceName string
	for _, tag := range tag {
		if *tag.Key == "Name" {
			instanceName = *tag.Value
		}
	}
	return instanceName
}
