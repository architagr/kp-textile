package hsncodeinfra

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

type HsnCodeStackProps struct {
	common.InfraStackProps
}

func NewHsnCodeStack(scope constructs.Construct, id string, props *HsnCodeStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	hsnCodeTable := buildTable(stack, props)
	buildLambda(stack, hsnCodeTable, props)
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
func buildLambda(stack awscdk.Stack, hsnCodeTable dynamodb.Table, props *HsnCodeStackProps) {

	env := common.GetEnv()
	env["hsnCodeTable"] = hsnCodeTable.TableName()

	hsnCodeFunction := lambda.NewFunction(stack, jsii.String("hsn-code-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../hsn-code-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("hsn-code-lambda-fn"),
	})

	hsnCodeTable.GrantFullAccess(hsnCodeFunction)

	hsnCodeApi := apigateway.NewLambdaRestApi(stack, jsii.String("HsnCodeApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     hsnCodeFunction,
		RestApiName:                 jsii.String("HsnCodeRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.HsnCodeApiDomain.Url),
		},
	})

	integration := apigateway.NewLambdaIntegration(hsnCodeFunction, &apigateway.LambdaIntegrationOptions{})

	apis := hsnCodeApi.Root().AddResource(jsii.String("hsncode"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	apis.AddMethod(jsii.String("GET"), integration, &apigateway.MethodOptions{})

	apis.AddMethod(jsii.String("POST"), integration, &apigateway.MethodOptions{})

	api := apis.AddResource(jsii.String("{id}"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	api.AddMethod(jsii.String("GET"), integration, &apigateway.MethodOptions{})

	api2 := apis.AddResource(jsii.String("addmultiple"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	api2.AddMethod(jsii.String("POST"), integration, &apigateway.MethodOptions{})

	hostedZone := common.GetHostedZone(stack, jsii.String("hsncodeHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("hsncodeArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.HsnCodeApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(hsnCodeApi)),
	})
}
