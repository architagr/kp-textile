package clientinfra

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

type ClientStackProps struct {
	common.InfraStackProps
}

func NewClientStack(scope constructs.Construct, id string, props *ClientStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	clientTable := buildTable(stack, props)
	buildLambda(stack, clientTable, props)
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
func buildLambda(stack awscdk.Stack, clientTable dynamodb.Table, props *ClientStackProps) {

	env := common.GetEnv()
	env["ClientTable"] = clientTable.TableName()

	clientFunction := lambda.NewFunction(stack, jsii.String("client-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../client-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("client-lambda-fn"),
	})

	clientTable.GrantFullAccess(clientFunction)

	clientApi := apigateway.NewLambdaRestApi(stack, jsii.String("ClientApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     clientFunction,
		RestApiName:                 jsii.String("ClientRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.ClientApiDomain.Url),
		},
	})

	integration := apigateway.NewLambdaIntegration(clientFunction, &apigateway.LambdaIntegrationOptions{})

	apis := clientApi.Root().AddResource(jsii.String("client"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	apis.AddMethod(jsii.String("POST"), integration, nil)

	api := apis.AddResource(jsii.String("{clientId}"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})

	api.AddMethod(jsii.String("GET"), integration, nil)
	api.AddMethod(jsii.String("DELETE"), integration, nil)
	api.AddMethod(jsii.String("PUT"), integration, nil)

	api2 := apis.AddResource(jsii.String("getall"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})

	api2.AddMethod(jsii.String("POST"), integration, nil)

	hostedZone := common.GetHostedZone(stack, jsii.String("clientHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("clientArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.ClientApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(clientApi)),
	})
}
