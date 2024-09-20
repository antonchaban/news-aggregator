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

type EKSStackProps struct {
	awscdk.StackProps
}

func NewEKSStack(scope constructs.Construct, id string, props *EKSStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Create VPC with public and private subnets
	vpc := awsec2.NewVpc(stack, jsii.String("antonvpc-cdk"), &awsec2.VpcProps{
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String("10.0.0.0/16")),
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
		ClusterName:     jsii.String("anton-eks-cluster-cdk"),
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
			awsec2.NewInstanceType(jsii.String("t2.medium")),
		},
		NodeRole:    nodegRole,
		MinSize:     jsii.Number(1),
		MaxSize:     jsii.Number(10),
		DesiredSize: jsii.Number(2),
		RemoteAccess: &awseks.NodegroupRemoteAccess{
			SshKeyName: jsii.String("anton"),
		},
		Subnets: &awsec2.SubnetSelection{
			Subnets: publicSubnets,
		},
		AmiType:  awseks.NodegroupAmiType_AL2_X86_64,
		DiskSize: jsii.Number(20),
	})

	awseks.NewKubernetesManifest(stack, jsii.String("CertManagerNamespace"), &awseks.KubernetesManifestProps{
		Cluster: cluster,
		Manifest: &[]*map[string]interface{}{
			{
				"apiVersion": "v1",
				"kind":       "Namespace",
				"metadata": map[string]interface{}{
					"name": "cert-manager",
				},
			},
		},
	})

	awseks.NewKubernetesManifest(stack, jsii.String("NewsAlligatorNamespace"), &awseks.KubernetesManifestProps{
		Cluster: cluster,
		Manifest: &[]*map[string]interface{}{
			{
				"apiVersion": "v1",
				"kind":       "Namespace",
				"metadata": map[string]interface{}{
					"name": "news-alligator",
				},
			},
		},
	})

	awseks.NewKubernetesManifest(stack, jsii.String("OperatorNamespace"), &awseks.KubernetesManifestProps{
		Cluster: cluster,
		Manifest: &[]*map[string]interface{}{
			{
				"apiVersion": "v1",
				"kind":       "Namespace",
				"metadata": map[string]interface{}{
					"name": "operator-system",
				},
			},
		},
	})

	awseks.NewHelmChart(stack, jsii.String("CertManagerHelmChart"), &awseks.HelmChartProps{
		Cluster:    cluster,
		Chart:      jsii.String("cert-manager"),
		Repository: jsii.String("https://charts.jetstack.io"),
		Release:    jsii.String("cert-manager"),
		Version:    jsii.String("v1.15.3"),
		Namespace:  jsii.String("cert-manager"),
		Values: &map[string]interface{}{
			"crds": map[string]interface{}{
				"enabled": true,
			},
		},
	})

	awseks.NewHelmChart(stack, jsii.String("VPAHelmChart"), &awseks.HelmChartProps{
		Cluster:    cluster,
		Chart:      jsii.String("vertical-pod-autoscaler"),
		Repository: jsii.String("https://cowboysysop.github.io/charts/"),
		Release:    jsii.String("my-release"),
		Namespace:  jsii.String("default"),
	})

	awseks.NewHelmChart(stack, jsii.String("EbsCsiDriverHelmChart"), &awseks.HelmChartProps{
		Cluster:    cluster,
		Chart:      jsii.String("aws-ebs-csi-driver"),
		Repository: jsii.String("https://kubernetes-sigs.github.io/aws-ebs-csi-driver"),
		Release:    jsii.String("aws-ebs-csi"),
		Namespace:  jsii.String("kube-system"),
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

func main() {
	app := awscdk.NewApp(nil)

	NewEKSStack(app, "anton-EKSStack", &EKSStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("406477933661"),
		Region:  jsii.String("us-west-2"),
	}
}
