# Nephelê
[![Build Status](https://travis-ci.org/bharath-srinivas/nephele.svg?branch=master)](https://travis-ci.org/bharath-srinivas/nephele)
[![GoDoc](https://godoc.org/github.com/bharath-srinivas/nephele?status.svg)](https://godoc.org/github.com/bharath-srinivas/nephele)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Nephelê (Nephele, formerly known as AWS Go) is a CLI tool for managing [AWS](https://aws.amazon.com) services without the need
to login to the AWS console, built to be fast and easy to use. Currently Nephele supports services like
EC2, Lambda, RDS etc.

## Installation

Currently Nephele is available only for linux and windows. Support for other operating systems will be added later.

On linux run the following command to install nephele:

```
curl -sL https://raw.githubusercontent.com/bharath-srinivas/nephele/master/setup_nephele | sudo -E bash -
```

If already installed, upgrade with:

```bash
$ sudo nephele upgrade
```

For windows, download the binary from [here](https://github.com/bharath-srinivas/nephele/releases).

## AWS credentials

Nephele requires AWS Credentials to perform operations and to manage resources. You can provide your credentials
to Nephele and manage them by using the `env` command.

### Managing environments

You can manage your environment profiles and create new profiles with the `env` command. It stores your config
so that you can transition between different profiles seamlessly without the need to enter your credentials
every time you switch to different environment.

For creating a new profile, use the following command:

```bash
$ nephele env create --profile production --region us-west-1
```

You'll require the following details for creating a new profile:

* `AWS Access Key ID` your AWS account's access key
* `AWS Secret Access Key` your AWS account's secret key

You can switch between environments with the following command:

```bash
$ nephele env use --profile staging --region eu-west-1
```

In both the above commands, the `--region` flag is `optional` and the default value will be `us-east-1` if the value
for the flag is not provided. 

### Listing profiles

For listing all the stored profiles:

```bash
$ nephele env --list
```

### Deleting profile

For deleting a profile:

```bash
$ nephele env --delete staging
``` 

### Minimum IAM policy

Below is the [AWS IAM](https://aws.amazon.com/iam) policy which provides the minimum required permissions for `nephele`
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

For RDS:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "rds:Describe*",
                "rds:ListTagsForResource",
                "ec2:DescribeAccountAttributes",
                "ec2:DescribeAvailabilityZones",
                "ec2:DescribeInternetGateways",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeSubnets",
                "ec2:DescribeVpcAttribute",
                "ec2:DescribeVpcs"
            ],
            "Effect": "Allow",
            "Resource": "*"
        }
    ]
}
```

#### IAM policy for starting and stopping instances

The following additional IAM policy is needed to start and stop the EC2 instances using nephele.

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

#### IAM policy for invoking Lambda functions

The following additional IAM policy is required to invoke Lambda functions using nephele.

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

**Note:** Before using nephele, [AWS credentials](#aws-credentials) are required for using the CLI.

### Listing EC2 instances

For listing all the EC2 instances in the current selected profile, you just have to run `list`. This will list all the
available EC2 instances in a table like structure excluding the ones that are being terminated or already terminated.
You can get entire info about all the instances in `JSON` format using the `-a` flag. You can also apply filters to the
list with the `--filters` or `-f` flag. The filtering functionality is supported only on normal table listing and not 
on the `-a` flag.

#### Supported filters:
The following filters are supported by the list command. It's important to note that with the `--filters` flag, you
cannot search for multiple instance names or multiple availability zones etc. In that case you can use the `JSON` file
which allows you to filter based on multiple values. 

Note that every filter is case insensitive:

* `name` instance name
* `id` instance ID
* `state` instance state
* `type` instance type
* `az` availability zone of the instance

#### Precedence:
The precedence for loading filters is:

* filters from the flag
* filters from JSON file

#### Example filters file:
     
```json
 [
   {
     "name": "name",
     "values": ["web"]
   },
   {
     "name": "az",
     "values": ["us-east-1a", "us-east-1d"]
   }
 ]
```

#### Example

List all the available EC2 instances in a table format:

```bash
$ nephele ec2 list
```

Apply filters to the list:

```bash
$ nephele ec2 list --filters name=web,az=us-east-1a
```

Applying filters from a `JSON` file:

```bash
$ nephele ec2 list -F filters.json
```

Listing everything in `JSON` format:

```bash
$ nephele ec2 list --all
```

Piping `JSON` to a `JSON` file:

```bash
$ nephele ec2 list --all > ec2.json
```

Performing `less` on `JSON` output:

```bash
$ nephele ec2 list --all | less
```

### Starting EC2 instances

Nephele supports starting or stopping multiple instances. For starting an instance, you have to use the `start` command
along with the `instance-id` of the instance you want to start as the argument.

#### Example

Starting an EC2 instance:

```bash
$ nephele ec2 start i-0a12b345c678de
```

Starting multiple EC2 instances:

```bash
$ nephele ec2 start i-0a12b345c678de i-0b12c345d678ef
```

Performing a `--dry-run` operation:

```bash
$ nephele ec2 start --dry-run i-0a12b345c678de
```

### Stopping EC2 instances

To stop an EC2 instance, use the `stop` command along with the `instance-id` of the instance you want to stop 
as the argument.

#### Example

Stopping an EC2 instance:

```bash
$ nephele ec2 stop i-0a12b345c678de
```

Stopping multiple EC2 instances:

```bash
$ nephele ec2 stop i-0a12b345c678de i-0b12c345d678ef
```

Performing a `--dry-run` operation:

```bash
$ nephele ec2 stop --dry-run i-0a12b345c678de
```

### Listing Lambda functions

Nephele lists all the available Lambda functions and their configurations in a human friendly terminal output.

#### Example

Listing the Lambda functions and their configurations:

```bash
$ nephele lambda list
```

### Invoking Lambda functions

Nephele allows you to invoke the specified AWS Lambda function from the command-line and it returns the status code of 
the function call. It's important to note that `invoke` command invokes the `$LATEST` version of the lambda function
available with RequestResponse invocation type.

#### Example

Invoking a Lambda function:

```bash
$ nephele lambda invoke testLambdaFunction
```

### Listing RDS instances

Nephele lists only the available RDS instances excluding the ones that are being terminated or already terminated. 
Nephele provides only the basic information about RDS instances since the terminal cannot accommodate all the information
about RDS instances. This might be improved in the future.

#### Example

Listing the RDS instances in a table:

```bash
$ nephele rds list
```

### Starting a RDS instance

To start a RDS instance, use the `start` command along with the `db-instance-id` of the RDS instance you want to start
as the argument.

#### Example

Starting a RDS instance:

```bash
$ nephele rds start test-rds-instance
```

### Stopping a RDS instance

To stop a RDS instance, use the `stop` command along with the `db-instance-id` of the RDS instance you want to stop as
the argument.

#### Example

Stopping a RDS instance:

```bash
$ nephele rds stop test-rds-instance
```

Taking snapshot of RDS instance before stopping:

```bash
$ nephele rds stop test-rds-instance --snapshot test-rds-instance-snapshot
```

### Listing S3 buckets

Nephele lists all buckets with their name and creation date in an ascii table.

#### Example

Listing the S3 buckets in a table:

```bash
$ nephele s3 list
```

### Listing S3 objects

To list all the S3 objects present in a bucket, just specify the `bucket name` along with the `list` command. This will
render the list of S3 objects in a pager so that you can perform operations like search, scroll through the list etc.

#### Example

To list the S3 objects in a bucket:

```bash
$ nephele s3 list test-bucket
```

To list the next set or previous set of objects in a bucket:

```bash
$ nephele s3 list test-bucket -t [token]
```

You'll get this token string from the pager.

To list more objects than the default limit of 100:

```bash
$ nephele s3 list test-bucket -c 500
```

Note that the maximum number of objects you can fetch per request is limited to `1000`.

### Downloading S3 objects

You can download either a single S3 object from a bucket or download multiple S3 objects from multiple buckets
concurrently in batch. To download S3 objects in batch, you have to provide the list of objects you want to download in
`JSON` file format to the `-o` flag.

The file should contain the following keys and their respective values:
* `bucket_name` the name of the S3 bucket
* `object_name` the S3 object name
* `file_name` the path to file name. The S3 object will be downloaded at this path with the provided file name

#### Example objects file:

```json
[
  {
    "bucket_name": "bucket-1",
    "object_name": "object-1.png",
    "file_name": "images/hello.png"
  },
  {
    "bucket_name": "bucket-1/docs/",
    "object_name": "hello.doc",
    "file_name": "hello.doc"
  }
]
```

#### Example

To download a S3 object:

```bash
$ nephele s3 download test-bucket:test.png test.png
```

To download an object from sub directory of a bucket:

```bash
$ nephele s3 download test-bucket/images/:hello.png hello.png
```

Note: Sub-directory name is case-sensitive and requires '/' at the end

To download multiple objects concurrently:

```bash
$ nephele s3 download -o objects-file.json
```