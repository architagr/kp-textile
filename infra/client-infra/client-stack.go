package clientinfra

import (
	common "infra/common"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type ClientStackProps struct {
	common.CommonStackProps
	awscdk.StackProps
}

func NewClientStack(scope constructs.Construct, id string, props *ClientStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	clientTable := buildTable(stack, props)
	buildLambda(stack, clientTable)
	return stack
}
func buildTable(stack awscdk.Stack, props *ClientStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("ClientTable"), &dynamodb.TableProps{
		TableName: jsii.String("client-table"),
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
func buildLambda(stack awscdk.Stack, clientTable dynamodb.Table) {

	env := make(map[string]*string)
	env["ClientTable"] = clientTable.TableName()
	env["GIN_MODE"] = jsii.String("release")

	function := lambda.NewFunction(stack, jsii.String("client-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../client-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("client-int-lambda-fn"),
	})

	clientTable.GrantFullAccess(function)
}
