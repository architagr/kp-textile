package documentinfra

import (
	common "infra/common"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	route53targets "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type DocumentStackProps struct {
	common.InfraStackProps
}

func NewDocumentStack(scope constructs.Construct, id string, props *DocumentStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	buildLambda(stack, props)
	return stack
}

func buildLambda(stack awscdk.Stack, props *DocumentStackProps) {

	env := common.GetEnv()

	docuementFunction := lambda.NewFunction(stack, jsii.String("document-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("internal-api"),
		Code:         lambda.Code_FromAsset(jsii.String("./../document-service/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("document-lambda-fn"),
	})

	documentApi := apigateway.NewLambdaRestApi(stack, jsii.String("DocumentApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     docuementFunction,
		RestApiName:                 jsii.String("DocumentRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
		DomainName: &apigateway.DomainNameOptions{
			Certificate: common.CreateAcmCertificate(stack, &props.InfraEnv),
			DomainName:  jsii.String(props.Domains.DocumentApiDomain.Url),
		},
	})
	integration := apigateway.NewLambdaIntegration(docuementFunction, &apigateway.LambdaIntegrationOptions{})

	apis := documentApi.Root().AddResource(jsii.String("challan"), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: common.GetCorsPreflightOptions(),
	})
	apis.AddMethod(jsii.String("POST"), integration, nil)

	hostedZone := common.GetHostedZone(stack, jsii.String("documentHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("documentArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.DocumentApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(documentApi)),
	})
}
