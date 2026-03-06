package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecspatterns"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func MockServerStack(scope constructs.Construct, id string, props *awscdk.StackProps) (awscdk.Stack, *string) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	vpc := awsec2.NewVpc(stack, jsii.String("MockServerVpc"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
	})

	cluster := awsecs.NewCluster(stack, jsii.String("MockServerCluster"), &awsecs.ClusterProps{
		Vpc: vpc,
	})

	service := awsecspatterns.NewApplicationLoadBalancedFargateService(
		stack,
		jsii.String("MockServerService"),
		&awsecspatterns.ApplicationLoadBalancedFargateServiceProps{
			Cluster:        cluster,
			Cpu:            jsii.Number(256),
			MemoryLimitMiB: jsii.Number(512),
			TaskImageOptions: &awsecspatterns.ApplicationLoadBalancedTaskImageOptions{
				Image: awsecs.ContainerImage_FromAsset(
					jsii.String("../"),
					&awsecs.AssetImageProps{
						File: jsii.String("Dockerfile.mockserver"),
					},
				),
				ContainerPort: jsii.Number(3001),
			},
			ListenerPort:       jsii.Number(80),
			PublicLoadBalancer: jsii.Bool(true),
		},
	)

	mockServerUrl := awscdk.Fn_Join(jsii.String(""), &[]*string{
		jsii.String("http://"),
		service.LoadBalancer().LoadBalancerDnsName(),
	})

	return stack, mockServerUrl
}
