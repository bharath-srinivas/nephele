//go:generate mockgen -destination ec2/ec2iface.go github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API
//go:generate mockgen -destination lambda/lambdaiface.go github.com/aws/aws-sdk-go/service/lambda/lambdaiface LambdaAPI
//go:generate mockgen -destination rds/rdsiface.go github.com/aws/aws-sdk-go/service/rds/rdsiface RDSAPI
//go:generate mockgen -destination s3/s3iface.go github.com/aws/aws-sdk-go/service/s3/s3iface S3API

package mock
