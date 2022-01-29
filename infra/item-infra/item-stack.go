package iteminfra

import (
	common "infra/common"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	route53targets "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
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
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	itemTable := buildItemTable(stack, props, removalPolicy)
	bailInfoTable := buildBailInfoTable(stack, props, removalPolicy)
	inventoryTable := buildInventoryTable(stack, props, removalPolicy)
	buildLambda(stack, itemTable, bailInfoTable, inventoryTable, props)
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
func buildLambda(stack awscdk.Stack, itemTable, bailInfoTable, inventoryTable dynamodb.Table, props *ItemStackProps) {

	env := common.GetEnv()
	env["ItemTable"] = itemTable.TableName()
	env["BailInfoTable"] = bailInfoTable.TableName()
	env["InventoryTable"] = inventoryTable.TableName()

	itemFunction := lambda.NewFunction(stack, jsii.String("item-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../item-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("item-lambda-fn"),
	})

	itemTable.GrantFullAccess(itemFunction)
	bailInfoTable.GrantFullAccess(itemFunction)
	inventoryTable.GrantFullAccess(itemFunction)

	itemApi := apigateway.NewLambdaRestApi(stack, jsii.String("ItemApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:             props.Stage,
		Handler:                   itemFunction,
		RestApiName:               jsii.String("ItemRestApi"),
		Proxy:                     jsii.Bool(false),
		Deploy:                    jsii.Bool(true),
		DisableExecuteApiEndpoint: jsii.Bool(false),
		EndpointTypes:             &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.ItemApiDomain.Url),
		},
	})
	bailApis := itemApi.Root().AddResource(jsii.String("bailInfo"), &apigateway.ResourceOptions{})

	bailInfoApi := bailApis.AddResource(jsii.String("{bailNo}"), &apigateway.ResourceOptions{})
	bailInfoApi.AddMethod(jsii.String("GET"), itemApi.Root().DefaultIntegration(), nil)

	bailInfoQuantityApis := bailApis.AddResource(jsii.String("quality"), &apigateway.ResourceOptions{})
	bailInfoQuantityApi := bailInfoQuantityApis.AddResource(jsii.String("{quality}"), &apigateway.ResourceOptions{})
	bailInfoQuantityApi.AddMethod(jsii.String("GET"), itemApi.Root().DefaultIntegration(), nil)

	salesApis := itemApi.Root().AddResource(jsii.String("sales"), &apigateway.ResourceOptions{})
	salesApis.AddMethod(jsii.String("POST"), itemApi.Root().DefaultIntegration(), nil)

	salesGetAllApi := salesApis.AddResource(jsii.String("getall"), &apigateway.ResourceOptions{})
	salesGetAllApi.AddMethod(jsii.String("POST"), itemApi.Root().DefaultIntegration(), nil)

	salesGetApi := salesApis.AddResource(jsii.String("{salesBillNumber}"), &apigateway.ResourceOptions{})
	salesGetApi.AddMethod(jsii.String("GET"), itemApi.Root().DefaultIntegration(), nil)
	salesGetApi.AddMethod(jsii.String("PUT"), itemApi.Root().DefaultIntegration(), nil)
	salesGetApi.AddMethod(jsii.String("DELETE"), itemApi.Root().DefaultIntegration(), nil)

	purchaseGetAllApi := itemApi.Root().AddResource(jsii.String("purchase"), &apigateway.ResourceOptions{})
	purchaseGetAllApi.AddMethod(jsii.String("POST"), itemApi.Root().DefaultIntegration(), nil)

	purchaseApis := purchaseGetAllApi.AddResource(jsii.String("getall"), &apigateway.ResourceOptions{})
	purchaseApis.AddMethod(jsii.String("POST"), itemApi.Root().DefaultIntegration(), nil)

	purchaseGetApi := purchaseApis.AddResource(jsii.String("{purchaseBillNumber}"), &apigateway.ResourceOptions{})
	purchaseGetApi.AddMethod(jsii.String("GET"), itemApi.Root().DefaultIntegration(), nil)
	purchaseGetApi.AddMethod(jsii.String("PUT"), itemApi.Root().DefaultIntegration(), nil)
	purchaseGetApi.AddMethod(jsii.String("DELETE"), itemApi.Root().DefaultIntegration(), nil)

	hostedZone := common.GetHostedZone(stack, jsii.String("itemHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("itemArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.ItemApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(itemApi)),
	})
}
