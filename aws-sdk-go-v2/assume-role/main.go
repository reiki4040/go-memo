package main

import (
	"context"
	"flag"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	optRoleArn      string
	optMFAArn       string
	optDebugSigning bool
)

func init() {
	flag.StringVar(&optRoleArn, "role-arn", "", "assume role arn. like format 'arn:aws:iam::<account id>:role/<role-name>'")
	flag.StringVar(&optMFAArn, "mfa-arn", "", "mfa arn like format 'arn:aws:iam::<account id>:mfa/<user name>'")
	flag.BoolVar(&optDebugSigning, "debug-signing", false, "show debug log about signing.")

	flag.Parse()
}

/*
https://stackoverflow.com/questions/65585709/how-to-assume-role-with-the-new-aws-go-sdk-v2-for-cross-account-access
https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/credentials/stscreds#pkg-overview
*/
func main() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-northeast-1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if optDebugSigning {
		cfg.ClientLogMode = aws.LogSigning
	}

	if optRoleArn != "" {
		stsClient := sts.NewFromConfig(cfg)
		var provider *stscreds.AssumeRoleProvider
		if optMFAArn != "" {
			provider = stscreds.NewAssumeRoleProvider(stsClient, optRoleArn, func(o *stscreds.AssumeRoleOptions) {
				o.SerialNumber = aws.String(optMFAArn)
				o.TokenProvider = stscreds.StdinTokenProvider
			})
		} else {
			provider = stscreds.NewAssumeRoleProvider(stsClient, optRoleArn)
		}
		cfg.Credentials = aws.NewCredentialsCache(provider)
	}

	cli := ec2.NewFromConfig(cfg)
	// first time, need to input MFA code if set -mfa-arn
	resp, err := cli.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range resp.Reservations {
		for _, ins := range r.Instances {
			log.Printf("%s %v", *ins.InstanceId, ins.InstanceType)
		}
	}

	// second time, not need to input MFA code because credential is cached.
	resp, err = cli.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range resp.Reservations {
		for _, ins := range r.Instances {
			log.Printf("%s %v", *ins.InstanceId, ins.InstanceType)
		}
	}
}
