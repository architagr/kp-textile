package qualityinfra

import (
	common "infra/common"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type QualityStackProps struct {
	common.InfraStackProps
}

func NewQualityStack(scope constructs.Construct, id string, props *QualityStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	qualityTable := buildTable(stack, props)
	buildLambda(stack, qualityTable)
	return stack
}
func buildTable(stack awscdk.Stack, props *QualityStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("QualityTable"), &dynamodb.TableProps{
		TableName: jsii.String("quality-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("id"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("name"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
}
func buildLambda(stack awscdk.Stack, qualityTable dynamodb.Table) {

	env := make(map[string]*string)
	env["qualityTable"] = qualityTable.TableName()
	env["GIN_MODE"] = jsii.String("release")

	function := lambda.NewFunction(stack, jsii.String("QualityServiceLambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../quality-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("quality-service-int-lambda-fn"),
	})

	qualityTable.GrantFullAccess(function)
}
