package transporterinfra

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

type TransporterStackProps struct {
	common.InfraStackProps
}

func NewTransporterStack(scope constructs.Construct, id string, props *TransporterStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	transporterTable := buildTable(stack, props)
	buildLambda(stack, transporterTable, props)
	return stack
}
func buildTable(stack awscdk.Stack, props *TransporterStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("TransporterTable"), &dynamodb.TableProps{
		TableName: jsii.String("transporter-table"),
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
func buildLambda(stack awscdk.Stack, transporterTable dynamodb.Table, props *TransporterStackProps) {

	env := common.GetEnv()
	env["TransporterTable"] = transporterTable.TableName()

	transporterFunction := lambda.NewFunction(stack, jsii.String("transporter-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../transportor-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("transporter-lambda-fn"),
	})

	transporterTable.GrantFullAccess(transporterFunction)

	transporterApi := apigateway.NewLambdaRestApi(stack, jsii.String("TransporterApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     transporterFunction,
		RestApiName:                 jsii.String("TransporterRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.TransporterApiDomain.Url),
		},
	})

	integration := apigateway.NewLambdaIntegration(transporterFunction, &apigateway.LambdaIntegrationOptions{})

	apis := transporterApi.Root().AddResource(jsii.String("transporter"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	apis.AddMethod(jsii.String("POST"), integration, nil)

	api := apis.AddResource(jsii.String("{transporterId}"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})

	api.AddMethod(jsii.String("GET"), integration, nil)
	api.AddMethod(jsii.String("DELETE"), integration, nil)
	api.AddMethod(jsii.String("PUT"), integration, nil)

	api2 := apis.AddResource(jsii.String("getall"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	api2.AddMethod(jsii.String("POST"), integration, nil)

	hostedZone := common.GetHostedZone(stack, jsii.String("transporterHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("transporterArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.TransporterApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(transporterApi)),
	})
}
