This is a stack for starting EKS cluster with CDK for news aggregator app, written in Go and deployed to AWS.
This CDK stack will create an EKS cluster with a managed node group and a VPC.
It will also install all required charts for the news aggregator app.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests
