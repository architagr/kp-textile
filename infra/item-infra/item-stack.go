package iteminfra

import (
	common "infra/common"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type ItemStackProps struct {
	common.InfraStackProps
}

func NewItemStack(scope constructs.Construct, id string, props *ItemStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	buildTable(stack, props)

	return stack
}

func buildTable(stack awscdk.Stack, props *ItemStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	itemTable := dynamodb.NewTable(stack, jsii.String("ItemTable"), &dynamodb.TableProps{
		TableName: jsii.String("item-table"),
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

	itemTable.AddGlobalSecondaryIndex(&dynamodb.GlobalSecondaryIndexProps{
		IndexName:      jsii.String("billingIndex"),
		ProjectionType: dynamodb.ProjectionType_ALL,
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("branchId"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("billNo"),
			Type: dynamodb.AttributeType_STRING,
		},
	})

	return itemTable
}
