# AWS Go

AWS Go is a CLI tool for managing [AWS](https://aws.amazon.com) services without the need
to login to the AWS console, built to be fast and easy to use. Currently AWS Go supports services like
EC2 and Lambda only. Support for more services will be added later.

## Installation

Currently AWS Go is available only for Linux amd64 architecture. Support for other operating systems and 
architectures will be added later once the tool is stable.

On Linux run the following command to install aws-go:

`curl -sL https://raw.githubusercontent.com/bharath-srinivas/aws-go/master/setup_aws_go | sudo sh`

### AWS credentials

AWS Go requires AWS Credentials to perform operations and to manage resources. You can provide your credentials
to AWS Go and manage them by using the `env` command.

#### Managing environments

You can manage your environment profiles and create new profiles with the `env` command. It stores your config
so that you can transition between different profiles seamlessly without the need to enter your credentials
every time you switch to different environment.

For creating a new profile, use the following command:

```bash
$ aws-go env create --profile production --region us-west-1
```

You'll require the following details for creating a new profile:

* `AWS Access Key ID` your AWS account's access key
* `AWS Secret Access Key` your AWS account's secret key

You can switch between environments with the following command:

```bash
$ aws-go env use --profile staging --region eu-west-1
```

In both the above commands, the `--region` flag is `optional` and the default value will be `us-east-1` if the value
for the flag is not provided. 

#### Listing profiles

For listing all the stored profiles:

```bash
$ aws-go env --list
```

#### Deleting profile

For deleting a profile:

```bash
$ aws-go env --delete staging
``` 

### Minimum IAM policy

Below is the [AWS IAM](https://aws.amazon.com/iam) policy which provides the minimum required permissions for `aws-go`
to function.

For EC2:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "ec2:Describe*",
            "Resource": "*"
        }
    ]
}
```

For Lambda:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:GetRole",
                "iam:GetRolePolicy",
                "iam:ListAttachedRolePolicies",
                "iam:ListRolePolicies",
                "iam:ListRoles",
                "lambda:Get*",
                "lambda:List*"
            ],
            "Resource": "*"
        }
    ]
}
```

### IAM policy for starting and stopping instances

The following additional IAM policy is needed to start and stop the EC2 instances using aws-go.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "ec2:*",
            "Effect": "Allow",
            "Resource": "*"
        }
    ]
}
```

### IAM policy for invoking Lambda functions

The following additional IAM policy is required to invoke Lambda functions using aws-go.

**Note:** The following IAM policy provides permissions to a minimalistic amount of AWS resources and may vary 
according to the type of Lambda function your're invoking as your function might require access to additional
resources like EC2, CloudWatch, S3 etc. Please refer to the [official documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/introduction.html)
for more information on how to set the required policies. 

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "cloudwatch:*",
                "iam:CreateRole",
                "iam:CreatePolicy",
                "iam:AttachRolePolicy",
                "iam:PassRole",
                "lambda:*",
                "logs:*",
                "s3:*"
            ],
            "Resource": "*"
        }
    ]
}
``` 

## Getting started

**Note:** Before using aws-go, [AWS credentials](#aws-credentials) are required for using the CLI.

### Listing EC2 instances

For listing all the EC2 instances in the current selected profile, you just have to run `list`. This will list all the
available EC2 instances in a table like structure excluding the ones that are being terminated or already terminated.

#### Example

List all the available EC2 instances:

```bash
$ aws-go list
```

### Starting EC2 instances

AWS Go currently supports starting or stopping a single instance at a time. For starting an instance, you have to use
the `start` command along with the `instance-id` of the instance you want to start as the argument.

#### Example

Starting an EC2 instance:

```bash
$ aws-go start i-0a12b345c678de
```

Performing a `--dry-run` operation:

```bash
$ aws-go start --dry-run i-0a12b345c678de
```

### Stopping EC2 instances

To stop an EC2 instance, use the `stop` command along with the `instance-id` of the instance you want to start 
as the argument.

#### Example

Stopping an EC2 instance:

```bash
$ aws-go stop i-0a12b345c678de
```

Performing a `--dry-run` operation:

```bash
$ aws-go stop --dry-run i-0a12b345c678de
```

### Listing Lambda functions

AWS Go lists all the available Lambda functions and their configurations in a human friendly terminal output.

#### Example

Listing the Lambda functions and their configurations:

```bash
$ aws-go lambda list
```

### Invoking Lambda functions

AWS Go allows you to invoke the specified AWS Lambda function from the command-line and it returns the status code of 
the function call. It's important to note that `invoke` command invokes the `$LATEST` version of the lambda function
available with RequestResponse invocation type.

#### Example

Invoking a Lambda function:

```bash
$ aws-go lambda invoke testLambdaFunction
```
