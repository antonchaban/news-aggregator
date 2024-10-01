package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestEKSStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewEKSStack(app, "TestEKSStack", &EKSStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
	})

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	// Validate VPC creation
	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"CidrBlock": "10.0.0.0/16",
	})

	// Validate EKS Cluster creation
	template.HasResourceProperties(jsii.String("Custom::AWSCDK-EKS-Cluster"), map[string]interface{}{
		"Config": map[string]interface{}{
			"name": "anton-eks-cluster-cdk",
		},
	})

	// Validate the creation of Node Group with expected instance type
	template.HasResourceProperties(jsii.String("AWS::EKS::Nodegroup"), map[string]interface{}{
		"InstanceTypes": []interface{}{"t2.medium"},
		"ScalingConfig": map[string]interface{}{
			"MinSize":     1,
			"MaxSize":     10,
			"DesiredSize": 2,
		},
	})

	// Validate the creation of IAM role for EKS Cluster
	template.HasResourceProperties(jsii.String("AWS::IAM::Role"), map[string]interface{}{
		"AssumeRolePolicyDocument": map[string]interface{}{
			"Statement": []interface{}{
				map[string]interface{}{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": "eks.amazonaws.com",
					},
				},
			},
			"Version": "2012-10-17",
		},
	})

	// Validate Security Group for EKS cluster
	template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), map[string]interface{}{
		"GroupDescription": "Allow traffic to EKS",
	})
}
