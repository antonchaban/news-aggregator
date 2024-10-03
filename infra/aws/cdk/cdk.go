package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	vpcCIDR       = "10.0.0.0/16"
	myClusterName = "anton-eks-cluster-cdk"
	nodeInstance  = "t2.medium"
	minNodeSize   = 1
	maxNodeSize   = 10
	desiredSize   = 2
	sshKey        = "anton"
)

type EKSStackProps struct {
	awscdk.StackProps
}

// Params is a struct to hold parameters for the EKS stack configuration from context
type Params struct {
	VpcCidr          string
	ClusterName      string
	NodeInstanceType string
	MinNodeSize      float64
	MaxNodeSize      float64
	DesiredNodeSize  float64
	SshKey           string
}

type EnvParams struct {
	Account string
	Region  string
}

func NewEKSStack(scope constructs.Construct, id string, props *EKSStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Fetch parameters as a struct
	params := fetchParams(stack)

	// Create VPC with public and private subnets
	vpc := awsec2.NewVpc(stack, jsii.String("antonvpc-cdk"), &awsec2.VpcProps{
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String(params.VpcCidr)),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:                jsii.String("anton-public-subnet"),
				SubnetType:          awsec2.SubnetType_PUBLIC,
				CidrMask:            jsii.Number(20),
				MapPublicIpOnLaunch: jsii.Bool(true),
			},
			{
				Name:       jsii.String("anton-private-subnet"),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
				CidrMask:   jsii.Number(20),
			},
		},
		MaxAzs: jsii.Number(3),
	})

	igw := vpc.InternetGatewayId()

	// Create Route Table for public subnets
	routeTable := awsec2.NewCfnRouteTable(stack, jsii.String("antonroutetable-cdk"), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
	})

	// Add a route to the internet gateway
	awsec2.NewCfnRoute(stack, jsii.String("antonroute-cdk"), &awsec2.CfnRouteProps{
		RouteTableId:         routeTable.Ref(),
		DestinationCidrBlock: jsii.String("0.0.0.0/0"),
		GatewayId:            igw,
	})

	// Associate subnets with route table
	publicSubnets := vpc.PublicSubnets()
	for i, subnet := range *publicSubnets {
		awsec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String(fmt.Sprintf("SubnetAssociation%d", i)), &awsec2.CfnSubnetRouteTableAssociationProps{
			SubnetId:     subnet.SubnetId(),
			RouteTableId: routeTable.Ref(),
		})
	}

	// IAM Role for EKS Cluster
	clusterRole := awsiam.NewRole(stack, jsii.String("antonclusterrole-cdk"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSServicePolicy")),
		},
	})

	// Security Group for EKS cluster
	clusterSG := awsec2.NewSecurityGroup(stack, jsii.String("antonclusterSG-cdk"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		SecurityGroupName: jsii.String("EKSClusterSG-cdk"),
		Description:       jsii.String("Allow traffic to EKS"),
		AllowAllOutbound:  jsii.Bool(true),
	})

	clusterSG.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(443)), jsii.String("Allow HTTPS"), jsii.Bool(true))
	clusterSG.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow HTTP"), jsii.Bool(true))
	clusterSG.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(22)), jsii.String("Allow SSH"), jsii.Bool(true))

	// Create EKS Cluster
	cluster := awseks.NewCluster(stack, jsii.String("antonekscluster-cdk"), &awseks.ClusterProps{
		EndpointAccess:  awseks.EndpointAccess_PUBLIC(),
		ClusterName:     jsii.String(params.ClusterName),
		Vpc:             vpc,
		DefaultCapacity: jsii.Number(0),
		SecurityGroup:   clusterSG,
		Version:         awseks.KubernetesVersion_V1_30(),
		Role:            clusterRole,
		VpcSubnets: &[]*awsec2.SubnetSelection{
			{SubnetType: awsec2.SubnetType_PUBLIC},
		},
	})

	iamUserArn := "arn:aws:iam::406477933661:user/anton"
	cluster.AwsAuth().AddUserMapping(awsiam.User_FromUserArn(stack, jsii.String("anton"), jsii.String(iamUserArn)), &awseks.AwsAuthMapping{
		Username: jsii.String("anton"),
		Groups: &[]*string{
			jsii.String("system:masters"),
		},
	})

	// IAM Role for Node Group
	nodegRole := awsiam.NewRole(stack, jsii.String("antonnodegrouprole-cdk"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSWorkerNodePolicy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryReadOnly")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2FullAccess")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKS_CNI_Policy")),
		},
	})

	cluster.AddNodegroupCapacity(jsii.String("anton-node-group-cdk"), &awseks.NodegroupOptions{
		InstanceTypes: &[]awsec2.InstanceType{
			awsec2.NewInstanceType(jsii.String(params.NodeInstanceType)),
		},
		NodeRole:    nodegRole,
		MinSize:     jsii.Number(params.MinNodeSize),
		MaxSize:     jsii.Number(params.MaxNodeSize),
		DesiredSize: jsii.Number(params.DesiredNodeSize),
		RemoteAccess: &awseks.NodegroupRemoteAccess{
			SshKeyName: jsii.String(params.SshKey),
		},
		Subnets: &awsec2.SubnetSelection{
			Subnets: publicSubnets,
		},
		AmiType:  awseks.NodegroupAmiType_AL2_X86_64,
		DiskSize: jsii.Number(20),
	})

	cluster.AddHelmChart(jsii.String("argo-cd"), &awseks.HelmChartOptions{
		Chart:      jsii.String("argo-cd"),
		Repository: jsii.String("https://argoproj.github.io/argo-helm"),
		Release:    jsii.String("argo-cd"),
		Namespace:  jsii.String("argocd"),
		Values: &map[string]interface{}{
			"installCRDs": true,
			"server": map[string]interface{}{
				"service": map[string]interface{}{
					"type": "LoadBalancer",
				},
			},
		},
	})

	// EKS Add-ons
	awseks.NewCfnAddon(stack, jsii.String("VPCCNIAddon"), &awseks.CfnAddonProps{
		ClusterName:      cluster.ClusterName(),
		AddonName:        jsii.String("vpc-cni"),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("CoreDNSAddon"), &awseks.CfnAddonProps{
		ClusterName:      cluster.ClusterName(),
		AddonName:        jsii.String("coredns"),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("KubeProxyAddon"), &awseks.CfnAddonProps{
		ClusterName:      cluster.ClusterName(),
		AddonName:        jsii.String("kube-proxy"),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("PodIdentityAddon"), &awseks.CfnAddonProps{
		ClusterName:      cluster.ClusterName(),
		AddonName:        jsii.String("eks-pod-identity-agent"),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	// Outputs
	awscdk.NewCfnOutput(stack, jsii.String("EKSClusterName"), &awscdk.CfnOutputProps{
		Value: cluster.ClusterName(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("VPCId"), &awscdk.CfnOutputProps{
		Value: vpc.VpcId(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("ClusterSecurityGroupId"), &awscdk.CfnOutputProps{
		Value: clusterSG.SecurityGroupId(),
	})

	return stack
}

func getString(stack awscdk.Stack, key string, defaultValue string) string {
	if value := stack.Node().TryGetContext(jsii.String(key)); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getFloat64(stack awscdk.Stack, key string, defaultValue float64) float64 {
	if value := stack.Node().TryGetContext(jsii.String(key)); value != nil {
		if num, ok := value.(float64); ok {
			return num
		}
	}
	return defaultValue
}

func fetchParams(stack awscdk.Stack) Params {
	return Params{
		VpcCidr:          getString(stack, "vpcCidr", vpcCIDR),
		ClusterName:      getString(stack, "eksClusterName", myClusterName),
		NodeInstanceType: getString(stack, "nodeInstanceType", nodeInstance),
		MinNodeSize:      getFloat64(stack, "minNodeSize", minNodeSize),
		MaxNodeSize:      getFloat64(stack, "maxNodeSize", maxNodeSize),
		DesiredNodeSize:  getFloat64(stack, "desiredNodeSize", desiredSize),
		SshKey:           getString(stack, "ec2SshKey", sshKey),
	}
}

func fetchEnvParams(app constructs.Construct) EnvParams {
	getString := func(key string, defaultValue string) string {
		if value := app.Node().TryGetContext(jsii.String(key)); value != nil {
			if str, ok := value.(string); ok {
				return str
			}
		}
		return defaultValue
	}

	return EnvParams{
		Account: getString("account", "406477933661"),
		Region:  getString("region", "us-west-2"),
	}
}

func main() {
	app := awscdk.NewApp(nil)

	envParams := fetchEnvParams(app)

	stackProps := awscdk.StackProps{
		Env: env(envParams),
	}

	NewEKSStack(app, "MyEKSStack", &EKSStackProps{
		StackProps: stackProps,
	})

	app.Synth(nil)
}

func env(params EnvParams) *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(params.Account),
		Region:  jsii.String(params.Region),
	}
}
