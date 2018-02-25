package function

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/bharath-srinivas/aws-go/mock/ec2"
	"github.com/bharath-srinivas/aws-go/mock/lambda"
	"github.com/bharath-srinivas/aws-go/mock/rds"
	"github.com/bharath-srinivas/aws-go/mock/s3"
)

func TestNewSession(t *testing.T) {
	s := NewSession()
	got := reflect.TypeOf(s).String()
	want := "*session.Session"
	if got != want {
		t.Errorf("NewSession returned incorrect type, got: %s, want: %s", got, want)
	}
}

func TestEC2Service_GetAllInstances(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ec2ServiceMock := mock_ec2iface.NewMockEC2API(mockController)

	params := &ec2.DescribeInstancesInput{}
	ec2ServiceMock.EXPECT().DescribeInstances(params).Return(&ec2.DescribeInstancesOutput{}, nil)

	ec2Service := &EC2Service{
		Service: ec2ServiceMock,
	}

	_, err := ec2Service.GetAllInstances()

	if err != nil {
		t.Fail()
	}
}

func TestEC2Service_GetInstances(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ec2ServiceMock := mock_ec2iface.NewMockEC2API(mockController)

	var params *ec2.DescribeInstancesInput
	ec2ServiceMock.EXPECT().DescribeInstances(params).Return(&ec2.DescribeInstancesOutput{}, nil)

	ec2Service := &EC2Service{
		Service: ec2ServiceMock,
	}

	out, err := ec2Service.GetInstances()

	if reflect.TypeOf(out).String() != "[][]string" {
		t.Fail()
	} else if err != nil {
		t.Fail()
	}
}

func TestEC2Service_SetFilters(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ec2ServiceMock := mock_ec2iface.NewMockEC2API(mockController)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: aws.StringSlice([]string{
					"*web*",
					"*Web*",
					"*web*",
				}),
			},
		},
	}

	ec2ServiceMock.EXPECT().DescribeInstances(params).Return(&ec2.DescribeInstancesOutput{}, nil)

	ec2Service := &EC2Service{
		Service: ec2ServiceMock,
	}

	ec2Service.SetFilters([]string{"name=web"})

	out, err := ec2Service.GetInstances()

	if reflect.TypeOf(out).String() != "[][]string" {
		t.Fail()
	} else if err != nil {
		t.Fail()
	}
}

func TestEC2Service_LoadFiltersFromFile(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ec2ServiceMock := mock_ec2iface.NewMockEC2API(mockController)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: aws.StringSlice([]string{
					"*web*",
					"*Web*",
					"*web*",
				}),
			},
			{
				Name: aws.String("availability-zone"),
				Values: aws.StringSlice([]string{
					"*us-east-1d*",
					"*Us-East-1d*",
					"*us-east-1d*",
				}),
			},
		},
	}

	ec2ServiceMock.EXPECT().DescribeInstances(params).Return(&ec2.DescribeInstancesOutput{}, nil)

	ec2Service := &EC2Service{
		Service: ec2ServiceMock,
	}

	ec2Service.LoadFiltersFromFile("testData/filters.json")

	out, err := ec2Service.GetInstances()

	if reflect.TypeOf(out).String() != "[][]string" {
		t.Fail()
	} else if err != nil {
		t.Fail()
	}
}

func TestEC2Service_StartInstances(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ec2ServiceMock := mock_ec2iface.NewMockEC2API(mockController)

	params := &ec2.StartInstancesInput{
		DryRun:      aws.Bool(false),
		InstanceIds: aws.StringSlice([]string{"i-0a12b345c678de"}),
	}

	ec2ServiceMock.EXPECT().StartInstances(params).Return(&ec2.StartInstancesOutput{}, nil)

	ec2Input := EC2{
		IDs: []string{"i-0a12b345c678de"},
	}

	ec2Service := &EC2Service{
		EC2:     ec2Input,
		Service: ec2ServiceMock,
	}

	_, err := ec2Service.StartInstances(false)

	if err != nil {
		t.Fail()
	}
}

func TestEC2Service_StopInstances(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ec2ServiceMock := mock_ec2iface.NewMockEC2API(mockController)

	params := &ec2.StopInstancesInput{
		DryRun:      aws.Bool(false),
		InstanceIds: aws.StringSlice([]string{"i-0a12b345c678de"}),
	}

	ec2ServiceMock.EXPECT().StopInstances(params).Return(&ec2.StopInstancesOutput{}, nil)

	ec2Input := EC2{
		IDs: []string{"i-0a12b345c678de"},
	}

	ec2Service := &EC2Service{
		EC2:     ec2Input,
		Service: ec2ServiceMock,
	}

	_, err := ec2Service.StopInstances(false)

	if err != nil {
		t.Fail()
	}
}

func TestLambdaService_GetFunctions(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	lambdaServiceMock := mock_lambdaiface.NewMockLambdaAPI(mockController)

	params := &lambda.ListFunctionsInput{}
	lambdaServiceMock.EXPECT().ListFunctions(params).Return(&lambda.ListFunctionsOutput{}, nil)

	lambdaService := &LambdaService{
		Service: lambdaServiceMock,
	}

	_, err := lambdaService.GetFunctions()

	if err != nil {
		t.Fail()
	}
}

func TestLambdaService_InvokeFunction(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	lambdaServiceMock := mock_lambdaiface.NewMockLambdaAPI(mockController)

	params := &lambda.InvokeInput{
		FunctionName:   aws.String("testFunction"),
		InvocationType: aws.String("RequestResponse"),
	}
	lambdaServiceMock.EXPECT().Invoke(params).Return(&lambda.InvokeOutput{}, nil)

	functionInput := Function{
		Name: "testFunction",
	}

	lambdaService := &LambdaService{
		Function: functionInput,
		Service:  lambdaServiceMock,
	}

	_, err := lambdaService.InvokeFunction()

	if err != nil {
		t.Fail()
	}
}

func TestRDSService_GetRDSInstances(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	rdsServiceMock := mock_rdsiface.NewMockRDSAPI(mockController)

	params := &rds.DescribeDBInstancesInput{}
	rdsServiceMock.EXPECT().DescribeDBInstances(params).Return(&rds.DescribeDBInstancesOutput{}, nil)

	rdsService := &RDSService{
		Service: rdsServiceMock,
	}

	out, err := rdsService.GetRDSInstances()

	if reflect.TypeOf(out).String() != "[][]string" {
		t.Fail()
	} else if err != nil {
		t.Fail()
	}
}

func TestS3Service_GetBuckets(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	s3ServiceMock := mock_s3iface.NewMockS3API(mockController)

	params := &s3.ListBucketsInput{}
	s3ServiceMock.EXPECT().ListBuckets(params).Return(&s3.ListBucketsOutput{}, nil)

	s3Service := &S3Service{
		Service: s3ServiceMock,
	}

	_, err := s3Service.GetBuckets()

	if err != nil {
		t.Fail()
	}
}

func TestS3Service_GetObjects(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	s3ServiceMock := mock_s3iface.NewMockS3API(mockController)

	params := &s3.ListObjectsV2Input{
		Bucket:  aws.String(""),
		MaxKeys: aws.Int64(0),
	}
	s3ServiceMock.EXPECT().ListObjectsV2(params).Return(&s3.ListObjectsV2Output{}, nil)

	s3Service := &S3Service{
		Service: s3ServiceMock,
	}

	_, err := s3Service.GetObjects()

	if err != nil {
		t.Fail()
	}
}
