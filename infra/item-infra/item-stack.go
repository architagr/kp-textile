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

const (
	purchaseIdIndexName     = "purchaseid-index"
	purchaseBillNoIndexName = "purchase-billno-index"
	salesIdIndexName        = "salesid-index"
	salesBillNoIndexName    = "sales-billno-index"
	challanNoIndexName      = "challanno-index"
	baleNoIndexName         = "baleno-index"
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

	salesTable := buildSalesTable(stack, props, removalPolicy)
	purchaseTable := buildPurchaseTable(stack, props, removalPolicy)
	baleTable := buildBaleTable(stack, props, removalPolicy)
	buildLambda(stack, purchaseTable, salesTable, baleTable, props)
	return stack
}

func buildSalesTable(stack awscdk.Stack, props *ItemStackProps, removalPolicy awscdk.RemovalPolicy) dynamodb.Table {
	salesTable := dynamodb.NewTable(stack, jsii.String("SalesTable"), &dynamodb.TableProps{
		TableName: jsii.String("sales-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("godownId"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("sortKey"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
	salesTable.AddGlobalSecondaryIndex(&dynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String(salesIdIndexName),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("salesId"),
			Type: dynamodb.AttributeType_STRING,
		},
		ProjectionType: dynamodb.ProjectionType_ALL,
	})
	salesTable.AddGlobalSecondaryIndex(&dynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String(salesBillNoIndexName),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("salesBillNo"),
			Type: dynamodb.AttributeType_STRING,
		},
		ProjectionType: dynamodb.ProjectionType_KEYS_ONLY,
	})
	salesTable.AddGlobalSecondaryIndex(&dynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String(challanNoIndexName),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("challanNo"),
			Type: dynamodb.AttributeType_STRING,
		},
		ProjectionType: dynamodb.ProjectionType_KEYS_ONLY,
	})
	return salesTable
}

func buildPurchaseTable(stack awscdk.Stack, props *ItemStackProps, removalPolicy awscdk.RemovalPolicy) dynamodb.Table {
	puchaseTable := dynamodb.NewTable(stack, jsii.String("PurchaseTable"), &dynamodb.TableProps{
		TableName: jsii.String("purchase-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("godownId"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("sortKey"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
	puchaseTable.AddGlobalSecondaryIndex(&dynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String(purchaseIdIndexName),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("purchaseId"),
			Type: dynamodb.AttributeType_STRING,
		},
		ProjectionType: dynamodb.ProjectionType_ALL,
	})
	puchaseTable.AddGlobalSecondaryIndex(&dynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String(purchaseBillNoIndexName),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("purchaseBillNo"),
			Type: dynamodb.AttributeType_STRING,
		},
		ProjectionType: dynamodb.ProjectionType_KEYS_ONLY,
	})
	return puchaseTable
}

func buildBaleTable(stack awscdk.Stack, props *ItemStackProps, removalPolicy awscdk.RemovalPolicy) dynamodb.Table {
	baleTable := dynamodb.NewTable(stack, jsii.String("BaleTable"), &dynamodb.TableProps{
		TableName: jsii.String("bale-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("godownId"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("sortKey"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
	baleTable.AddGlobalSecondaryIndex(&dynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String(baleNoIndexName),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("baleNo"),
			Type: dynamodb.AttributeType_STRING,
		},
		ProjectionType: dynamodb.ProjectionType_ALL,
	})
	return baleTable
}

func buildLambda(stack awscdk.Stack, purchaseTable, salesTable, baleTable dynamodb.Table, props *ItemStackProps) {

	env := common.GetEnv()
	env["PurchaseTable"] = purchaseTable.TableName()
	env["SalesTable"] = salesTable.TableName()
	env["BaleTable"] = baleTable.TableName()
	env["PurchaseIdIndex"] = jsii.String(purchaseIdIndexName)
	env["SalesIdIndex"] = jsii.String(salesIdIndexName)
	env["PurchaseBillNoIndex"] = jsii.String(purchaseBillNoIndexName)
	env["SalesBillNoIndex"] = jsii.String(salesBillNoIndexName)
	env["ChallanNoIndex"] = jsii.String(challanNoIndexName)
	env["BaleNoIndex"] = jsii.String(baleNoIndexName)

	itemFunction := lambda.NewFunction(stack, jsii.String("item-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../item-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("item-lambda-fn"),
	})

	purchaseTable.GrantFullAccess(itemFunction)
	salesTable.GrantFullAccess(itemFunction)
	baleTable.GrantFullAccess(itemFunction)

	itemApi := apigateway.NewLambdaRestApi(stack, jsii.String("ItemApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     itemFunction,
		RestApiName:                 jsii.String("ItemRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.ItemApiDomain.Url),
		},
	})

	integration := apigateway.NewLambdaIntegration(itemFunction, &apigateway.LambdaIntegrationOptions{})

	// bailApis := itemApi.Root().AddResource(jsii.String("bailInfo"), &apigateway.ResourceOptions{
	// 	DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	// })

	// bailInfoApi := bailApis.AddResource(jsii.String("{bailNo}"), &apigateway.ResourceOptions{
	// 	DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	// })
	// bailInfoApi.AddMethod(jsii.String("GET"), integration, nil)

	// bailInfoQuantityApis := bailApis.AddResource(jsii.String("quality"), &apigateway.ResourceOptions{
	// 	DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	// })
	// bailInfoQuantityApi := bailInfoQuantityApis.AddResource(jsii.String("{quality}"), &apigateway.ResourceOptions{
	// 	DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	// })
	// bailInfoQuantityApi.AddMethod(jsii.String("GET"), integration, nil)

	salesApis := itemApi.Root().AddResource(jsii.String("sales"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	salesApis.AddMethod(jsii.String("POST"), integration, nil)

	salesGetAllApi := salesApis.AddResource(jsii.String("getall"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	salesGetAllApi.AddMethod(jsii.String("POST"), integration, nil)

	// salesGetApi := salesApis.AddResource(jsii.String("{salesBillNumber}"), &apigateway.ResourceOptions{
	// 	DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	// })
	// salesGetApi.AddMethod(jsii.String("GET"), integration, nil)
	// salesGetApi.AddMethod(jsii.String("PUT"), integration, nil)
	// salesGetApi.AddMethod(jsii.String("DELETE"), integration, nil)

	purchaseGetAllApi := itemApi.Root().AddResource(jsii.String("purchase"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	purchaseGetAllApi.AddMethod(jsii.String("POST"), integration, nil)

	purchaseApis := purchaseGetAllApi.AddResource(jsii.String("getall"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	purchaseApis.AddMethod(jsii.String("POST"), integration, nil)

	// purchaseGetApi := purchaseApis.AddResource(jsii.String("{purchaseBillNumber}"), &apigateway.ResourceOptions{
	// 	DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	// })
	// purchaseGetApi.AddMethod(jsii.String("GET"), integration, nil)
	// purchaseGetApi.AddMethod(jsii.String("PUT"), integration, nil)
	// purchaseGetApi.AddMethod(jsii.String("DELETE"), integration, nil)

	hostedZone := common.GetHostedZone(stack, jsii.String("itemHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("itemArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.ItemApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(itemApi)),
	})
}
