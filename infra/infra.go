package main

import (
	"fmt"
	clientinfra "infra/client-infra"
	documentinfra "infra/document-infra"

	common "infra/common"
	hsncodeinfra "infra/hsn-code-infra"
	iteminfra "infra/item-infra"
	qualityinfra "infra/quality-infra"
	transporterinfra "infra/transportor-infra"
	vendorinfra "infra/vendor-infra"
	"os"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	jsii "github.com/aws/jsii-runtime-go"
)

const baseDomain = "inventory-management.click"

var infraStackProps = common.InfraStackProps{
	StackProps: awscdk.StackProps{
		Env: env(),
	},
	InfraEnv: common.InfraEnv{
		HostedZoneId:    "Z00125422MSTTNZDZ5PRV",
		CertificateArn:  "arn:aws:acm:ap-south-1:675174225340:certificate/4995bc1a-da1b-42b1-b8d8-504ad0b9826b",
		StackNamePrefix: "kp-textile",
		Domains: common.Domain{
			BaseApi: baseDomain,
			ClientApiDomain: common.DomainDetails{
				RecordName: "client-api",
				Url:        fmt.Sprintf("client-api.%s", baseDomain),
			},
			VendorApiDomain: common.DomainDetails{
				RecordName: "vendor-api",
				Url:        fmt.Sprintf("vendor-api.%s", baseDomain),
			},
			TransporterApiDomain: common.DomainDetails{
				RecordName: "transporter-api",
				Url:        fmt.Sprintf("transporter-api.%s", baseDomain),
			},
			HsnCodeApiDomain: common.DomainDetails{
				RecordName: "hsncode-api",
				Url:        fmt.Sprintf("hsncode-api.%s", baseDomain),
			},
			QualityApiDomain: common.DomainDetails{
				RecordName: "quality-api",
				Url:        fmt.Sprintf("quality-api.%s", baseDomain),
			},
			ItemApiDomain: common.DomainDetails{
				RecordName: "item-api",
				Url:        fmt.Sprintf("item-api.%s", baseDomain),
			},
			DocumentApiDomain: common.DomainDetails{
				RecordName: "doc-api",
				Url:        fmt.Sprintf("doc-api.%s", baseDomain),
			},
		},
		CommonStackProps: common.CommonStackProps{
			IsLocal: os.Getenv("isLocal"),
			Stage: &apigateway.StageOptions{
				StageName: jsii.String("Dev"),
			},
		},
	},
}

func main() {
	app := awscdk.NewApp(nil)

	clientinfra.NewClientStack(app, "ClientStack", &clientinfra.ClientStackProps{
		InfraStackProps: infraStackProps,
	})
	documentinfra.NewDocumentStack(app, "DocumentStack", &documentinfra.DocumentStackProps{
		InfraStackProps: infraStackProps,
	})
	hsncodeinfra.NewHsnCodeStack(app, "HsnCodeStack", &hsncodeinfra.HsnCodeStackProps{
		InfraStackProps: infraStackProps,
	})
	iteminfra.NewItemStack(app, "ItemStack", &iteminfra.ItemStackProps{
		InfraStackProps: infraStackProps,
	})
	qualityinfra.NewQualityStack(app, "QualityStack", &qualityinfra.QualityStackProps{
		InfraStackProps: infraStackProps,
	})
	transporterinfra.NewTransporterStack(app, "TransporterStack", &transporterinfra.TransporterStackProps{
		InfraStackProps: infraStackProps,
	})

	vendorinfra.NewVendorStack(app, "VendorStack", &vendorinfra.VendorStackProps{
		InfraStackProps: infraStackProps,
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
