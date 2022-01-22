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
	BaseApi                                                                                    string
	ClientApiDomain, VendorApiDomain, TransporterApiDomain, HsnCodeApiDomain, QualityApiDomain DomainDetails
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
