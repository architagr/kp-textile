package vendorinfra

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

type VendorStackProps struct {
	common.InfraStackProps
}

func NewVendorStack(scope constructs.Construct, id string, props *VendorStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	vendorTable := buildTable(stack, props)
	buildLambda(stack, vendorTable, props)
	return stack
}
func buildTable(stack awscdk.Stack, props *VendorStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("VendorTable"), &dynamodb.TableProps{
		TableName: jsii.String("vendor-table"),
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
func buildLambda(stack awscdk.Stack, vendorTable dynamodb.Table, props *VendorStackProps) {

	env := common.GetEnv()
	env["VendorTable"] = vendorTable.TableName()

	vendorFunction := lambda.NewFunction(stack, jsii.String("vendor-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../vendor-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("vendor-lambda-fn"),
	})

	vendorTable.GrantFullAccess(vendorFunction)

	vendorApi := apigateway.NewLambdaRestApi(stack, jsii.String("VendorApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     vendorFunction,
		RestApiName:                 jsii.String("VendorRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.VendorApiDomain.Url),
		},
	})
	integration := apigateway.NewLambdaIntegration(vendorFunction, &apigateway.LambdaIntegrationOptions{})

	apis := vendorApi.Root().AddResource(jsii.String("vendor"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	apis.AddMethod(jsii.String("POST"), integration, nil)

	api := apis.AddResource(jsii.String("{vendorId}"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})

	api.AddMethod(jsii.String("GET"), integration, nil)
	api.AddMethod(jsii.String("DELETE"), integration, nil)
	api.AddMethod(jsii.String("PUT"), integration, nil)

	api2 := apis.AddResource(jsii.String("getall"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	api2.AddMethod(jsii.String("POST"), integration, nil)

	hostedZone := common.GetHostedZone(stack, jsii.String("vendorHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("vendorArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.VendorApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(vendorApi)),
	})

}
