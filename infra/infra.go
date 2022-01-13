package main

import (
	clientinfra "infra/client-infra"
	common "infra/common"
	hsncodeinfra "infra/hsn-code-infra"
	qualityinfra "infra/quality-infra"

	"os"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type InfraEnv struct {
	StackNamePrefix string
}

type InfraStackProps struct {
	awscdk.StackProps
}

var customEnv InfraEnv = InfraEnv{
	StackNamePrefix: "kp-textile",
}

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)
	commonStackProps := common.CommonStackProps{
		IsLocal: os.Getenv("isLocal"),
	}
	hsncodeinfra.NewHsnCodeStack(app, "HsnCodeStack", &hsncodeinfra.HsnCodeStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		CommonStackProps: commonStackProps,
	})
	qualityinfra.NewQualityStack(app, "QualityStack", &qualityinfra.QualityStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		CommonStackProps: commonStackProps,
	})
	clientinfra.NewClientStack(app, "ClientStack", &clientinfra.ClientStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		CommonStackProps: commonStackProps,
	})
	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String("675174225340"),
		Region:  jsii.String("ap-south-1"),
	}

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
