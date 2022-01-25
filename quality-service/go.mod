module quality-service

go 1.17

replace commonpkg => ./../commonpkg

require (
	commonpkg v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.28.0
	github.com/aws/aws-sdk-go v1.42.35
	github.com/awslabs/aws-lambda-go-api-proxy v0.12.0
	github.com/gin-gonic/gin v1.7.7
	github.com/iris-contrib/go.uuid v2.0.0+incompatible
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
