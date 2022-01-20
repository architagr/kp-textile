package vendorinfra

import (
	common "infra/common"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type VendorStackProps struct {
	common.CommonStackProps
	awscdk.StackProps
}

func NewVendorStack(scope constructs.Construct, id string, props *VendorStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	vendorTable := buildTable(stack, props)
	buildLambda(stack, vendorTable)
	return stack
}
func buildTable(stack awscdk.Stack, props *VendorStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("VendorTable"), &dynamodb.TableProps{
		TableName: jsii.String("vendor-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("branchId"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("sortKey"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
}
func buildLambda(stack awscdk.Stack, vendorTable dynamodb.Table) {

	env := make(map[string]*string)
	env["VendorTable"] = vendorTable.TableName()
	env["GIN_MODE"] = jsii.String("release")

	function := lambda.NewFunction(stack, jsii.String("vendor-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../vendor-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("vendor-int-lambda-fn"),
	})

	vendorTable.GrantFullAccess(function)
}
