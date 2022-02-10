package qualityinfra

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

type QualityStackProps struct {
	common.InfraStackProps
}

func NewQualityStack(scope constructs.Construct, id string, props *QualityStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	qualityTable := buildQualityTable(stack, props)
	productTable := buildProductTable(stack, props)
	buildLambda(stack, qualityTable, productTable, props)
	return stack
}
func buildQualityTable(stack awscdk.Stack, props *QualityStackProps) dynamodb.Table {
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

func buildProductTable(stack awscdk.Stack, props *QualityStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("ProductTable"), &dynamodb.TableProps{
		TableName: jsii.String("product-table"),
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
func buildLambda(stack awscdk.Stack, qualityTable dynamodb.Table, productTable dynamodb.Table, props *QualityStackProps) {

	env := common.GetEnv()
	env["qualityTable"] = qualityTable.TableName()
	env["productTable"] = productTable.TableName()

	qualityFunction := lambda.NewFunction(stack, jsii.String("QualityServiceLambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../quality-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("quality-service-lambda-fn"),
	})

	qualityTable.GrantFullAccess(qualityFunction)
	productTable.GrantFullAccess(qualityFunction)

	qualityApi := apigateway.NewLambdaRestApi(stack, jsii.String("QualityApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     qualityFunction,
		RestApiName:                 jsii.String("QualityRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.QualityApiDomain.Url),
		},
	})

	integration := apigateway.NewLambdaIntegration(qualityFunction, &apigateway.LambdaIntegrationOptions{})

	qualityApis := qualityApi.Root().AddResource(jsii.String("quality"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	qualityApis.AddMethod(jsii.String("GET"), integration, nil)
	qualityApis.AddMethod(jsii.String("POST"), integration, nil)

	qApi := qualityApis.AddResource(jsii.String("{id}"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	qApi.AddMethod(jsii.String("GET"), integration, nil)

	qApi2 := qualityApis.AddResource(jsii.String("addmultiple"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	qApi2.AddMethod(jsii.String("POST"), integration, nil)

	productApis := qualityApi.Root().AddResource(jsii.String("product"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	productApis.AddMethod(jsii.String("GET"), integration, nil)
	productApis.AddMethod(jsii.String("POST"), integration, nil)

	pApi := productApis.AddResource(jsii.String("{id}"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	pApi.AddMethod(jsii.String("GET"), integration, nil)

	pApi2 := productApis.AddResource(jsii.String("addmultiple"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	pApi2.AddMethod(jsii.String("POST"), integration, nil)

	hostedZone := common.GetHostedZone(stack, jsii.String("qualityHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("qualityArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.QualityApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(qualityApi)),
	})
}
