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

const (
	BailInfoIndexName  = "BailInfoIndex"
	InventoryIndexName = "InventoryIndex"
)

func NewItemStack(scope constructs.Construct, id string, props *ItemStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	buildItemTable(stack, props, removalPolicy)
	buildBailInfoTable(stack, props, removalPolicy)
	buildInventoryTable(stack, props, removalPolicy)
	return stack
}

func buildItemTable(stack awscdk.Stack, props *ItemStackProps, removalPolicy awscdk.RemovalPolicy) dynamodb.Table {
	return dynamodb.NewTable(stack, jsii.String("ItemTable"), &dynamodb.TableProps{
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
}

func buildBailInfoTable(stack awscdk.Stack, props *ItemStackProps, removalPolicy awscdk.RemovalPolicy) dynamodb.Table {
	return dynamodb.NewTable(stack, jsii.String("BailInfoTable"), &dynamodb.TableProps{
		TableName: jsii.String("bail-info-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("branchId"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("bailInfoSortKey"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
}
func buildInventoryTable(stack awscdk.Stack, props *ItemStackProps, removalPolicy awscdk.RemovalPolicy) dynamodb.Table {
	return dynamodb.NewTable(stack, jsii.String("InventoryTable"), &dynamodb.TableProps{
		TableName: jsii.String("inventory-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("branchId"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("inventorySortKey"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
}
