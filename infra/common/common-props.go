package common

import (
	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	acm "github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/jsii-runtime-go"
)

type CommonStackProps struct {
	IsLocal string
	Stage   *apigateway.StageOptions
}

type InfraEnv struct {
	StackNamePrefix string
	Domains         Domain
	CommonStackProps
	HostedZoneId   string
	CertificateArn string
}
type DomainDetails struct {
	RecordName, Url string
}
type Domain struct {
	BaseApi                                                                  string
	ClientApiDomain, VendorApiDomain, TransporterApiDomain, HsnCodeApiDomain DomainDetails
	QualityApiDomain, ItemApiDomain, DocumentApiDomain                       DomainDetails
}
type InfraStackProps struct {
	awscdk.StackProps
	InfraEnv
}

func CreateAcmCertificate(stack awscdk.Stack, props *InfraEnv) acm.ICertificate {
	return acm.Certificate_FromCertificateArn(stack, jsii.String("clientApiCertificate"), &props.CertificateArn)
}

func GetHostedZone(stack awscdk.Stack, id *string, props InfraEnv) route53.IHostedZone {
	return route53.PublicHostedZone_FromHostedZoneAttributes(stack, id, &route53.HostedZoneAttributes{
		HostedZoneId: jsii.String(props.HostedZoneId),
		ZoneName:     jsii.String(props.Domains.BaseApi),
	})
}

func GetEnv() map[string]*string {
	env := make(map[string]*string)
	env["GIN_MODE"] = jsii.String("release")

	return env
}

func AllowCors() *apigateway.ResourceOptions {
	return &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: GetCorsPreflightOptions(),
		DefaultMethodOptions:        &apigateway.MethodOptions{},
	}
}

func GetCorsPreflightOptions() *apigateway.CorsOptions {
	return &apigateway.CorsOptions{
		AllowOrigins:     apigateway.Cors_ALL_ORIGINS(),
		AllowMethods:     apigateway.Cors_ALL_METHODS(),
		AllowHeaders:     jsii.Strings("Content-Type", "Authorization", "X-Amz-Date", "X-Api-Key"),
		AllowCredentials: jsii.Bool(true),
	}
}
