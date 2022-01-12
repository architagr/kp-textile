package hsncodeinfra

import (
	common "infra/common"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type HsnCodeStackProps struct {
	common.CommonStackProps
	awscdk.StackProps
}

func NewHsnCodeStack(scope constructs.Construct, id string, props *HsnCodeStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	hsnCodeTable := buildTable(stack, props)
	buildLambda(stack, hsnCodeTable)
	return stack
}
func buildTable(stack awscdk.Stack, props *HsnCodeStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("hsn-code-table"), &dynamodb.TableProps{
		TableName: jsii.String("hsn-code"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("id"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("hnsCode"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
}
func buildLambda(stack awscdk.Stack, hsnCodeTable dynamodb.Table) {

	env := make(map[string]*string)
	env["hsnCodeTable"] = hsnCodeTable.TableName()

	function := lambda.NewFunction(stack, jsii.String("hsn-code-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../hsn-code-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("hsn-code-int-lambda-fn"),
	})

	hsnCodeTable.GrantFullAccess(function)
}
