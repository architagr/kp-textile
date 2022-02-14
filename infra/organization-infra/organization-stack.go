package organizationinfra

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

type OrganizationStackProps struct {
	common.InfraStackProps
}

func NewOrganizationStack(scope constructs.Construct, id string, props *OrganizationStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	godownTable := buildGodownTable(stack, props)
	userTable := buildUserTable(stack, props)
	buildLambda(stack, godownTable, userTable, props)
	return stack
}
func buildGodownTable(stack awscdk.Stack, props *OrganizationStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("GodownTable"), &dynamodb.TableProps{
		TableName: jsii.String("godown-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("id"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
}

func buildUserTable(stack awscdk.Stack, props *OrganizationStackProps) dynamodb.Table {
	var removalPolicy awscdk.RemovalPolicy = awscdk.RemovalPolicy_RETAIN

	if props.IsLocal == "" {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	} else {
		removalPolicy = awscdk.RemovalPolicy_DESTROY

	}

	return dynamodb.NewTable(stack, jsii.String("UserTable"), &dynamodb.TableProps{
		TableName: jsii.String("user-table"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("userName"),
			Type: dynamodb.AttributeType_STRING,
		},
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
}
func buildLambda(stack awscdk.Stack, godownTable dynamodb.Table, userTable dynamodb.Table, props *OrganizationStackProps) {

	env := common.GetEnv()
	env["godownTable"] = godownTable.TableName()
	env["userTable"] = userTable.TableName()

	organizationFunction := lambda.NewFunction(stack, jsii.String("OrganizationServiceLambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../organization-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("organization-service-lambda-fn"),
	})

	userTable.GrantFullAccess(organizationFunction)
	godownTable.GrantFullAccess(organizationFunction)

	organizationApi := apigateway.NewLambdaRestApi(stack, jsii.String("OrganizationApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     organizationFunction,
		RestApiName:                 jsii.String("OrganizationRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.OrganizationApiDomain.Url),
		},
	})

	integration := apigateway.NewLambdaIntegration(organizationFunction, &apigateway.LambdaIntegrationOptions{})

	godownApis := organizationApi.Root().AddResource(jsii.String("godown"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	godownApis.AddMethod(jsii.String("GET"), integration, nil)
	godownApis.AddMethod(jsii.String("POST"), integration, nil)

	userApis := organizationApi.Root().AddResource(jsii.String("user"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	pApi := userApis.AddResource(jsii.String("login"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	pApi.AddMethod(jsii.String("POST"), integration, nil)

	hostedZone := common.GetHostedZone(stack, jsii.String("organizationHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("organizationArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.OrganizationApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(organizationApi)),
	})
}
